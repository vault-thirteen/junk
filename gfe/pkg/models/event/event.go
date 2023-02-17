package event

import (
	"time"

	"github.com/pkg/errors"
)

const ErrRawDataNotSet = "raw data is not set"

// Event -- событие по файлу, для выборки со склейками.
// Используется для HTTP запросов.
type Event struct {
	// День склейки события по часовому поясу клиента, выполняющего выборку.
	// Для раздельных событий, не являющихся склейкой, равен NULL.
	// См. Примечание ниже.
	LocalDay *string `db:"local_day" json:"day"`

	// Количество раздельных событий в склейке по указанному дню по часовому
	// поясу клиента, выполняющего выборку. Для раздельных событий, не
	// являющихся склейкой, равно NULL.
	SeparateEventsCount *int `db:"count" json:"subEventsCount"`

	// Момент времени (по часовому поясу клиента, выполняющего выборку), в
	// который произошло событие или склейка событий. Временем склейки событий
	// считается время самого позднего события из склейки.
	LocalTime time.Time `db:"local_time" json:"time"`

	// Идентификатор типа события.
	EventTypeID TypeID `db:"event_type_id" json:"eventTypeId"`

	// Примечание.
	// Тип string используется для быстродействия. Поскольку сервис отдаёт
	// ответ в текстовом виде и читает данные из хранилища также в текстовом
	// виде, то в тех местах, где не нужны никакие преобразования,  лишние
	// действия (парсинг и форматирование) не используются.
}

// NewFromRawData -- конструктор типа Event из типа RawEvent.
func NewFromRawData(rawData *RawEvent, timeZoneName string) (e *Event, err error) {
	const (
		// DateStringLength -- длина строки даты.
		DateStringLength = 10 // '2001-01-01'.

		// TimeStringLength -- длина строки времени.
		TimeStringLength = 19 // '2001-01-01T01:00:00'.

		// TimeLayout -- разметка времени.
		TimeLayout = "2006-01-02T15:04:05"
	)

	if rawData == nil {
		return nil, errors.New(ErrRawDataNotSet)
	}

	e = new(Event)

	// 1. День.
	if rawData.LocalDay != nil {
		e.LocalDay = new(string)
		*e.LocalDay = (*rawData.LocalDay)[:DateStringLength]
	}

	// 2. Количество событий в склейке.
	e.SeparateEventsCount = rawData.SeparateEventsCount

	// 3. Время.
	var location *time.Location
	location, err = time.LoadLocation(timeZoneName)
	if err != nil {
		return nil, err
	}

	// Из-за особенностей работы SQL драйвера в Go, время без часового пояса
	// приходит "с" часовым поясом UTC, но на самом деле оно не UTC.
	// Парсим время в другом часовом поясе, отбрасывая часовой пояс исходного
	// времени.
	e.LocalTime, err = time.ParseInLocation(TimeLayout, (rawData.LocalTime)[:TimeStringLength], location)
	if err != nil {
		return nil, err
	}

	// 4. Тип события.
	e.EventTypeID = rawData.EventTypeID

	return e, nil
}
