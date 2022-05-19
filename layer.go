package Falcon8

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

// Call this to update the active layer on the device
func (f *Falcon8) UpdateLayer() error {
	if !f.ActiveLayer.Valid() {
		return ErrInvalidLayer
	}

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
	data[0x01] = 0x06 // Same for keys

	return f.setReport(data) // SET 3
}
