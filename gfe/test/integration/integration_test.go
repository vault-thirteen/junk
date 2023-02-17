package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vault-thirteen/junk/gfe/internal/application/config"
	"github.com/vault-thirteen/junk/gfe/internal/random"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
	"github.com/vault-thirteen/junk/gfe/pkg/models/user"
)

const (
	// IntegrationTestIdPrefix -- префикс идентификатора текущего экземпляра
	// интеграционного теста.
	IntegrationTestIdPrefix = "integration_test"

	// UserId -- используемое в тесте значение идентификатора пользователя.
	UserId user.ID = "User-111"

	// FileId -- используемое в тесте значение идентификатора файла.
	FileId file.ID = "File-222"
)

// TestFunc -- тип (сигнатура) функции теста.
type TestFunc = func(t *testing.T)

// Глобальный объект, хранящий параметры тестов.
// Используется для передачи информации между разными тестами, поскольку
// нельзя менять сигнатуру функции тестирования.
var test *Test

// Общий тест, запускающий остальные тесты как суб-тесты.
// Для защиты от неправильного запуска, в каждом суб-тесте есть проверка
// способа его запуска. Если суб-тест запущен не как суб-тест, он просто
// пропускается, а когда он запущен правильно, как суб-тест, он проходит.
// Важно, что для проведения интеграционных тестов нужно заранее, перед
// проведением тестов, сгенерировать публичный RSA ключ и JWT токен. Скрипт,
// осуществляющий это генерирование, входит в исходный код сервиса.
func TestAll(t *testing.T) {
	var err error

	// Arrange.
	testId := random.MakeUniqueRandomString()
	rsaPublicKeyFilePath := `./../web-token/key/jwt_rsa_RS256.key.pub`
	jwtFilePath := `./../web-token/key/jwt_token.txt`

	test, err = NewTest(testId, rsaPublicKeyFilePath, jwtFilePath)
	assert.NoError(t, err)

	defer func() {
		// Finalization.
		err = test.Stop()
		assert.NoError(t, err)
	}()

	testFunctions := []TestFunc{
		Test_ProduceMessages,
		Test_ReadStorage,
		Test_UseAPI,
	}

	// Act & assert.
	for i, testFunction := range testFunctions {
		n := i + 1
		t.Run(strconv.Itoa(n), testFunction)

		time.Sleep(time.Millisecond * 0) // Debug.
	}
}

// Суб-тест, отправляющий сообщения в Kafka.
func Test_ProduceMessages(t *testing.T) {
	skipTestIfNotSubtest(t)

	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Printf("recovery: %v\r\n", exception)
		}
	}()

	var err error
	err = test.Kafka.ProduceMessages()
	assert.NoError(t, err)
	assert.Equal(t, 5, len(test.CreatedEventTypes))
	assert.Equal(t, 5, len(test.SentEventMessages))
}

