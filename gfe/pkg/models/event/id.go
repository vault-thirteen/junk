package event

import (
	"github.com/pkg/errors"
)

// Значения идентификатора типа события.
const (
	// Creation -- создание.
	Creation TypeID = 1

	// Upload -- загрузка.
	Upload TypeID = 2

	// Download -- скачивание.
	Download TypeID = 3

	// Modification -- изменение.
	Modification TypeID = 4
)

// TypeID -- идентификатор типа события.
type TypeID int16

// IsValid проверяет правильность идентификатора типа события.
func (id TypeID) IsValid() (bool, error) {
	switch id {
	case Creation,
		Upload,
		Download,
		Modification:
		return true, nil
	}

	return false, errors.Errorf(ErrFTypeUnsupported, id)
}

// IsCreation проверяет, является ли тип события типом "Создание".
func (id TypeID) IsCreation() bool {
	return id.isEqualTo(Creation)
}

// IsUpload проверяет, является ли тип события типом "Загрузка".
func (id TypeID) IsUpload() bool {
	return id.isEqualTo(Upload)
}

// IsDownload проверяет, является ли тип события типом "Скачивание".
func (id TypeID) IsDownload() bool {
	return id.isEqualTo(Download)
}

// IsModification проверяет, является ли тип события типом "Изменение".
func (id TypeID) IsModification() bool {
	return id.isEqualTo(Modification)
}

// IsSimple проверяет, является ли тип события типом простого события.
func (id TypeID) IsSimple() bool {
	switch id {
	case Creation,
		Upload,
		Modification:
		return true
	}

	return false
}

// IsAggregated проверяет, является ли тип события типом агрегируемого события.
func (id TypeID) IsAggregated() bool {
	switch id {
	case Download:
		return true
	}

	return false
}

// isEqualTo сравнивает текущий TypeID с TypeID, заданным в качестве аргумента.
func (id TypeID) isEqualTo(that TypeID) bool {
	return id == that
}
