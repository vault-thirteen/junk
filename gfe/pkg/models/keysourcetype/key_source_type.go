package keysourcetype

// Возможные значения типа источника публичного RSA ключа.
const (
	// EnvironmentVariable -- переменная окружения.
	EnvironmentVariable KeySourceType = 1

	// File -- файл.
	File KeySourceType = 2

	// Vault -- Vault.
	Vault KeySourceType = 3
)

// KeySourceType -- тип источника публичного RSA ключа.
// Используется для проверки подписи JWT токена и для других целей.
type KeySourceType byte

// IsValid проверяет является ли тип источника публичного RSA ключа допустимым.
func (kst KeySourceType) IsValid() bool {
	switch kst {
	case EnvironmentVariable,
		File,
		Vault:
		return true
	}

	return false
}

// IsEnvironmentVariable проверяет является ли тип источника публичного RSA
// ключа переменной окружения.
func (kst KeySourceType) IsEnvironmentVariable() bool {
	return kst.isEqualTo(EnvironmentVariable)
}

// IsFile проверяет является ли тип источника публичного RSA ключа файлом.
func (kst KeySourceType) IsFile() bool {
	return kst.isEqualTo(File)
}

// IsVault проверяет является ли тип источника публичного RSA ключа хранилищем
// Vault.
func (kst KeySourceType) IsVault() bool {
	return kst.isEqualTo(Vault)
}

// isEqualTo сравнивает текущий тип с типом, заданным в качестве аргумента.
func (kst KeySourceType) isEqualTo(that KeySourceType) bool {
	return kst == that
}
