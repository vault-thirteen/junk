package integration

import (
	"os"
	"strings"
	"time"

	"github.com/vault-thirteen/junk/gfe/internal/application"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/message"
)

// Вспомогательный объект для проведения интеграционного тестирования.
// Хранит информацию, которая нужна всем тестам, используется для подготовки и
// завершения процесса тестирования.
type Test struct {
	// UniqueTestId -- уникальный идентификатор теста.
	UniqueTestId string

	// Путь до файла, который содержит публичный RSA ключ для проверки JWT.
	RsaPublicKeyFilePath string

	// JwtText -- текстовое значение JSON веб токена для совершения HTTP
	// запросов к API.
	JwtText string

	// Путь до файла, который содержит текстовое значение JSON веб токена.
	JwtFilePath string

	// EnvironmentVariables -- список используемых переменных окружения.
	EnvironmentVariables []EnvironmentVariable

	// Kafka -- объекты, связанные с тестированием отправки сообщений о
	// событиях.
	Kafka *Kafka

	// Kafka -- объекты, связанные с базой данных.
	Storage *Storage

	// API -- объекты, связанные с тестированием программного интерфейса (API).
	API *API

	// Массив со списком типов генерируемых событий.
	// Значения могут повторяться, т.к. это не список уникальных типов событий,
	// а именно список типов для создаваемых сообщений.
	CreatedEventTypes []event.TypeID

	// Массив со списком времени  генерируемых событий.
	// Таблица времени используется для значительного упрощения проверки
	// результатов, в том числе для упрощения расчёта агрегированных событий.
	CreatedEventTimes []time.Time

	// SentEventMessage -- массив созданных и отправленных сообщений.
	// Используется для проверки ответов сервиса.
	SentEventMessages []*message.Message

	// Тестируемое приложение.
	App *application.Application
}

// NewTest -- конструктор объекта, выполняющего проведение тестов.
func NewTest(
	testId string,
	rsaPublicKeyFilePath string,
	jwtFilePath string,
) (test *Test, err error) {
	test = new(Test)

	test.UniqueTestId = testId
	test.RsaPublicKeyFilePath = rsaPublicKeyFilePath
	test.JwtFilePath = jwtFilePath

	err = test.init()
	if err != nil {
		return nil, err
	}

	return test, nil
}

// init -- инициализация объекта, выполняющего проведение тестов.
func (t *Test) init() (err error) {
	err = t.initCreatedEvents()
	if err != nil {
		return err
	}

	err = t.readJwt()
	if err != nil {
		return err
	}

	t.EnvironmentVariables, err = t.prepareEnvironmentVariablesList()
	if err != nil {
		return err
	}

	err = t.setEnvVars(t.EnvironmentVariables)
	if err != nil {
		return err
	}

	t.Kafka, err = NewKafka(t, t.UniqueTestId)
	if err != nil {
		return err
	}

	t.Storage, err = NewStorage()
	if err != nil {
		return err
	}

	t.API, err = NewAPI(t)
	if err != nil {
		return err
	}

	t.App, err = application.NewApplication()
	if err != nil {
		return err
	}

	err = t.App.Start()
	if err != nil {
		return err
	}

	return nil
}

// initCreatedEvents подготавливает параметры для создаваемых сообщений.
func (t *Test) initCreatedEvents() (err error) {
	t.CreatedEventTypes = []event.TypeID{
		event.Creation,
		event.Upload,
		event.Download,
		event.Modification,
		event.Download,
	}

	var (
		locationGreenwich *time.Location // UTC + 00:00, no DST.
		locationMoscow    *time.Location // UTC + 03:00, no DST.
		locationHawaii    *time.Location // UTC - 10:00, no DST.
	)

	locationGreenwich, err = time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	locationMoscow, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	locationHawaii, err = time.LoadLocation("Pacific/Honolulu")
	if err != nil {
		return err
	}

	// В 2021-ом году равноденствие было в 20-ый день марта в 09:37:27 по
	// Гринвичу. Проверяющий клиент будет запрашивать время по часовому поясу
	// 'Йоханесбург' (UTC+2). По часовому поясу клиента, одна часть сообщений
	// уложится до полуночи, другая часть сообщений придёт на следующий день.
	t.CreatedEventTimes = []time.Time{
		// Равноденствие.
		// По Гринвичу = 09:37, в Йоханесбурге = 11:37.
		time.Date(2021, 3, 20, 9, 37, 27, 0, locationGreenwich),

		// Через 11 часов после равноденствия.
		// По Гринвичу = 20:37, в Йоханесбурге = 22:37.
		time.Date(2021, 3, 20, 20, 37, 27, 0, locationGreenwich),

		// Через 12 часов после равноденствия.
		// По Гринвичу = 21:37, в Йоханесбурге = 23:37.
		// Для проверки сервиса, это сообщение приходит по часовому поясу Москвы (UTC+3).
		time.Date(2021, 3, 20, 21+3, 37, 27, 0, locationMoscow),

		// Через 13 часов после равноденствия.
		// По Гринвичу = 22:37, в Йоханесбурге = 00:37 следующего дня.
		// Для проверки сервиса, это сообщение приходит по часовому поясу Москвы (UTC+3).
		time.Date(2021, 3, 20, 22+3, 37, 27, 0, locationMoscow),

		// Через 14 часов после равноденствия.
		// По Гринвичу = 23:37, в Йоханесбурге = 01:37 следующего дня.
		// Для проверки сервиса, это сообщение приходит по часовому поясу Гавайи (UTC-10).
		time.Date(2021, 3, 20, 23-10, 37, 27, 0, locationHawaii),
	}

	// В итоге, для клиента из Йоханесбурга, первые три сообщения будут за 20 марта,
	// следующие два события будут за 21 марта.

	return nil
}

// readJwt читает текстовое значение JSON веб токена из файла.
func (t *Test) readJwt() (err error) {
	var buffer []byte
	buffer, err = os.ReadFile(t.JwtFilePath)
	if err != nil {
		return err
	}

	t.JwtText = strings.TrimSpace(string(buffer))

	return nil
}

// fin -- остановка объекта, выполняющего проведение тестов.
func (t *Test) fin() (err error) {
	err = t.Kafka.Stop()
	if err != nil {
		return err
	}

	err = t.Storage.Stop()
	if err != nil {
		return err
	}

	err = t.App.Stop()
	if err != nil {
		return err
	}

	err = t.unsetEnvVars()
	if err != nil {
		return err
	}

	return nil
}

// Stop останавливает тест.
func (t *Test) Stop() (err error) {
	err = t.fin()
	if err != nil {
		return err
	}

	return nil
}
