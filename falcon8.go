package Falcon8

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/gousb"
)

const (
	maxPacketSize = 64
)

type Falcon8 struct {
	Context    *gousb.Context
	Device     *gousb.Device
	Config     *gousb.Config
	Interfaces []*gousb.Interface
	// Endpoints  []*gousb.InEndpoint

	isOpen bool

	ActiveLayer Layer
	LEDControls *LEDControls
	KeyControls *KeyControls
}

func New() *Falcon8 {
	falcon8 := new(Falcon8)

	falcon8.Context = gousb.NewContext()

	falcon8.ActiveLayer = Layer1
	falcon8.LEDControls = new(LEDControls)
	falcon8.KeyControls = new(KeyControls)

	return falcon8
}

func (falcon8 *Falcon8) Open() error {
	if falcon8.isOpen {
		return errors.New("falcon8: device already open")
	}

	fmt.Println("Opening Falcon8")

	device, err := falcon8.Context.OpenDeviceWithVIDPID(VENDOR_ID, PRODUCT_ID)
	if err != nil {
		return err
	}
	falcon8.Device = device

	err = falcon8.Device.SetAutoDetach(true)
	if err != nil {
		return err
	}

	falcon8.Config, err = falcon8.Device.Config(1)
	if err != nil {
		return err
	}
	falcon8.Interfaces = make([]*gousb.Interface, len(falcon8.Config.Desc.Interfaces))

	// https://github.com/sferris/howler-controller/blob/443a2564a9475281d38ab3cc2758169a33ce920e/howler.go#L103
	// Claim all interfaces so that when we're done, they're all released
	// properly. (Or else the OS doesn't reclaim them)
	for n, desc := range falcon8.Config.Desc.Interfaces {
		iface, err := falcon8.Config.Interface(desc.Number, 0)
		if err != nil {
			fmt.Printf("Falcon-8: error claiming interface: %s\n", err.Error())
		} else {
			falcon8.Interfaces[n] = iface
		}
	}

	if falcon8.Interfaces[0] == nil {
		return fmt.Errorf("Falcon-8: failed to claim howler config interface")
	}

	return nil
}

func (falcon8 *Falcon8) Close() {
	var err error
	fmt.Println("Closing Falcon8")

	if falcon8.Interfaces != nil && len(falcon8.Interfaces) > 0 {
		for _, iface := range falcon8.Interfaces {
			iface.Close()
		}
	}

	if falcon8.Config != nil {
		err = falcon8.Config.Close()
		if err != nil {
			fmt.Println("Falcon-8 error closing config:", err)
		}
	}

	if falcon8.Device != nil {
		err = falcon8.Device.Close()
		if err != nil {
			fmt.Println("Falcon-8 error closing device:", err)
		}
	}

	if falcon8.Context != nil {
		err = falcon8.Context.Close()
		if err != nil {
			fmt.Println("Falcon-8 error closing context:", err)
		}
	}
}

// Read will read from the given endpoint and return only the data that was read.
func (falcon8 *Falcon8) Read(endpoint *gousb.InEndpoint) ([]byte, error) {
	buff := make([]byte, maxPacketSize)
	n, err := endpoint.Read(buff)
	if err != nil {
		return nil, err
	}
	return buff[:n], nil
}

// Loop through all interfaces and read from all endpoints that are directionally IN.
// pass in a byte channel to buffer to continuously read from if desired,
// pass in a stop channel to stop the loop, but bear in mind you have to send n bools to the stop channel
// where n is the number of total endpoints.
//
// Returns number of endpoints that are being read from. (This is useful for stopping the loop)
func (falcon8 *Falcon8) ReadAll(buffer chan []byte, stop chan bool, print bool) int {
	i := 0
	for _, iface := range falcon8.Interfaces {
		for _, endpoint := range iface.Setting.Endpoints {
			if endpoint.Direction == gousb.EndpointDirectionIn {
				endpoint, err := iface.InEndpoint(endpoint.Number)
				if err != nil {
					fmt.Println("Falcon-8: error getting in endpoint:", err)
					continue
				}
				if endpoint == nil {
					fmt.Println("Falcon-8: failed to get in endpoint")
					continue
				}
				i++
				go falcon8.ReadLoop(endpoint, buffer, stop, print)
			}
		}
	}
	return i
}

// Loop and continuously read from endpoint,
// pass in a byte channel to buffer to continuously read from if desired,
// pass in a stop channel to stop the loop when recevied from,
// and pass in true for print to print the buffer to stdout.
func (falcon8 *Falcon8) ReadLoop(endpoint *gousb.InEndpoint, buffer chan []byte, stop chan bool, print bool) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			buff := make([]byte, maxPacketSize)

			n, err := endpoint.Read(buff)
			if err != nil {
				break
			}
			data := buff[:n]

			if buffer != nil {
				go func() {
					buffer <- data
				}()
			}

			if print {
				fmt.Println(hex.Dump(data))
			}
		case <-stop:
			return
		}
	}
}

func (falcon8 *Falcon8) SaveConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(falcon8.ActiveLayer)
	err = falcon8.setReport(data) // SET
	if err != nil {
		return err
	}

	err = falcon8.getReport(data) // GET
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(data[:0xD0])
	hashSum := hash.Sum(nil)

	_, err = file.Write(data[:0xD0])
	if err != nil {
		return err
	}
	_, err = file.Write(hashSum)
	return err
}

func (falcon8 *Falcon8) LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	n, err := file.Read(data[:0xD0])
	if err != nil {
		return err
	}
	if n != 0xD0 {
		return fmt.Errorf("Falcon-8: failed to read enough data bytes from file")
	}

	hash := sha256.New()
	hash.Write(data[:0xD0])
	hashSum := hash.Sum(nil)

	hashSumCompare := make([]byte, 32)
	n, err = file.Read(hashSumCompare)
	if err != nil {
		return err
	}
	if n != 32 {
		return fmt.Errorf("Falcon-8: failed to read enough checksum bytes from file")
	}

	if !bytes.Equal(hashSum, hashSumCompare) {
		return fmt.Errorf("Falcon-8: checksum mismatch")
	}

	falcon8.ActiveLayer = Layer(data[0x02])
	data[0x01] = 0x02 // only thing to change

	err = falcon8.setReport(data) // SET
	return err
}

// No checksum
func (falcon8 *Falcon8) unsafeSaveConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(falcon8.ActiveLayer)
	err = falcon8.setReport(data) // SET
	if err != nil {
		return err
	}

	err = falcon8.getReport(data) // GET
	if err != nil {
		return err
	}

	_, err = file.Write(data[:0xD0])
	return err
}

// No checksum
func (falcon8 *Falcon8) unsafeLoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	n, err := file.Read(data[:0xD0])
	if err != nil {
		return err
	}
	if n != 0xD0 {
		return fmt.Errorf("Falcon-8: failed to read enough data bytes from file")
	}

	falcon8.ActiveLayer = Layer(data[0x02])
	data[0x01] = 0x02 // only thing to change

	err = falcon8.setReport(data) // SET
	return err
}
