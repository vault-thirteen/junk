package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vault-thirteen/junk/gfe/internal/helper"
)

// TestData -- тестовые данные.
type TestData struct {
	// Input -- входные параметры.
	Input TestDataInput

	// Output -- выходные параметры.
	Output TestDataOutput
}

// TestDataInput -- входные параметры теста.
type TestDataInput struct {
	// RawData -- событие по файлу, "сырые" данные.
	RawData *RawEvent

	// TimeZoneName -- название часового пояса.
	TimeZoneName string
}

// Выходные параметры теста.
type TestDataOutput struct {
	// ExpectedEvent -- событие.
	ExpectedEvent *Event

	// ExpectedEventLocalTimeString -- локальное время события в виде строки.
	// Дополнительная проверка для надёжности.
	ExpectedEventLocalTimeString string

	// IsErrorExpected -- ожидается ли в тесте ошибка.
	IsErrorExpected bool
}

func TestNewFromRawData(t *testing.T) {
	// Arrange.
	tests := make([]TestData, 0)

	var (
		separateEventsCount *int
		localTime           time.Time
		location            *time.Location
		err                 error
	)

	// Тест 1.
	// Правильные данные.
	separateEventsCount = helper.NewIntPointer(123)
	location, err = time.LoadLocation("Europe/Moscow")
	assert.NoError(t, err)
	localTime = time.Date(2000, 1, 1, 23, 55, 1, 0, location)
	tests = append(tests, TestData{
		Input: TestDataInput{
			RawData: &RawEvent{
				LocalDay:            helper.NewStringPointer("2000-01-01..."),
				SeparateEventsCount: separateEventsCount,
				LocalTime:           "2000-01-01T23:55:01Z",
				EventTypeID:         Creation,
			},
			TimeZoneName: "Europe/Moscow",
		},
		Output: TestDataOutput{
			ExpectedEvent: &Event{
				LocalDay:            helper.NewStringPointer("2000-01-01"),
				SeparateEventsCount: separateEventsCount,
				LocalTime:           localTime,
				EventTypeID:         Creation,
			},
			ExpectedEventLocalTimeString: "2000-01-01 23:55:01 +0300 MSK",
			IsErrorExpected:              false,
		},
	})

	// Тест 2-A.
	// Кривые данные.
	// Сырые данные не заданы.
	tests = append(tests, TestData{
		Input: TestDataInput{
			RawData:      nil,
			TimeZoneName: "Europe/Moscow",
		},
		Output: TestDataOutput{
			ExpectedEvent:   nil,
			IsErrorExpected: true,
		},
	})

	// Тест 2-B.
	// Кривые данные.
	// Часовой пояс задан неправильно.
	tests = append(tests, TestData{
		Input: TestDataInput{
			RawData: &RawEvent{
				LocalDay:            helper.NewStringPointer("2000-01-01..."),
				SeparateEventsCount: separateEventsCount,
				LocalTime:           "2000-01-01T23:55:01Z",
				EventTypeID:         Creation,
			},
			TimeZoneName: "Quake III",
		},
		Output: TestDataOutput{
			ExpectedEvent:   nil,
			IsErrorExpected: true,
		},
	})

	// Тест 2-C.
	// Кривые данные.
	// Местное время задано неправильно.
	tests = append(tests, TestData{
		Input: TestDataInput{
			RawData: &RawEvent{
				LocalDay:            helper.NewStringPointer("2000-01-01..."),
				SeparateEventsCount: separateEventsCount,
				LocalTime:           "2000-XX-01T23:55:01Z",
				EventTypeID:         Creation,
			},
			TimeZoneName: "Europe/Moscow",
		},
		Output: TestDataOutput{
			ExpectedEvent:   nil,
			IsErrorExpected: true,
		},
	})

	// Запуск тестов.
	for i, test := range tests {
		fmt.Printf("[%d] ", i+1)

		// Act. Запускаем тестируемый метод или функцию.
		eventActual, errActual := NewFromRawData(test.Input.RawData, test.Input.TimeZoneName)

		// Assert.1. Проверяем ошибку.
		switch test.Output.IsErrorExpected {
		case true:
			assert.Error(t, errActual)
		case false:
			assert.NoError(t, errActual)
		}

		// Assert.2. Проверяем объект.
		assert.Equal(t, test.Output.ExpectedEvent, eventActual)

		// Assert.3. Проверяем дополнительную строку локального времени.
		if eventActual != nil {
			assert.Equal(t, test.Output.ExpectedEventLocalTimeString, eventActual.LocalTime.String())
		}
	}
}
