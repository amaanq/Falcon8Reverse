package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/google/gousb"
)

const (
	VENDOR_ID  = 0x195D // Itron Technology iONE
	PRODUCT_ID = 0x6009 // Unknown

	DELAY_TIME = time.Millisecond * 50 // Someone else can test this for themselves, this works for me
)

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
	if false {
		falcon8.LEDControls.SetLEDMode(LEDMODE_BREATHING).SetBrightness(BRIGHTNESS_MAX).SetFlow(FLOW_SPINNING).SetColor(color.RGBA{255, 0, 0, 255})
		falcon8.UpdateLEDs()
	}

	// falcon8.SetLayer(Layer1)
	// err = falcon8.UpdateLayer()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// falcon8.UpdateLEDs()

	go func() {
		for {
			for i := KeyIndex1; i <= KeyIndex8; i++ {
				r := uint8(rand.Intn(255))
				g := uint8(rand.Intn(255))
				b := uint8(rand.Intn(255))
				falcon8.LEDControls.SetLEDMode(LEDMODE_CUSTOM).SetKeyColor(i, color.RGBA{r, g, b, 0})
				fmt.Println("Set color to:", color.White)
			}
			err = falcon8.UpdateLEDs()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

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
