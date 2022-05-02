package Falcon8

import (
	"fmt"
	"testing"
)

func Test_Debug(t *testing.T) {
	falcon8 := New()

	err := falcon8.Open()
	if err != nil {
		panic(err)
	}

	falcon8.LEDControls.SetLEDMode(LEDMODE_BREATHING).SetBrightness(BRIGHTNESS_MAX)
	falcon8.UpdateLEDs()

	falcon8.KeyControls.SetKey(KeyIndex5, KEY_KPASTERISK)
	falcon8.KeyControls.SetKey(KeyIndex1, KEY_KPSLASH)
	falcon8.UpdateKeys()

	err = falcon8.unsafeLoadConfig("test.bin")
	if err != nil {
		fmt.Printf("error loading config: %s\n", err)
	}

	falcon8.Close()
}
