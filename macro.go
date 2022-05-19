package Falcon8

import "time"

// set 0x07 0x82 0xLAYER PREPARE FOR GET
// get
// set DATA[0x01] = 0x02, KEY to 0xf8
// set MACRO data | 0x07 0x05 0xLAYER 0x00 0x01??
// set MACRO data | 0x07 0x05 0xLAYER 0x00 0x02??
// set MACRO data | 0x07 0x05 0xLAYER 0x00 0x03??
// set END data   | 0x07 0x06

// SINGLE KEY MACRO
// BYTE 0x08 - 0x09 DETERMINES # OF TIMES TO REPEAT, FF FF MEANS WHILE PRESSED, 00 00 MEANS UNTIL NEXT KEY IS PRESSED
// TYPE: 0x80 DOWN 0x00 UP
// LENGTH: 0xN * 10 MILLISECONDS
// LETTER: USB HID CODE OF LETTER
// TYPE MASKS WITH LENGTH LOL ITS A UINT16 !!!
// MSB SET, KEY DOWN | LSB SET, KEY UP
// OTHERWISE THE LAST 15 BYTES ARE THE TIME IN 10*MS MAKING THE MAX TIME 327670 MS = 327.670 S

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// BYTE 0x08 - 0x09
//_______________________________________________________________________

// Anything but 0x0000 or 0xFFFF means number of times to repeat the macro
type Repetition uint16

const (
	RepeatUntilNextKeyPressed Repetition = 0x0000
	RepeatWhilePressed        Repetition = 0xFFFF
)

func (r Repetition) ToBytes() []byte {
	MSB := byte(r >> 8)
	LSB := byte(r)
	// Big endian
	return []byte{MSB, LSB}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// BYTE 0x0A - END OF MACRO
//_______________________________________________________________________

type Delay int16     // max is 0x7FFF because MSB is reserved for KeyPress type
type KeyPress uint16 // 0x8000 or 0x0000, masked with Delay

const (
	KeyPressDOWN KeyPress = 0x8000
	KeyPressUP   KeyPress = 0x0000
)

const Duration = 10 * time.Millisecond

// Funny enough 0ms is actually 10ms in the Windows software.
func TimeToDelay(t time.Duration) Delay {
	if t < Duration {
		t = Duration
	}
	if t > 32767*Duration {
		t = 32767 * Duration
	}
	return Delay(t / Duration)
}

// In the form of 2 bytes, KeyPress | Delay
type KeyData struct {
	KeyPress KeyPress
	Delay    Delay
}

func (k KeyData) ToBytes() []byte {
	//  0 0 0 0 0 0 0 0  0 0 0 0 0 0 0 0
	//  1 1 1 1 1 1 1 1  1 1 1 1 1 1 1 1
	//  _ _ _ _ _ _ _ _  _ _ _ _ _ _ _ _  Last 15 bits are Delay as a factor of 10ms in big endian format (weird)
	//  🡅
	//  KeyPress bit 1 = DOWN 0 = UP

	MSB := byte(uint16(k.Delay>>8) | uint16(k.KeyPress>>8))
	LSB := byte(k.Delay)
	// Big Endian
	return []byte{MSB, LSB}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// 3-Byte Input
//_______________________________________________________________________

// This takes up 3 bytes, 2 from the KeyData and 1 from the Key itself
type Input struct {
	KeyData KeyData
	KeyCode KeyCode
}

// Return the input as a 3-byte array where the first two bytes are the KeyData and the third is the Key
func (i Input) ToBytes() [3]byte {
	//  Key Data + Key
	//  0 0 0 0 0 0 0 0  0 0 0 0 0 0 0 0  0 0 0 0 0 0 0 0 0
	//  First 2 bytes are key data        last byte is the key in its USB HID Code

	// First 2 bytes are key data
	keydataBytes := i.KeyData.ToBytes()
	// Last byte is the key
	keyByte := byte(i.KeyCode)
	// Concatenate
	return [3]byte{keydataBytes[0], keydataBytes[1], keyByte}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Entire Macro
//_______________________________________________________________________

type Macro struct {
	Repeats Repetition
	Inputs  []Input
}

func NewMacro() *Macro {
	m := new(Macro)
	m.Repeats = 1 // Do only once as default
	m.Inputs = make([]Input, 0)
	return m
}

func (m *Macro) SetRepeats(r Repetition) {
	m.Repeats = r
}

// Max limit of 240 inputs for any macro
func (m *Macro) AddInput(kp KeyPress, duration time.Duration, keyCode KeyCode) {
	if len(m.Inputs) >= 240 {
		return
	}
	input := Input{
		KeyData: KeyData{
			KeyPress: kp,
			Delay:    TimeToDelay(duration),
		},
		KeyCode: keyCode,
	}
	m.Inputs = append(m.Inputs, input)
}

// Take each input in the macro and convert it to its byte array, first 2 bytes are the repetition then each input represented as 3 bytes.
//
// The falcon 8 has 3 packets sending the macro, each packet is 264 bytes but has an 8 byte header so there are 256 bytes left effectively.
// However, the first packet also has the repetition packed in 2 bytes, so we need to subtract 2 bytes from the total (first packet has 254 usable bytes..)
// So, this gives us 766 bytes to use, which if divided into 3 for each input gives us 255 inputs, but the Falcon-8 software only permits a max of 240 so for safety we'll use 240.
// Note: returns an array of byte arrays of length 3 (3 packets)
func (m *Macro) ToBytes() [3][]byte {
	var b [3][]byte
	b[0] = make([]byte, 0)
	b[1] = make([]byte, 0)
	b[2] = make([]byte, 0)

	index := 0 // Index of the packet we're currently working on

	// First 2 bytes are the repetition
	b[index] = append(b[index], m.Repeats.ToBytes()...) // We don't have to check for packet overflow here

	// Then each input
	for i := range m.Inputs {
		bytes := m.Inputs[i].ToBytes()

		if len(b[index]) >= 256 {
			index++
		}
		b[index] = append(b[index], bytes[0])

		if len(b[index]) >= 256 {
			index++
		}
		b[index] = append(b[index], bytes[1])

		if len(b[index]) >= 256 {
			index++
		}
		b[index] = append(b[index], bytes[2])
	}
	return b
}