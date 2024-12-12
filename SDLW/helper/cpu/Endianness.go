package cpu

const (
	EndiannessUnknown      = Endianness(0)
	EndiannessBigEndian    = Endianness(1)
	EndiannessLittleEndian = Endianness(2)
)

type Endianness byte
