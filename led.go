package main

import (
	"errors"
	"image/color"
)

type KeyIndex uint8

const (
	// First row, first key (top left)
	KeyIndex1 KeyIndex = iota
	// First row, second key
	KeyIndex2
	// First row, third key
	KeyIndex3
	// First row, fourth key (top right)
	KeyIndex4
	// Second row, first key (bottom left)
	KeyIndex5
	// Second row, second key
	KeyIndex6
	// Second row, third key
	KeyIndex7
	// Second row, fourth key (bottom right)
	KeyIndex8
)

func (k KeyIndex) Valid() bool {
	return k >= KeyIndex1 && k <= KeyIndex8
}

type LEDMode byte

const (
	LEDMODE_NORMAL LEDMode = iota
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
	LEDMode    *LEDMode                 // 0x85
	Brightness *Brightness              // 0x86
	Flow       *Flow                    // 0x87
	Color      color.Color              // 0x88 - 0x8A
	KeyColors  map[KeyIndex]color.Color // sporadic occurrences
}

func (l *LEDControls) SetLEDMode(m LEDMode) *LEDControls {
	l.LEDMode = &m
	return l
}

func (l *LEDControls) SetBrightness(b Brightness) *LEDControls {
	l.Brightness = &b
	return l
}

func (l *LEDControls) SetFlow(f Flow) *LEDControls {
	l.Flow = &f
	return l
}

func (l *LEDControls) SetColor(c color.Color) *LEDControls {
	l.Color = c
	return l
}

// Pass in the key index and the color to set the LED to, mode must be set to LEDMODE_CUSTOM for this to work
func (l *LEDControls) SetKeyColor(k KeyIndex, c color.Color) *LEDControls {
	if !k.Valid() {
		return l
	}
	if l.KeyColors == nil {
		l.KeyColors = make(map[KeyIndex]color.Color)
	}
	l.KeyColors[k] = c
	return l
}

// Pass in the key index to turn off its LED, mode must be set to LEDMODE_CUSTOM for this to work
func (l *LEDControls) SetKeyColorDisabled(k KeyIndex) *LEDControls {
	if !k.Valid() {
		return l
	}
	if l.KeyColors == nil {
		l.KeyColors = make(map[KeyIndex]color.Color)
	}
	l.KeyColors[k] = color.Black
	return l
}

// Do not use or modify this function
func (l *LEDControls) SetKeyColorInByteArray(b []byte, k KeyIndex, c color.Color) {
	_r, _g, _b, _ := c.RGBA()
	switch k {
	case KeyIndex1:
		b[0x7A] = byte(_r)
		b[0x93] = byte(_g)
		b[0xAC] = byte(_b)
	case KeyIndex2:
		b[0x7B] = byte(_r)
		b[0x94] = byte(_g)
		b[0xAD] = byte(_b)
	case KeyIndex3:
		b[0x80] = byte(_r)
		b[0x99] = byte(_g)
		b[0xB2] = byte(_b)
	case KeyIndex4:
		b[0x7C] = byte(_r)
		b[0x95] = byte(_g)
		b[0xAE] = byte(_b)
	case KeyIndex5:
		b[0x7F] = byte(_r)
		b[0x98] = byte(_g)
		b[0xB1] = byte(_b)
	case KeyIndex6:
		b[0x85] = byte(_r)
		b[0x9E] = byte(_g)
		b[0xB7] = byte(_b)
	case KeyIndex7:
		b[0x8A] = byte(_r)
		b[0xA3] = byte(_g)
		b[0xBC] = byte(_b)
	case KeyIndex8:
		b[0x81] = byte(_r)
		b[0x9A] = byte(_g)
		b[0xB3] = byte(_b)
	}
}

func (l *LEDControls) SetByteArray(b []byte) error {
	if len(b) != 264 {
		return errors.New("invalid byte array length, must be 264 for Falcon8")
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
			l.SetKeyColorInByteArray(b, k, c)
		}
	}
	return nil
}
