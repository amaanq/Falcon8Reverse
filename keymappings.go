package Falcon8

type Key uint8

const (
	// First row, first key (top left)
	Key1 Key = iota
	// First row, second key
	Key2
	// First row, third key
	Key3
	// First row, fourth key (top right)
	Key4
	// Second row, first key (bottom left)
	Key5
	// Second row, second key
	Key6
	// Second row, third key
	Key7
	// Second row, fourth key (bottom right)
	Key8
)

func (k Key) Valid() bool {
	return k >= Key1 && k <= Key8
}

// Struct to map all keys to their byte indexes in the main and macro USB packets
type KeyByteMapping struct {
	// The byte index for the Key HID Code
	KeyIndex byte
	// The byte index for the Key's individual LED Colors
	KeyColorIndexes [3]byte
	// The byte index for the Key's macro
	KeyMacroIndex byte
}

var (
	KeyMappings map[Key]KeyByteMapping = map[Key]KeyByteMapping{
		Key1: {
			KeyIndex:        0x08,
			KeyColorIndexes: [3]byte{0x3A, 0x53, 0x6C},
			KeyMacroIndex:   0x00,
		},
		Key2: {
			KeyIndex:        0x0D,
			KeyColorIndexes: [3]byte{0x3B, 0x54, 0x6D},
			KeyMacroIndex:   0x05,
		},
		Key3: {
			KeyIndex:        0x12,
			KeyColorIndexes: [3]byte{0x40, 0x59, 0x72},
			KeyMacroIndex:   0x0A,
		},
		Key4: {
			KeyIndex:        0x17,
			KeyColorIndexes: [3]byte{0x3C, 0x55, 0x6E},
			KeyMacroIndex:   0x0F,
		},
		Key5: {
			KeyIndex:        0x09,
			KeyColorIndexes: [3]byte{0x3F, 0x58, 0x71},
			KeyMacroIndex:   0x01,
		},
		Key6: {
			KeyIndex:        0x0E,
			KeyColorIndexes: [3]byte{0x45, 0x5E, 0x77},
			KeyMacroIndex:   0x06,
		},
		Key7: {
			KeyIndex:        0x13,
			KeyColorIndexes: [3]byte{0x4A, 0x63, 0x7C},
			KeyMacroIndex:   0x0B,
		},
		Key8: {
			KeyIndex:        0x18,
			KeyColorIndexes: [3]byte{0x41, 0x5A, 0x73},
			KeyMacroIndex:   0x10,
		},
	}
)
