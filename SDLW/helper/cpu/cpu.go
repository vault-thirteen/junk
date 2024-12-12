package cpu

import (
	"unsafe"
)

// HardwareEndianness tries to detect the "endianness" of the CPU.
func HardwareEndianness() Endianness {
	var buf [2]byte
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return EndiannessLittleEndian
	case [2]byte{0xAB, 0xCD}:
		return EndiannessBigEndian
	default:
		return EndiannessUnknown
	}
}
