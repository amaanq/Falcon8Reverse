package Falcon8

type KeyControls struct {
	Keys map[Key]KeyCode
}

// Pass in the key index and the color to set the LED to, mode must be set to LEDMODE_CUSTOM otherwise this will have no effect.
func (k *KeyControls) SetKey(key Key, keyCode KeyCode) *KeyControls {
	if !key.Valid() {
		return k
	}
	if k.Keys == nil {
		k.Keys = make(map[Key]KeyCode)
	}
	k.Keys[key] = keyCode
	return k
}

func (k *KeyControls) setByteArray(b []byte) error {
	if len(b) != 264 {
		return ErrInvalidByteArrayLength
	}

	if k.Keys != nil {
		for k, v := range k.Keys {
			if !k.Valid() {
				continue
			}
			b[KeyMappings[k].KeyIndex] = byte(v) // set key to register when kth key on keypad is pressed
		}
	}
	return nil
}

// Call this to commit the changes to KeyControls to the device
func (f *Falcon8) UpdateKeys() error {
	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(f.ActiveLayer) // LAYER

	hexDump("KEY SET 1", data)
	err := f.setReport(data) // SET 1
	if err != nil {
		return err
	}

	err = f.getReport(data) // GET
	if err != nil {
		return err
	}
	hexDump("KEY GET 1", data)

	// Clear last 56 bytes, set byte 2 from 0x82 to 0x02 (read to write?)
	err = f.prepareSet2(data)
	if err != nil {
		return err
	}

	// Write Falcon-8 Key Controls to USB packet
	err = f.KeyControls.setByteArray(data)
	if err != nil {
		return err
	}

	hexDump("KEY SET 2", data)
	err = f.setReport(data) // SET 2
	if err != nil {
		return err
	}

	data = nil
	data = make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x06 // Same for layers

	hexDump("KEY SET 3", data)
	return f.setReport(data) // SET 3
}
