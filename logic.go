package main

import (
	"errors"
	"time"
)

// Sets the device report to the given data
func (f *Falcon8) setReport(data []byte) error {
	_, err := f.Device.Control(0x21, 0x09, 0x0307, 0x0002, data) // SET REPORT
	time.Sleep(DELAY_TIME)
	return err
}

// Modifies and fills data with the report from the device
func (f *Falcon8) getReport(data []byte) error {
	_, err := f.Device.Control(0xA1, 0x01, 0x0307, 0x0002, data) // GET REPORT
	time.Sleep(DELAY_TIME)
	return err
}

// Clear last 56 bytes of data
func (f *Falcon8) prepareSet2(data []byte) error {
	if len(data) != 264 {
		return errors.New("invalid byte array length, must be 264 for Falcon8")
	}
	data[0x01] = 0x02 // second byte goes from 0x82 to 0x02
	for i := range data[0xD0:] {
		data[0xD0+i] = 0x00 // clear last 56 bytes
	}
	return nil
}
