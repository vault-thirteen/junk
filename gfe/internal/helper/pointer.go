package helper

// NewStringPointer создаёт указатель на строку из текста.
// Этот метод в основном нужен для тестирования.
func NewStringPointer(s string) *string {
	ptr := new(string)
	*ptr = s

	return ptr
}

// NewIntPointer создаёт указатель на int из int.
// Этот метод в основном нужен для тестирования.
func NewIntPointer(i int) *int {
	ptr := new(int)
	*ptr = i

	return ptr
}
