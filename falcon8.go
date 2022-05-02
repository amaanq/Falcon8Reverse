package Falcon8

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/gousb"
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
			fmt.Println(iface)
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
		} else {
			fmt.Println("Falcon-8 config 1 closed")
		}
	}

	if falcon8.Device != nil {
		err = falcon8.Device.Close()
		if err != nil {
			fmt.Println("Falcon-8 error closing device:", err)
		} else {
			fmt.Println("Falcon-8 device closed")
		}
	}

	if falcon8.Context != nil {
		err = falcon8.Context.Close()
		if err != nil {
			fmt.Println("Falcon-8 error closing context:", err)
		} else {
			fmt.Println("Falcon-8 context closed")
		}
	}
}

func (falcon8 *Falcon8) read(endpoint *gousb.InEndpoint, interval time.Duration, maxSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		buff := make([]byte, maxSize)

		n, err := endpoint.Read(buff)
		if err != nil {
			break
		}
		data := buff[:n]

		fmt.Println(hex.Dump(data)) // Logger to be removed later..
	}
}
