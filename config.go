package Falcon8

import (
	"bytes"
	"crypto/sha256"
	"os"
)

// Save the 208 bytes of intact data from the device.
// A hash of the data is also saved for sanity checking.
// TODO: Save macro packets too
func (falcon8 *Falcon8) SaveConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(falcon8.ActiveLayer)
	err = falcon8.setReport(data) // SET
	if err != nil {
		return err
	}

	err = falcon8.getReport(data) // GET
	if err != nil {
		return err
	}

	hash := sha256.New()
	hash.Write(data[:0xD0])
	hashSum := hash.Sum(nil)

	_, err = file.Write(data[:0xD0])
	if err != nil {
		return err
	}
	_, err = file.Write(hashSum)
	return err
}

// Load the 208 bytes of intact data from the device.
// This checks the hash to make sure the data is valid.
func (falcon8 *Falcon8) LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	n, err := file.Read(data[:0xD0])
	if err != nil {
		return err
	}
	if n != 0xD0 {
		return ErrBadCFGRead
	}

	hash := sha256.New()
	hash.Write(data[:0xD0])
	hashSum := hash.Sum(nil)

	hashSumCompare := make([]byte, 32)
	n, err = file.Read(hashSumCompare)
	if err != nil {
		return err
	}
	if n != 32 {
		return ErrBadCFGChecksumRead
	}

	if !bytes.Equal(hashSum, hashSumCompare) {
		return ErrBadCFGChecksum
	}

	falcon8.ActiveLayer = Layer(data[0x02])
	data[0x01] = 0x02 // only thing to change

	err = falcon8.setReport(data) // SET
	return err
}

// No checksum
func (falcon8 *Falcon8) unsafeSaveConfig(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(falcon8.ActiveLayer)
	err = falcon8.setReport(data) // SET
	if err != nil {
		return err
	}

	err = falcon8.getReport(data) // GET
	if err != nil {
		return err
	}

	_, err = file.Write(data[:0xD0])
	return err
}

// No checksum
func (falcon8 *Falcon8) unsafeLoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, 264)
	n, err := file.Read(data[:0xD0])
	if err != nil {
		return err
	}
	if n != 0xD0 {
		return ErrBadCFGRead
	}

	falcon8.ActiveLayer = Layer(data[0x02])
	data[0x01] = 0x02 // only thing to change

	err = falcon8.setReport(data) // SET
	return err
}
