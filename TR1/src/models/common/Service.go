package cm

import (
	"sync"
	"sync/atomic"

	"github.com/vault-thirteen/TR1/src/interfaces"
)

type Service struct {
	cfg        *Configuration
	iC         *ServiceInternalControls
	eC         *ServiceExternalControls
	scs        []interfaces.IServiceComponent
	controller interfaces.IController
}

func NewService(cfg *Configuration, components []interfaces.IServiceComponent, controller interfaces.IController) (s *Service, err error) {
	s = &Service{
		cfg:        cfg,
		iC:         NewInternalControls(),
		eC:         NewExternalControls(),
		scs:        []interfaces.IServiceComponent{},
		controller: controller,
	}

	var c interfaces.IServiceComponent
	for _, component := range components {
		c, err = component.Init(cfg, controller)
		if err != nil {
			return nil, err
		}

		s.scs = append(s.scs, c)
	}

	return s, nil
}

func (s *Service) Start() (err error) {
	s.iC.ssp.Lock()
	defer s.iC.ssp.Unlock()

	err = s.iC.ssp.BeginStart()
	if err != nil {
		return err
	}

	// Components are started in a normal (forward) order.
	// ErrorListenerComponent must be the first one in the list if it is used.
	for _, sc := range s.scs {
		err = sc.Start(s)
		if err != nil {
			return err
		}

		s.iC.subRoutines.Add(1)
	}

	s.iC.ssp.CompleteStart()

	return nil
}

func (s *Service) GetStopChannel() *chan bool {
	return &s.eC.mustBeStopped
}

func (s *Service) ReportStart() {
	for _, sc := range s.scs {
		sc.ReportStart()
	}
}

func (s *Service) Stop() (err error) {
	s.iC.ssp.Lock()
	defer s.iC.ssp.Unlock()

	err = s.iC.ssp.BeginStop()
	if err != nil {
		return err
	}

	s.iC.mustStop.Store(true)

	// Components are stopped in a reversed (backward) order.
	// This is required for all the errors to be collected.
	for i := len(s.scs) - 1; i >= 0; i-- {
		err = s.scs[i].Stop(s)
		if err != nil {
			return err
		}
	}

	s.iC.subRoutines.Wait()

	s.iC.ssp.CompleteStop()

	return nil
}

func (s *Service) GetSubRoutinesWG() *sync.WaitGroup {
	return s.iC.subRoutines
}

func (s *Service) GetMustStopAB() *atomic.Bool {
	return s.iC.mustStop
}

func (s *Service) GetComponentByIndex(idx int) interfaces.IServiceComponent {
	return s.scs[idx]
}

func (s *Service) GetConfiguration() interfaces.IConfiguration {
	return s.cfg
}