// Суб-тест, читающий данные из хранилища.
func Test_ReadStorage(t *testing.T) {
	skipTestIfNotSubtest(t)

	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Printf("recovery: %v\r\n", exception)
		}
	}()

	var err error
	var storageData *StorageData
	storageData, err = test.Storage.ReadData()
	assert.NoError(t, err)

	// 1. Грубые проверки.
	assert.NotEqual(t, nil, storageData)
	assert.Equal(t, 4, len(storageData.EventTypes))
	assert.Equal(t, 3, len(storageData.SimpleEvents))
	assert.Equal(t, 2, len(storageData.DownloadEvents))
	assert.Equal(t, 5, len(test.SentEventMessages))

	// 2. Точные проверки.
	// Примечание.
	// Время в базе данных и в языке Go хранятся в разных форматах.
	// В данный момент, СУБД PostgreSQL v13 округляет время до 1 мкс.
	// Для надёжности и некоторого "запаса прочности", сравнение времени
	// производим, округляя его до 1 мс.

	// 2.1. Типы событий.
	expectedEventType := event.Type{
		ID:            1,
		DescriptionRu: "Создание",
		DescriptionEn: "Creation",
		IsSimple:      true,
		IsAggregated:  false,
	}
	assert.Equal(t, expectedEventType, storageData.EventTypes[0])

	expectedEventType = event.Type{
		ID:            2,
		DescriptionRu: "Загрузка",
		DescriptionEn: "Upload",
		IsSimple:      true,
		IsAggregated:  false,
	}
	assert.Equal(t, expectedEventType, storageData.EventTypes[1])

	expectedEventType = event.Type{
		ID:            3,
		DescriptionRu: "Скачивание",
		DescriptionEn: "Download",
		IsSimple:      false,
		IsAggregated:  true,
	}
	assert.Equal(t, expectedEventType, storageData.EventTypes[2])

	expectedEventType = event.Type{
		ID:            4,
		DescriptionRu: "Изменение",
		DescriptionEn: "Modification",
		IsSimple:      true,
		IsAggregated:  false,
	}
	assert.Equal(t, expectedEventType, storageData.EventTypes[3])

	// 2.2. Простые события.
	// В списке всех созданных событий, простые события имеют номера: 1, 2, 4.
	assertEqualSimpleEvent(t, UserId, FileId, event.Creation, test.SentEventMessages[0].EventTime, storageData.SimpleEvents[0])     // #1.
	assertEqualSimpleEvent(t, UserId, FileId, event.Upload, test.SentEventMessages[1].EventTime, storageData.SimpleEvents[1])       // #2.
	assertEqualSimpleEvent(t, UserId, FileId, event.Modification, test.SentEventMessages[3].EventTime, storageData.SimpleEvents[2]) // #4.

	// 2.3. События типа "скачивание".
	// В списке всех созданных событий, события типа "скачивание" имеют номера: 3, 5.
	assertEqualDownloadEvent(t, UserId, FileId, test.SentEventMessages[2].EventTime, storageData.DownloadEvents[0]) // #3.
	assertEqualDownloadEvent(t, UserId, FileId, test.SentEventMessages[4].EventTime, storageData.DownloadEvents[1]) // #5.
}

