package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/gousb"
)

const (
	VENDOR_ID  = 0x0B05
	PRODUCT_ID = 0x1827

	BUFFER_SIZE = 16
)

type SupremeHiFiX struct {
	Context  *gousb.Context
	Device   *gousb.Device
	Intf     *gousb.Interface
	Endpoint *gousb.InEndpoint
}

func New() (*SupremeHiFiX, error) {
	s := new(SupremeHiFiX)
	err := s.loadDevice()
	return s, err
}

func (s *SupremeHiFiX) Close() {
	fmt.Println("Closing SupremeHiFiX")
	if s.Device != nil {
		s.Device.Close()
	}
	if s.Context != nil {
		s.Context.Close()
	}
	if s.Intf != nil {
		s.Intf.Close()
	}
}

func (s *SupremeHiFiX) loadDevice() error {
	s.Context = gousb.NewContext()
	device, err := s.Context.OpenDeviceWithVIDPID(VENDOR_ID, PRODUCT_ID)
	if err != nil {
		return err
	}
	s.Device = device
	return nil
}

func (s *SupremeHiFiX) read(endpoint *gousb.InEndpoint, interval time.Duration, maxSize int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		buf := make([]byte, maxSize)
		n, err := endpoint.Read(buf)
		if err != nil {
			panic(err)
		}
		data := buf[:n]
		fmt.Println(hex.Dump(data))
	}
}

func (s *SupremeHiFiX) setReport(data []byte) error {
	_, err := s.Device.Control(0x21, 0x09, 0x0200, 0x0000, data)
	return err
}

func main() {
	hifiX, err := New()
	if err != nil {
		panic(err)
	}
	defer hifiX.Close()

	hifiX.Device.SetAutoDetach(true)

	for num := range hifiX.Device.Desc.Configs {
		config, _ := hifiX.Device.Config(num)

		defer config.Close()

		for _, desc := range config.Desc.Interfaces {
			intf, err := config.Interface(desc.Number, 0)
			if err != nil {
				panic(err)
			}
			hifiX.Intf = intf

			for _, endpointDesc := range intf.Setting.Endpoints {
				// We only want to read, so we're looking for IN endpoints.
				if endpointDesc.Direction == gousb.EndpointDirectionIn {
					endpoint, err := intf.InEndpoint(endpointDesc.Number)
					if err != nil {
						panic(err)
					}
					hifiX.Endpoint = endpoint
					go hifiX.read(endpoint, endpointDesc.PollInterval, endpointDesc.MaxPacketSize)
				}
			}
		}
	}

	data0 := make([]byte, BUFFER_SIZE)
	data0[0x00] = 0xF9
	data0[0x01] = 0x21
	data0[0x02] = 0x06
	data0[0x04] = 0x02

	data1 := make([]byte, BUFFER_SIZE)
	data1[0x00] = 0xF9
	data1[0x01] = 0x21
	data1[0x02] = 0x06

	data2 := make([]byte, BUFFER_SIZE)
	data2[0x00] = 0xFA
	data2[0x01] = 0x21
	data2[0x02] = 0x07
	data2[0x04] = 0xE8

	data3 := make([]byte, BUFFER_SIZE)
	data3[0x00] = 0xF9
	data3[0x01] = 0x21
	data3[0x02] = 0x0A
	data3[0x04] = 0x14

	hifiX.setReport(data0)
	hifiX.setReport(data1)
	hifiX.setReport(data2)
	hifiX.setReport(data3)
	hifiX.setReport(data3)
}
