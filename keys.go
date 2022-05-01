package Falcon8

import "errors"

var (
	KeyIndexKeyMap = map[KeyIndex]byte{
		KeyIndex1: 0x08,
	}
)

type KeyControls struct {
	Keys map[KeyIndex]Key // seems to be bytes 0x08-0x0F
}

func (k *KeyControls) setByteArray(b []byte) error {
	if len(b) != 264 {
		return errors.New("invalid byte array length, must be 264 for Falcon8")
	}

	if k.Keys != nil {
		for k, v := range k.Keys {
			if !k.Valid() {
				continue
			}
			b[KeyIndexKeyMap[k]] = byte(v) // set key to activate when kth is pressed
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

	// Write Falcon-8 Key Controls to USB packet
	err = f.LEDControls.setByteArray(data)
	if err != nil {
		return err
	}

	err = f.setReport(data) // SET 2
	if err != nil {
		return err
	}

	data = nil
	data = make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x06 // Same for layers

	return f.setReport(data) // SET 3
}
