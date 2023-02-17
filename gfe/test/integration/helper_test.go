package integration

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	iHttp "github.com/vault-thirteen/junk/gfe/internal/http"
	"github.com/vault-thirteen/junk/gfe/pkg/models/downloadevent"
	"github.com/vault-thirteen/junk/gfe/pkg/models/event"
	"github.com/vault-thirteen/junk/gfe/pkg/models/file"
	"github.com/vault-thirteen/junk/gfe/pkg/models/simpleevent"
	"github.com/vault-thirteen/junk/gfe/pkg/models/user"
)

// composeTopicName создаёт название темы (топика) для Kafka.
func composeTopicName(testId string) (topicName string) {
	return composeId(IntegrationTestIdPrefix, "topic", testId)
}

// composeId создаёт идентификатор из его частей.
func composeId(elements ...string) (id string) {
	return strings.Join(elements, "_")
}

// skipTestIfNotSubtest -- защита от неправильного запуска суб-тестов.
func skipTestIfNotSubtest(t *testing.T) {
	if !isSubTest(t) {
		t.Skip()
	}
}

// isSubTest проверяет, является ли тест суб-тестом.
func isSubTest(t *testing.T) bool {
	testNameParts := strings.Split(t.Name(), `/`)

	if len(testNameParts) == 1 {
		return false
	}

	return true
}

// assertEqualSimpleEvent производит сравнение параметров простого события.
func assertEqualSimpleEvent(
	t *testing.T,
	expectedUserId user.ID,
	expectedFileId file.ID,
	expectedEventType event.TypeID,
	expectedTime time.Time,
	simpleEvent simpleevent.SimpleEvent,
) {
	assert.Equal(t, expectedUserId, simpleEvent.UserID)
	assert.Equal(t, expectedFileId, simpleEvent.FileID)
	assert.Equal(t, expectedEventType, simpleEvent.EventTypeID)
	assertEqualTimeInMs(t, expectedTime, simpleEvent.EventTime)
}

// assertEqualDownloadEvent производит сравнение параметров события типа
// "скачивание".
func assertEqualDownloadEvent(
	t *testing.T,
	expectedUserId user.ID,
	expectedFileId file.ID,
	expectedTime time.Time,
	simpleEvent downloadevent.DownloadEvent,
) {
	assert.Equal(t, expectedUserId, simpleEvent.UserID)
	assert.Equal(t, expectedFileId, simpleEvent.FileID)
	assertEqualTimeInMs(t, expectedTime, simpleEvent.EventTime)
}

// assertEqualTimeInMs производит сравнение времён при округлении их до
// миллисекунды.
func assertEqualTimeInMs(
	t *testing.T,
	t1 time.Time,
	t2 time.Time,
) {
	const Million = 1_000_000 // Nano = 10^(-9), Milli = 10^(-3).

	t1ms := t1.UnixNano() / Million
	t2ms := t2.UnixNano() / Million

	assert.Equal(t, t1ms, t2ms)
}

// makeAuthorizationHeader возвращает название и значение HTTP заголовка для
// авторизации.
func (t *Test) makeAuthorizationHeader() (headerName string, headerValue string) {
	return iHttp.HttpHeaderAuthorization,
		iHttp.HttpHeaderAuthorizationTypeBearer + " " + t.JwtText
}
