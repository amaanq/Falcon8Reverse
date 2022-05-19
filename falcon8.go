package Falcon8

import (
	"encoding/hex"
	"fmt"
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

// Create a new struct to interact with the Falcon-8 RGB Keypad
func New() *Falcon8 {
	falcon8 := new(Falcon8)

	falcon8.Context = gousb.NewContext()

	falcon8.ActiveLayer = Layer1
	falcon8.LEDControls = new(LEDControls)
	falcon8.KeyControls = new(KeyControls)

	return falcon8
}

// Opens the device for usage
//
// Note: The device will not be usable outside of the program this was called in until it is closed.
//
// All interfaces are claimed so that when closed they are all released back to the kernel.
func (falcon8 *Falcon8) Open() error {
	if falcon8.isOpen {
		return ErrDeviceAlreadyOpen
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
			return fmt.Errorf(ErrClaimingInterface.Error(), err.Error())
		} else {
			falcon8.Interfaces[n] = iface
		}
	}

	if falcon8.Interfaces[0] == nil {
		return ErrClaimingConfig
	}

	return nil
}

// Closes the device and releases all interfaces so that it is usable in any application.
// This must be called before exiting the program.
func (falcon8 *Falcon8) Close() error {
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
			return fmt.Errorf(ErrClosingConfig.Error(), err.Error())
		}
	}

	if falcon8.Device != nil {
		err = falcon8.Device.Close()
		if err != nil {
			return fmt.Errorf(ErrClosingDevice.Error(), err.Error())
		}
	}

	if falcon8.Context != nil {
		err = falcon8.Context.Close()
		if err != nil {
			return fmt.Errorf(ErrClosingContext.Error(), err.Error())
		}
	}
	return nil
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
					fmt.Println("falcon-8: error getting in endpoint:", err)
					continue
				}
				if endpoint == nil {
					fmt.Println("falcon-8: failed to get in endpoint")
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
