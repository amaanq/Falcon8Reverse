package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/gousb"
)

type Falcon8 struct {
	Context  *gousb.Context
	Device   *gousb.Device
	Intf     *gousb.Interface
	Endpoint *gousb.InEndpoint

	ActiveLayer Layer
	LEDControls *LEDControls
}

func New() (*Falcon8, error) {
	f := new(Falcon8)
	f.LEDControls = new(LEDControls)
	f.ActiveLayer = 0x01 // Layer 1
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
			break
		}
		data := buff[:n]
		fmt.Println(hex.Dump(data))
	}
}
