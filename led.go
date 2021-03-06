package Falcon8

import (
	"image/color"
)

type LEDMode byte

const (
	LEDMODE_STATIC LEDMode = iota
	LEDMODE_BREATHING
	LEDMODE_FADE_IN
	LEDMODE_FADE_OUT
	LEDMODE_LAST_KEYSTROKE
	LEDMODE_RGB_WAVE
	LEDMODE_RGB_RANDOM_SINGLE_KEY
	LEDMODE_CUSTOM

	LEDModeIndex = 0x85
)

type Brightness byte

const (
	BRIGHTNESS_OFF Brightness = iota
	BRIGHTNESS_LOW
	BRIGHTNESS_MEDIUM
	BRIGHTNESS_HIGH
	BRIGHTNESS_MAX

	BrightnessIndex = 0x86
)

type Flow byte

const (
	FLOW_RIGHT_TO_LEFT Flow = iota
	FLOW_LEFT_TO_RIGHT
	FLOW_TOP_TO_BOTTOM
	FLOW_SPINNING
	FLOW_CONSTANT

	FlowIndex = 0x87
)

type LEDControls struct {
	LEDMode    *LEDMode            // 0x85
	Brightness *Brightness         // 0x86
	Flow       *Flow               // 0x87
	Color      color.Color         // 0x88 - 0x8A
	KeyColors  map[Key]color.Color // sporadic occurrences, refer to KeyMappings
}

func (l *LEDControls) SetLEDMode(m LEDMode) *LEDControls {
	l.LEDMode = &m
	return l
}

func (l *LEDControls) SetBrightness(b Brightness) *LEDControls {
	l.Brightness = &b
	return l
}

// Only works with LEDMODE_RGB_WAVE
func (l *LEDControls) SetFlow(f Flow) *LEDControls {
	l.Flow = &f
	return l
}

// Set the overall color of the LEDs, works with LEDMODE_{NORMAL, BREATHING, FADE_IN, FADE_OUT, LAST_KEYSTROKE}
func (l *LEDControls) SetColor(c color.Color) *LEDControls {
	l.Color = c
	return l
}

// Pass in the key index and the color to set the LED to, mode must be set to LEDMODE_CUSTOM otherwise this will have no effect.
func (l *LEDControls) SetKeyColor(k Key, c color.Color) *LEDControls {
	if !k.Valid() {
		return l
	}
	if l.KeyColors == nil {
		l.KeyColors = make(map[Key]color.Color)
	}
	l.KeyColors[k] = c
	return l
}

// Pass in the key index to turn off its LED, mode must be set to LEDMODE_CUSTOM otherwise this will have no effect.
func (l *LEDControls) DisableKeyColor(k Key) *LEDControls {
	if !k.Valid() {
		return l
	}
	if l.KeyColors == nil {
		l.KeyColors = make(map[Key]color.Color)
	}
	l.KeyColors[k] = color.Black
	return l
}

// Do not use or modify this function
func (l *LEDControls) setKeyColorsInByteArray(b []byte, key Key, c color.Color) {
	_r, _g, _b, _ := c.RGBA()

	// Get the Key Color Indexes of a particular key where order is R G B, then set that byte in b to the color value
	b[KeyMappings[key].KeyColorIndexes[0x00]] = byte(_r)
	b[KeyMappings[key].KeyColorIndexes[0x01]] = byte(_g)
	b[KeyMappings[key].KeyColorIndexes[0x02]] = byte(_b)
}

// Writes Mode, Brightness, Flow, Color and KeyColors to the byte array
func (l *LEDControls) setByteArray(b []byte) error {
	if len(b) != 264 {
		return ErrInvalidByteArrayLength
	}

	if l.LEDMode != nil {
		b[LEDModeIndex] = byte(*l.LEDMode)
	}

	if l.Brightness != nil {
		b[BrightnessIndex] = byte(*l.Brightness)
	}

	if l.Flow != nil {
		b[FlowIndex] = byte(*l.Flow)
	}

	if l.Color != nil {
		_r, _g, _b, _ := l.Color.RGBA()
		b[0x88] = byte(_r)
		b[0x89] = byte(_g)
		b[0x8A] = byte(_b)
	}

	if l.KeyColors != nil {
		for k, c := range l.KeyColors {
			if !k.Valid() {
				continue
			}
			l.setKeyColorsInByteArray(b, k, c)
		}
	}
	return nil
}

// Call this to commit the changes to LEDControls to the device
func (f *Falcon8) UpdateLEDs() error {
	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(f.ActiveLayer) // LAYER

	err := f.setReport(data) // SET 1
	if err != nil {
		return err
	}

	err = f.getReport(data) // GET
	if err != nil {
		return err
	}

	// Clear last 56 bytes, set byte 2 from 0x82 to 0x02 (read to write?)
	err = f.prepareSet2(data)
	if err != nil {
		return err
	}

	// Write Falcon-8 LED controls to USB packet
	err = f.LEDControls.setByteArray(data)
	if err != nil {
		return err
	}

	return f.setReport(data) // SET 2
}
