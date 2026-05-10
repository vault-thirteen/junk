package elc

import (
	"fmt"
	"log"
	"sync"

	"github.com/vault-thirteen/TR1/src/interfaces"
)

const (
	Msg_ServiceError = "Service error: "
)

type ErrorListenerComponent struct {
	cfg           interfaces.IConfiguration
	errorsChan    *chan error
	mustBeStopped *chan bool
	listenerWG    *sync.WaitGroup
}

func (c *ErrorListenerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	elc := &ErrorListenerComponent{
		cfg:        cfg,
		errorsChan: controller.GetErrorsChan(),
		listenerWG: new(sync.WaitGroup),
	}

	return elc, nil
}
func (c *ErrorListenerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *ErrorListenerComponent) Start(s interfaces.IService) (err error) {
	c.mustBeStopped = s.GetStopChannel()

	c.listenerWG.Add(1)
	go c.listenForHttpErrors()

	return nil
}
func (c *ErrorListenerComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	close(*c.errorsChan)
	c.listenerWG.Wait()
	c.ReportStop()

	return nil
}

func (c *ErrorListenerComponent) listenForHttpErrors() {
	defer c.listenerWG.Done()

	go func() {
		for svcErr := range *c.errorsChan {
			log.Println(Msg_ServiceError + svcErr.Error())
			*c.mustBeStopped <- true
		}
	}()
}

func (c *ErrorListenerComponent) ReportStart() {
	fmt.Println("ErrorListenerComponent has started")
}
func (c *ErrorListenerComponent) ReportStop() {
	fmt.Println("ErrorListenerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *ErrorListenerComponent) {
	return x.(*ErrorListenerComponent)
}

// Non-standard methods.
