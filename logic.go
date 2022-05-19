package Falcon8

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/gousb"
)

// Sets the device report to the given data
func (f *Falcon8) setReport(data []byte) error {
	_, err := f.Device.Control(gousb.ControlOut|gousb.ControlClass|gousb.ControlInterface, 0x09, 0x0307, 0x0002, data) // SET REPORT
	time.Sleep(DELAY_TIME)
	return err
}

// Modifies and fills data with the report from the device
func (f *Falcon8) getReport(data []byte) error {
	_, err := f.Device.Control(gousb.ControlIn|gousb.ControlClass|gousb.ControlInterface, 0x01, 0x0307, 0x0002, data) // GET REPORT
	time.Sleep(DELAY_TIME)
	return err
}

// Clear last 56 bytes of data and set data[0x01] to 0x02
func (f *Falcon8) prepareSet2(data []byte) error {
	if len(data) != 264 {
		return ErrInvalidByteArrayLength
	}

	data[0x01] = 0x02 // second byte goes from 0x82 to 0x02

	for i := range data[0xD0:] {
		data[0xD0+i] = 0x00 // clear last 56 bytes
	}

	return nil
}

// Pretty hex dump of packet data
func hexDump(name string, data []byte) {
	fmt.Printf(name + "\n" + hex.Dump(data))
}
