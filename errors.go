package Falcon8

import "errors"

var (
	ErrDeviceAlreadyOpen = errors.New("falcon8: device already open")
	ErrDeviceNotOpen     = errors.New("falcon8: device not open")

	ErrClaimingInterface = errors.New("falcon8: error claiming interface %s")
	ErrClaimingConfig    = errors.New("falcon-8: failed to claim config interface")

	ErrClosingConfig  = errors.New("falcon-8: error closing config %s")
	ErrClosingDevice  = errors.New("falcon-8: error closing device %s")
	ErrClosingContext = errors.New("falcon-8: error closing context %s")

	ErrBadCFGRead         = errors.New("falcon-8: failed to read enough data bytes from file")
	ErrBadCFGChecksumRead = errors.New("falcon-8: failed to read enough checksum bytes from file")
	ErrBadCFGChecksum     = errors.New("falcon-8: checksum mismatch")

	ErrInvalidByteArrayLength = errors.New("falcon-8: invalid byte array length, must be 264 for Falcon8")
	ErrInvalidLayer           = errors.New("falcon8: invalid layer")
)
