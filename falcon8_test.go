package Falcon8

import (
	//"fmt"
	//"os"
	//"os/signal"
	"testing"
)

func Test_Debug(t *testing.T) {
	falcon8 := New()

	err := falcon8.Open()
	if err != nil {
		panic(err)
	}

	// // Iterate through available interfaces for this configuration
	// for _, desc := range config.Desc.Interfaces {
	// 	if desc.Number != 2 {
	// 		continue
	// 	}
	// 	intf, err := config.Interface(desc.Number, 0)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	falcon8.Intf = intf

	// 	// Iterate through endpoints available for this interface.
	// 	for _, endpointDesc := range intf.Setting.Endpoints {
	// 		// We only want to read, so we're looking for IN endpoints.
	// 		if endpointDesc.Direction == gousb.EndpointDirectionIn {
	// 			endpoint, err := intf.InEndpoint(endpointDesc.Number)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			falcon8.Endpoint = endpoint
	// 			go falcon8.read(endpoint, endpointDesc.PollInterval, endpointDesc.MaxPacketSize)
	// 			// When we get here, we have an endpoint where we can
	// 			// read data from the USB device
	// 		}
	// 	}
	// }

	// Red
	// if false {
	// 	falcon8.LEDControls.SetLEDMode(LEDMODE_BREATHING).SetBrightness(BRIGHTNESS_MAX).SetFlow(FLOW_SPINNING).SetColor(color.RGBA{255, 0, 0, 255})
	// 	falcon8.UpdateLEDs()
	// }

	// falcon8.SetLayer(Layer1)
	// err = falcon8.UpdateLayer()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	falcon8.LEDControls.SetLEDMode(LEDMODE_BREATHING).SetBrightness(BRIGHTNESS_MAX)
	falcon8.UpdateLEDs()

	falcon8.KeyControls.SetKey(KeyIndex5, KEY_KPASTERISK)
	falcon8.KeyControls.SetKey(KeyIndex1, KEY_KPSLASH)
	falcon8.UpdateKeys()

	// go func() {
	// 	for {
	// 		for i := KeyIndex1; i <= KeyIndex8; i++ {
	// 			r := uint8(rand.Intn(255))
	// 			g := uint8(rand.Intn(255))
	// 			b := uint8(rand.Intn(255))
	// 			falcon8.LEDControls.SetLEDMode(LEDMODE_CUSTOM).SetKeyColor(i, color.RGBA{r, g, b, 0})
	// 		}
	// 		err = falcon8.UpdateLEDs()
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 	}
	// }()

	// create interrupt handler
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)

	// sig := <-c
	// fmt.Printf("\nSignal received: %s\n", sig)
	falcon8.Close()
}
