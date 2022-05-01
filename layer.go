package main

import (
	"errors"
)

type Layer byte

const (
	Layer1 Layer = iota + 1
	Layer2
	Layer3
	Layer4
	Layer5
)

func (l Layer) Valid() bool {
	return l >= Layer1 && l <= Layer5
}

func (f *Falcon8) SetLayer(layer Layer) {
	if layer.Valid() {
		f.ActiveLayer = layer
	}
}
func (f *Falcon8) UpdateLayer() error {
	// Set the active layer
	if !f.ActiveLayer.Valid() {
		return errors.New("falcon8: invalid layer")
	}

	data := make([]byte, 264)
	data[0x00] = 0x07
	data[0x01] = 0x82
	data[0x02] = byte(f.ActiveLayer) // LAYER
	err := f.setReport(data)         // SET 1
	if err != nil {
		return err
	}

	err = f.getReport(data) // GET
	if err != nil {
		return err
	}

	err = f.prepareSet2(data)
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
	data[0x01] = 0x06
	err = f.setReport(data) // SET 3
	if err != nil {
		return err
	}

	return nil
}
