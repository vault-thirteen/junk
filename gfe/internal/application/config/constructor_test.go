package config

// В этом файле описаны типы, общие для всех тестов конструктора настроек.

// Общий для тестов префикс названия переменной окружения.
const EnvPrefixCommon = "test"

// TestDataForConfigConstructor -- параметры теста.
type TestDataForConfigConstructor struct {
	// VariablesData -- список настроек переменных окружения.
	// Перед началом теста, мы задаём значения указанным переменным
	// окружения.
	VariablesData []TestEnvVarDataForConfigConstructor

	// Параметр для тестируемого конструктора.
	EnvPrefix string

	// Ожидаемый результат теста.
	ExpectedResult ExpectedResultForConfigConstructor
}

// TestEnvVarDataForConfigConstructor -- объект, хранящий настройки для
// тестирования переменной окружения.
type TestEnvVarDataForConfigConstructor struct {
	// EnvVarName -- название переменной окружения.
	EnvVarName string

	// EnvVarValue -- значение переменной окружения.
	EnvVarValue string
}

// ExpectedResultForConfigConstructor -- ожидаемый результат теста.
type ExpectedResultForConfigConstructor struct {
	// Полученный объект.
	Object interface{}

	// Флаг, показывающий ожидается ли ошибка.
	IsErrorExpected bool
}