// Суб-тест, проверяющий правильность работы основного API сервиса.
func Test_UseAPI(t *testing.T) {
	skipTestIfNotSubtest(t)

	defer func() {
		exception := recover()
		if exception != nil {
			fmt.Printf("recovery: %v\r\n", exception)
		}
	}()

	var err error
	var apiData *APIData
	apiData, err = test.API.ReadData()
	assert.NoError(t, err)
	assert.NotEqual(t, nil, apiData)

	// Проверка Readiness.
	assert.Equal(t, http.StatusOK, apiData.ReadinessStatus)

	// Проверка метрик.
	assert.Equal(t, http.StatusOK, apiData.MetricsStatus)

	// Проверка доступности сервиса.
	assert.Equal(t, http.StatusOK, apiData.LivenessStatus)

	// Проверка типов событий.
	assert.NotEqual(t, nil, apiData.FileEventTypes)
	assert.Equal(t, 4, len(apiData.FileEventTypes.Data))
	expectedEventTypes := []event.Type{
		{
			ID:            1,
			DescriptionRu: "Создание",
			DescriptionEn: "Creation",
			IsSimple:      true,
			IsAggregated:  false,
		},
		{
			ID:            2,
			DescriptionRu: "Загрузка",
			DescriptionEn: "Upload",
			IsSimple:      true,
			IsAggregated:  false,
		},
		{
			ID:            3,
			DescriptionRu: "Скачивание",
			DescriptionEn: "Download",
			IsSimple:      false,
			IsAggregated:  true,
		},
		{
			ID:            4,
			DescriptionRu: "Изменение",
			DescriptionEn: "Modification",
			IsSimple:      true,
			IsAggregated:  false,
		},
	}
	assert.Equal(t, expectedEventTypes, apiData.FileEventTypes.Data)

	// Проверка всех событий.
	// Этот хендлер отдаёт события по убыванию, сначала самые новые, поэтому
	// порядок индексов при переборе массива -- обратный (4, 3, 2, 1, 0).
	assert.Equal(t, FileId, apiData.FileEventsAll.Data.FileID)
	assert.Equal(t, TimeZoneForTests, apiData.FileEventsAll.Data.TimeZoneName)
	// Для данного запроса агрегированные события расположены в разных днях,
	// поэтому количество событий в ответе совпадает с количеством записей в
	// базе данных.
	assert.Equal(t, len(test.SentEventMessages), len(apiData.FileEventsAll.Data.Records))

	assert.Equal(t, event.Creation, apiData.FileEventsAll.Data.Records[4].EventTypeID)
	assert.Equal(t, (*string)(nil), apiData.FileEventsAll.Data.Records[4].LocalDay)
	assert.Equal(t, (*int)(nil), apiData.FileEventsAll.Data.Records[4].SeparateEventsCount)
	assert.Equal(t, "2021-03-20 11:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[4].LocalTime.String())

	assert.Equal(t, event.Upload, apiData.FileEventsAll.Data.Records[3].EventTypeID)
	assert.Equal(t, (*string)(nil), apiData.FileEventsAll.Data.Records[3].LocalDay)
	assert.Equal(t, (*int)(nil), apiData.FileEventsAll.Data.Records[3].SeparateEventsCount)
	assert.Equal(t, "2021-03-20 22:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[3].LocalTime.String())

	assert.Equal(t, event.Download, apiData.FileEventsAll.Data.Records[2].EventTypeID)
	assert.Equal(t, "2021-03-20", *apiData.FileEventsAll.Data.Records[2].LocalDay)
	assert.Equal(t, 1, *apiData.FileEventsAll.Data.Records[2].SeparateEventsCount)
	assert.Equal(t, "2021-03-20 23:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[2].LocalTime.String())

	assert.Equal(t, event.Modification, apiData.FileEventsAll.Data.Records[1].EventTypeID)
	assert.Equal(t, (*string)(nil), apiData.FileEventsAll.Data.Records[1].LocalDay)
	assert.Equal(t, (*int)(nil), apiData.FileEventsAll.Data.Records[1].SeparateEventsCount)
	assert.Equal(t, "2021-03-21 00:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[1].LocalTime.String())

	assert.Equal(t, event.Download, apiData.FileEventsAll.Data.Records[0].EventTypeID)
	assert.Equal(t, "2021-03-21", *apiData.FileEventsAll.Data.Records[0].LocalDay)
	assert.Equal(t, 1, *apiData.FileEventsAll.Data.Records[0].SeparateEventsCount)
	assert.Equal(t, "2021-03-21 01:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[0].LocalTime.String())

	// Проверка недавних событий.
	// Этот хендлер отдаёт события по убыванию, сначала самые новые, поэтому
	// порядок индексов при переборе массива -- обратный (4, 3, 2, 1, 0).
	assert.Equal(t, FileId, apiData.FileEventsAll.Data.FileID)
	assert.Equal(t, TimeZoneForTests, apiData.FileEventsAll.Data.TimeZoneName)
	assert.Equal(t, config.FileLastEventsCount, len(apiData.FileEventsLastN.Data.Records))

	assert.Equal(t, event.Upload, apiData.FileEventsAll.Data.Records[3].EventTypeID)
	assert.Equal(t, (*string)(nil), apiData.FileEventsAll.Data.Records[3].LocalDay)
	assert.Equal(t, (*int)(nil), apiData.FileEventsAll.Data.Records[3].SeparateEventsCount)
	assert.Equal(t, "2021-03-20 22:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[3].LocalTime.String())

	assert.Equal(t, event.Download, apiData.FileEventsAll.Data.Records[2].EventTypeID)
	assert.Equal(t, "2021-03-20", *apiData.FileEventsAll.Data.Records[2].LocalDay)
	assert.Equal(t, 1, *apiData.FileEventsAll.Data.Records[2].SeparateEventsCount)
	assert.Equal(t, "2021-03-20 23:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[2].LocalTime.String())

	assert.Equal(t, event.Modification, apiData.FileEventsAll.Data.Records[1].EventTypeID)
	assert.Equal(t, (*string)(nil), apiData.FileEventsAll.Data.Records[1].LocalDay)
	assert.Equal(t, (*int)(nil), apiData.FileEventsAll.Data.Records[1].SeparateEventsCount)
	assert.Equal(t, "2021-03-21 00:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[1].LocalTime.String())

	assert.Equal(t, event.Download, apiData.FileEventsAll.Data.Records[0].EventTypeID)
	assert.Equal(t, "2021-03-21", *apiData.FileEventsAll.Data.Records[0].LocalDay)
	assert.Equal(t, 1, *apiData.FileEventsAll.Data.Records[0].SeparateEventsCount)
	assert.Equal(t, "2021-03-21 01:37:27 +0200 +0200", apiData.FileEventsAll.Data.Records[0].LocalTime.String())
}
