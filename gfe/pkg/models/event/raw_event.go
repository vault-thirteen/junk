package event

// RawEvent -- событие по файлу, для выборки со склейками.
// Raw -- "сырые" данные из хранилища.
// Используется для HTTP запросов.
type RawEvent struct {
	// День склейки события по часовому поясу клиента, выполняющего выборку.
	// Для раздельных событий, не являющихся склейкой, равен NULL.
	// См. Примечание ниже.
	LocalDay *string `db:"local_day"`

	// Количество раздельных событий в склейке по указанному дню по часовому
	// поясу клиента, выполняющего выборку. Для раздельных событий, не
	// являющихся склейкой, равно NULL.
	SeparateEventsCount *int `db:"count"`

	// Момент времени (по часовому поясу клиента, выполняющего выборку), в
	// который произошло событие или склейка событий. Временем склейки событий
	// считается время самого позднего события из склейки.
	// См. Примечание ниже.
	LocalTime string `db:"local_time"`

	// Идентификатор типа события.
	EventTypeID TypeID `db:"event_type_id"`

	// Примечание.
	// Go драйвер PostgreSQL читает время без часового пояса
	// как будто оно имеет часовой пояс UTC. Чтобы не запутаться, читаем его
	// сначала в string, затем правим время и распознаём (парсим).
}
