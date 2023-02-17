package application

import (
	"time"

	"github.com/vault-thirteen/junk/gfe/internal/message"
)

// Stop останавливает приложение.
func (a *Application) Stop() (err error) {
	// Защита от дурака.
	a.stopperLock.Lock()
	defer a.stopperLock.Unlock()
	if !a.isStarted.Load() {
		return message.ErrIsNotStarted
	}

	err = a.kafka.Close()
	if err != nil {
		return err
	}

	err = a.storage.Close()
	if err != nil {
		return err
	}

	err = a.stopHttpServers()
	if err != nil {
		return err
	}

	// Изменяем состояние.
	a.isStarted.Store(false)

	return nil
}

// stopHttpServers останавливает HTTP серверы.
func (a *Application) stopHttpServers() (err error) {
	const delayAfterHttpServerShutdownMs = 100

	err = a.businessHttpServer.Stop()
	if err != nil {
		return err
	}

	err = a.metricsHttpServer.Stop()
	if err != nil {
		return err
	}

	// Чтобы сообщения о завершении попали в журнал, нужно сделать паузу.
	time.Sleep(time.Millisecond * delayAfterHttpServerShutdownMs)

	return nil
}
