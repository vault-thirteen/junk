package shc

import (
	"fmt"

	"github.com/vault-thirteen/TR1/src/interfaces"
	"github.com/vault-thirteen/TR1/src/libraries/scheduler"
)

type SchedulerComponent struct {
	cfg       interfaces.IConfiguration
	functions []sch.ScheduledFn
	scheduler *sch.Scheduler
}

func (c *SchedulerComponent) Init(cfg interfaces.IConfiguration, controller interfaces.IController) (sc interfaces.IServiceComponent, err error) {
	shc := &SchedulerComponent{
		cfg:       cfg,
		functions: controller.GetScheduledFunctions(),
	}

	shc.scheduler = sch.NewScheduler(shc.functions)

	return shc, nil
}
func (c *SchedulerComponent) GetConfiguration() interfaces.IConfiguration {
	return c.cfg
}

func (c *SchedulerComponent) Start(s interfaces.IService) (err error) {
	c.scheduler.Start()
	return nil
}
func (c *SchedulerComponent) Stop(s interfaces.IService) (err error) {
	wg := s.GetSubRoutinesWG()
	defer wg.Done()

	c.scheduler.Stop()
	c.ReportStop()

	return nil
}

func (c *SchedulerComponent) ReportStart() {
	fmt.Println("SchedulerComponent has started")
}
func (c *SchedulerComponent) ReportStop() {
	fmt.Println("SchedulerComponent has stopped")
}

// Other methods.

func FromAny(x any) (c *SchedulerComponent) {
	return x.(*SchedulerComponent)
}

// Non-standard methods.
