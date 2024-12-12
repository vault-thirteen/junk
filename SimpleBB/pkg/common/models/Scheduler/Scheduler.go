package sch

import (
	cmi "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base1"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	"log"
	"time"
)

type Scheduler struct {
	srv       cmi.IServer
	funcs60   []simple.ScheduledFn
	funcs600  []simple.ScheduledFn
	funcs3600 []simple.ScheduledFn
}

func NewScheduler(srv cmi.IServer, funcs60 []simple.ScheduledFn, funcs600 []simple.ScheduledFn, funcs3600 []simple.ScheduledFn) (s *Scheduler) {
	s = &Scheduler{
		srv:       srv,
		funcs60:   funcs60,
		funcs600:  funcs600,
		funcs3600: funcs3600,
	}

	return s
}

func (s *Scheduler) Run() {
	subRoutinesWG := s.srv.GetSubRoutinesWG()
	defer subRoutinesWG.Done()

	// Time counter.
	// It counts seconds and resets every 24 hours.
	var tc uint = 1
	const SecondsInDay = 86400 // 60*60*24.
	var err error

	for {
		if s.srv.GetMustStopAB().Load() {
			break
		}

		// Periodical tasks (every minute).
		if tc%60 == 0 {
			for _, fn := range s.funcs60 {
				err = fn()
				if err != nil {
					s.log(err)
				}
			}

			// Periodical tasks (every 10 minutes).
			if tc%600 == 0 {
				for _, fn := range s.funcs600 {
					err = fn()
					if err != nil {
						s.log(err)
					}
				}

				// Periodical tasks (every hour).
				if tc%3600 == 0 {
					for _, fn := range s.funcs3600 {
						err = fn()
						if err != nil {
							s.log(err)
						}
					}
				}
			}
		}

		// Next tick.
		if tc == SecondsInDay {
			tc = 0
		}
		tc++
		time.Sleep(time.Second)
	}

	s.log(c.MsgSchedulerHasStopped)
}

func (s *Scheduler) log(v ...any) {
	log.Println(v...)
}
