package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image/color"
	"os"
	"os/signal"
	"time"

	"github.com/google/gousb"
)

const (
	VENDOR_ID  = 0x195D // Itron Technology iONE
	PRODUCT_ID = 0x6009 // Unknown
)

type Falcon8 struct {
	Context  *gousb.Context
	Device   *gousb.Device
	Intf     *gousb.Interface
	Endpoint *gousb.InEndpoint

	LEDControls *LEDControls
}

func New() (*Falcon8, error) {
	f := new(Falcon8)
	f.LEDControls = new(LEDControls)
	err := f.loadDevice()
	return f, err
}

func (f *Falcon8) Close() {
	fmt.Println("Closing Falcon8")
	if f.Device != nil {
		f.Device.Close()
	}
	if f.Context != nil {
		f.Context.Close()
	}
	if f.Intf != nil {
		f.Intf.Close()
	}
}

func (f *Falcon8) loadDevice() error {
	f.Context = gousb.NewContext()
	device, err := f.Context.OpenDeviceWithVIDPID(VENDOR_ID, PRODUCT_ID)
	if err != nil {
		return err
	}
	f.Device = device
	return nil
}

func (f *Falcon8) read(endpoint *gousb.InEndpoint, interval time.Duration, maxSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		buff := make([]byte, maxSize)
		n, err := endpoint.Read(buff)
		if err != nil {
			panic(err)
		}
		data := buff[:n]
		fmt.Println(hex.Dump(data))
	}
}

// Sets the device report to the given data
func (f *Falcon8) setReport(data []byte) error {
	_, err := f.Device.Control(0x21, 0x09, 0x0307, 0x0002, data) // SET REPORT
	return err
}

// Modifies and fills data with the report from the device
func (f *Falcon8) getReport(data []byte) error {
	_, err := f.Device.Control(0xA1, 0x01, 0x0307, 0x0002, data) // GET REPORT
	return err
}

func (f *Falcon8) prepareLEDWrite(data []byte) error {
	if len(data) != 264 {
		return errors.New("invalid byte array length, must be 264 for Falcon8")
	}
	data[0x01] = 0x02 // second byte goes from 0x82 to 0x02
	for i := range data[0xD0:] {
		data[0xD0+i] = 0x00 // clear last 56 bytes
	}
	return nil
}

func (f *Falcon8) UpdateLEDs() error {
	data := make([]byte, 264)
	data[0] = 0x07
	data[1] = 0x82
	data[2] = 0x01

	fmt.Println(hex.Dump(data))
	err := f.setReport(data)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 100)
	err = f.getReport(data)
	if err != nil {
		return err
	}
	fmt.Println(hex.Dump(data))

	err = f.prepareLEDWrite(data)
	if err != nil {
		return err
	}
	err = f.LEDControls.SetByteArray(data)
	if err != nil {
		return err
	}

	fmt.Println(hex.Dump(data))
	err = f.setReport(data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	falcon8, err := New()

	if err != nil {
		panic(err)
	}
	fmt.Println(falcon8.Device.Desc)

	falcon8.Device.SetAutoDetach(true)

	for num := range falcon8.Device.Desc.Configs {
		config, _ := falcon8.Device.Config(num)

		// In a scenario where we have an error, we can continue
		// to the next config. Same is true for interfaces and
		// endpoints.
		defer config.Close()

		// Iterate through available interfaces for this configuration
		for _, desc := range config.Desc.Interfaces {
			intf, err := config.Interface(desc.Number, 0)
			if err != nil {
				panic(err)
			}
			falcon8.Intf = intf

			// Iterate through endpoints available for this interface.
			for _, endpointDesc := range intf.Setting.Endpoints {
				// We only want to read, so we're looking for IN endpoints.
				if endpointDesc.Direction == gousb.EndpointDirectionIn {
					endpoint, err := intf.InEndpoint(endpointDesc.Number)
					if err != nil {
						panic(err)
					}
					falcon8.Endpoint = endpoint
					go falcon8.read(endpoint, endpointDesc.PollInterval, endpointDesc.MaxPacketSize)
					// When we get here, we have an endpoint where we can
					// read data from the USB device
				}
			}
		}
	}

	// blue
	falcon8.LEDControls.SetLEDMode(LEDMODE_BREATHING).SetBrightness(BRIGHTNESS_LOW).SetFlow(FLOW_SPINNING).SetColor(color.RGBA{0, 0, 255, 255})
	falcon8.UpdateLEDs()

	// create interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		sig := <-c
		fmt.Printf("\nSignal received: %s\n", sig)
		falcon8.Close()
		os.Exit(0)
	}()

	<-make(chan struct{})
}
