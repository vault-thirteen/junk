package event

// Type -- тип события.
type Type struct {
	// Идентификатор типа события.
	ID TypeID `db:"id" json:"id"`

	// Описание на русском языке.
	DescriptionRu string `db:"description_ru" json:"descriptionRu"`

	// Описание на английском языке.
	DescriptionEn string `db:"description_en" json:"descriptionEn"`

	// Является ли событие такого типа простым, т.е. не агрегируемым, т.е.
	// таким, которое не склеивается с другими.
	IsSimple bool `db:"is_simple" json:"isSimple"`

	// Является ли событие такого типа агрегируемым, т.е. таким, которое
	// склеивается с другими.
	IsAggregated bool `db:"is_aggregated" json:"isAggregated"`
}
