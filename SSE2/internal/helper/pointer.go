package helper

func NewStringPointer(s string) *string {
	ptr := new(string)
	*ptr = s

	return ptr
}

func NewIntPointer(i int) *int {
	ptr := new(int)
	*ptr = i

	return ptr
}

func NewBoolPointer(b bool) *bool {
	ptr := new(bool)
	*ptr = b

	return ptr
}
