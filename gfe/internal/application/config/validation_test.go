package config

// В этом файле описаны типы, общие для всех тестов валидации объекта с
// настройками (конфига).

// TestDataForIsValid -- параметры теста валидации объекта с настройками.
type TestDataForIsValid struct {
	// Тестируемый объект с настройками.
	TestedConfig interface{}

	// Ожидаемый результат теста.
	ExpectedResult ExpectedResultForIsValid
}

// ExpectedResultForIsValid -- ожидаемый результат теста.
type ExpectedResultForIsValid struct {
	// Флаг, показывающий годность конфига.
	IsValid bool

	// Флаг, показывающий ожидается ли ошибка.
	IsErrorExpected bool
}
