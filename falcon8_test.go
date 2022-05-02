package Falcon8

import (
	"testing"
)

func Test_Debug(t *testing.T) {
	falcon8 := New()

	err := falcon8.Open()
	if err != nil {
		panic(err)
	}

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

	falcon8.Close()
}
