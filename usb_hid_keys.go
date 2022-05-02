package Falcon8

const (
	KEY_NONE    Key = iota // No key pressed
	KEY_ERR_OVF            //  Keyboard Error Roll Over - used for all slots if too many keys are pressed ("Phantom key")
	_                      //  //  Keyboard POST Fail
	_                      //  //  Keyboard Error Undefined
	KEY_A                  // Keyboard a and A
	KEY_B                  // Keyboard b and B
	KEY_C                  // Keyboard c and C
	KEY_D                  // Keyboard d and D
	KEY_E                  // Keyboard e and E
	KEY_F                  // Keyboard f and F
	KEY_G                  // Keyboard g and G
	KEY_H                  // Keyboard h and H
	KEY_I                  // Keyboard i and I
	KEY_J                  // Keyboard j and J
	KEY_K                  // Keyboard k and K
	KEY_L                  // Keyboard l and L
	KEY_M                  // Keyboard m and M
	KEY_N                  // Keyboard n and N
	KEY_O                  // Keyboard o and O
	KEY_P                  // Keyboard p and P
	KEY_Q                  // Keyboard q and Q
	KEY_R                  // Keyboard r and R
	KEY_S                  // Keyboard s and S
	KEY_T                  // Keyboard t and T
	KEY_U                  // Keyboard u and U
	KEY_V                  // Keyboard v and V
	KEY_W                  // Keyboard w and W
	KEY_X                  // Keyboard x and X
	KEY_Y                  // Keyboard y and Y
	KEY_Z                  // Keyboard z and Z

	KEY_1 // Keyboard 1 and !
	KEY_2 // Keyboard 2 and @
	KEY_3 // Keyboard 3 and #
	KEY_4 // Keyboard 4 and $
	KEY_5 // Keyboard 5 and %
	KEY_6 // Keyboard 6 and ^
	KEY_7 // Keyboard 7 and &
	KEY_8 // Keyboard 8 and *
	KEY_9 // Keyboard 9 and (
	KEY_0 // Keyboard 0 and )

	KEY_ENTER      // Keyboard Return (ENTER)
	KEY_ESC        // Keyboard ESCAPE
	KEY_BACKSPACE  // Keyboard DELETE (Backspace)
	KEY_TAB        // Keyboard Tab
	KEY_SPACE      // Keyboard Spacebar
	KEY_MINUS      // Keyboard - and _
	KEY_EQUAL      // Keyboard = and +
	KEY_LEFTBRACE  // Keyboard [ and {
	KEY_RIGHTBRACE // Keyboard ] and }
	KEY_BACKSLASH  // Keyboard \ and |
	KEY_HASHTILDE  // Keyboard Non-US # and ~
	KEY_SEMICOLON  // Keyboard ; and :
	KEY_APOSTROPHE // Keyboard ' and "
	KEY_GRAVE      // Keyboard ` and ~
	KEY_COMMA      // Keyboard , and <
	KEY_DOT        // Keyboard . and >
	KEY_SLASH      // Keyboard / and ?
	KEY_CAPSLOCK   // Keyboard Caps Lock

	KEY_F1  // Keyboard F1
	KEY_F2  // Keyboard F2
	KEY_F3  // Keyboard F3
	KEY_F4  // Keyboard F4
	KEY_F5  // Keyboard F5
	KEY_F6  // Keyboard F6
	KEY_F7  // Keyboard F7
	KEY_F8  // Keyboard F8
	KEY_F9  // Keyboard F9
	KEY_F10 // Keyboard F10
	KEY_F11 // Keyboard F11
	KEY_F12 // Keyboard F12

	KEY_SYSRQ      // Keyboard Print Screen
	KEY_SCROLLLOCK // Keyboard Scroll Lock
	KEY_PAUSE      // Keyboard Pause
	KEY_INSERT     // Keyboard Insert
	KEY_HOME       // Keyboard Home
	KEY_PAGEUP     // Keyboard Page Up
	KEY_DELETE     // Keyboard Delete Forward
	KEY_END        // Keyboard End
	KEY_PAGEDOWN   // Keyboard Page Down
	KEY_RIGHT      // Keyboard Right Arrow
	KEY_LEFT       // Keyboard Left Arrow
	KEY_DOWN       // Keyboard Down Arrow
	KEY_UP         // Keyboard Up Arrow

	KEY_NUMLOCK    // Keyboard Num Lock and Clear
	KEY_KPSLASH    // Keypad /
	KEY_KPASTERISK // Keypad *
	KEY_KPMINUS    // Keypad -
	KEY_KPPLUS     // Keypad +
	KEY_KPENTER    // Keypad ENTER
	KEY_KP1        // Keypad 1 and End
	KEY_KP2        // Keypad 2 and Down Arrow
	KEY_KP3        // Keypad 3 and PageDn
	KEY_KP4        // Keypad 4 and Left Arrow
	KEY_KP5        // Keypad 5
	KEY_KP6        // Keypad 6 and Right Arrow
	KEY_KP7        // Keypad 7 and Home
	KEY_KP8        // Keypad 8 and Up Arrow
	KEY_KP9        // Keypad 9 and Page Up
	KEY_KP0        // Keypad 0 and Insert
	KEY_KPDOT      // Keypad . and Delete

	KEY_102ND   // Keyboard Non-US \ and |
	KEY_COMPOSE // Keyboard Application
	KEY_POWER   // Keyboard Power
	KEY_KPEQUAL // Keypad =

	KEY_F13 // Keyboard F13
	KEY_F14 // Keyboard F14
	KEY_F15 // Keyboard F15
	KEY_F16 // Keyboard F16
	KEY_F17 // Keyboard F17
	KEY_F18 // Keyboard F18
	KEY_F19 // Keyboard F19
	KEY_F20 // Keyboard F20
	KEY_F21 // Keyboard F21
	KEY_F22 // Keyboard F22
	KEY_F23 // Keyboard F23
	KEY_F24 // Keyboard F24

	KEY_OPEN             // Keyboard Execute
	KEY_HELP             // Keyboard Help
	KEY_PROPS            // Keyboard Menu
	KEY_FRONT            // Keyboard Select
	KEY_STOP             // Keyboard Stop
	KEY_AGAIN            // Keyboard Again
	KEY_UNDO             // Keyboard Undo
	KEY_CUT              // Keyboard Cut
	KEY_COPY             // Keyboard Copy
	KEY_PASTE            // Keyboard Paste
	KEY_FIND             // Keyboard Find
	KEY_MUTE             // Keyboard Mute
	KEY_VOLUMEUP         // Keyboard Volume Up
	KEY_VOLUMEDOWN       // Keyboard Volume Down
	_                    //   Keyboard Locking Caps Lock
	_                    //   Keyboard Locking Num Lock
	_                    //   Keyboard Locking Scroll Lock
	KEY_KPCOMMA          // Keypad Comma
	_                    //   Keypad Equal Sign
	KEY_RO               // Keyboard International1
	KEY_KATAKANAHIRAGANA // Keyboard International2
	KEY_YEN              // Keyboard International3
	KEY_HENKAN           // Keyboard International4
	KEY_MUHENKAN         // Keyboard International5
	KEY_KPJPCOMMA        // Keyboard International6
	_                    //   Keyboard International7
	_                    //   Keyboard International8
	_                    //   Keyboard International9
	KEY_HANGEUL          // Keyboard LANG1
	KEY_HANJA            // Keyboard LANG2
	KEY_KATAKANA         // Keyboard LANG3
	KEY_HIRAGANA         // Keyboard LANG4
	KEY_ZENKAKUHANKAKU   // Keyboard LANG5
	_                    //   Keyboard LANG6
	_                    //   Keyboard LANG7
	_                    //   Keyboard LANG8
	_                    //   Keyboard LANG9
	_                    //   Keyboard Alternate Erase
	_                    //   Keyboard SysReq/Attention
	_                    //   Keyboard Cancel
	_                    //   Keyboard Clear
	_                    //   Keyboard Prior
	_                    //   Keyboard Return
	_                    //   Keyboard Separator
	_                    //   Keyboard Out
	_                    //   Keyboard Oper
	_                    //   Keyboard Clear/Again
	_                    //   Keyboard CrSel/Props
	_                    //   Keyboard ExSel

	_                //   Keypad 00
	_                //   Keypad 000
	_                //   Thousands Separator
	_                //   Decimal Separator
	_                //   Currency Unit
	_                //   Currency Sub-unit
	KEY_KPLEFTPAREN  // Keypad (
	KEY_KPRIGHTPAREN // Keypad )
	_                //   Keypad {
	_                //   Keypad }
	_                //   Keypad Tab
	_                //   Keypad Backspace
	_                //   Keypad A
	_                //   Keypad B
	_                //   Keypad C
	_                //   Keypad D
	_                //   Keypad E
	_                //   Keypad F
	_                //   Keypad XOR
	_                //   Keypad ^
	_                //   Keypad %
	_                //   Keypad <
	_                //   Keypad >
	_                //   Keypad &
	_                //   Keypad &&
	_                //   Keypad |
	_                //   Keypad ||
	_                //   Keypad :
	_                //   Keypad #
	_                //   Keypad Space
	_                //   Keypad @
	_                //   Keypad !
	_                //   Keypad Memory Store
	_                //   Keypad Memory Recall
	_                //   Keypad Memory Clear
	_                //   Keypad Memory Add
	_                //   Keypad Memory Subtract
	_                //   Keypad Memory Multiply
	_                //   Keypad Memory Divide
	_                //   Keypad +/-
	_                //   Keypad Clear
	_                //   Keypad Clear Entry
	_                //   Keypad Binary
	_                //   Keypad Octal
	_                //   Keypad Decimal
	_                //   Keypad Hexadecimal

	KEY_LEFTCTRL   // Keyboard Left Control
	KEY_LEFTSHIFT  // Keyboard Left Shift
	KEY_LEFTALT    // Keyboard Left Alt
	KEY_LEFTMETA   // Keyboard Left GUI
	KEY_RIGHTCTRL  // Keyboard Right Control
	KEY_RIGHTSHIFT // Keyboard Right Shift
	KEY_RIGHTALT   // Keyboard Right Alt
	KEY_RIGHTMETA  // Keyboard Right GUI

	KEY_MEDIA_PLAYPAUSE
	KEY_MEDIA_STOPCD
	KEY_MEDIA_PREVIOUSSONG
	KEY_MEDIA_NEXTSONG
	KEY_MEDIA_EJECTCD
	KEY_MEDIA_VOLUMEUP
	KEY_MEDIA_VOLUMEDOWN
	KEY_MEDIA_MUTE
	KEY_MEDIA_WWW
	KEY_MEDIA_BACK
	KEY_MEDIA_FORWARD
	KEY_MEDIA_STOP
	KEY_MEDIA_FIND
	KEY_MEDIA_SCROLLUP
	KEY_MEDIA_SCROLLDOWN
	KEY_MEDIA_EDIT
	KEY_MEDIA_SLEEP
	KEY_MEDIA_COFFEE
	KEY_MEDIA_REFRESH
	KEY_MEDIA_CALC

	// KEY_MUSIC      // Music
	// KEY_MEDIA_STOP // Media Stop
	// KEY_MEDIA_PREV // Media Previous
	// KEY_PAUSE_PLAY // Pause/Play
	// KEY_MEDIA_NEXT // Media Next
	// KEY_MUTE       // Mute
	// KEY_VOLUME_DN  // Volume Up
	// KEY_VOLUME_UP  // Volume Down

	// KEY_MEDIA_CALC // Media Calculator
)
