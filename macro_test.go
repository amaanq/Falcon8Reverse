package Falcon8

import (
	"reflect"
	"testing"
)

func TestKeyData_ToBytes(t *testing.T) {
	type fields struct {
		KeyPress KeyPress
		Delay    Delay
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{name: "4660ms Down", fields: fields{KeyPress: KeyPressDOWN, Delay: 466}, want: []byte{0x81, 0xD2}},
		{name: "4660ms Up", fields: fields{KeyPress: KeyPressUP, Delay: 466}, want: []byte{0x01, 0xD2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KeyData{
				KeyPress: tt.fields.KeyPress,
				Delay:    tt.fields.Delay,
			}
			if got := k.ToBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyData.ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
