package sch

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Msg_SchedulerHasStarted = "scheduler has started"
	Msg_SchedulerHasStopped = "scheduler has stopped"
)

type Scheduler struct {
	funcs60  []ScheduledFn
	mustStop atomic.Bool
	wg       *sync.WaitGroup
}

func NewScheduler(funcs60 []ScheduledFn) (s *Scheduler) {
	s = &Scheduler{
		funcs60: funcs60,
		wg:      new(sync.WaitGroup),
	}

	s.mustStop.Store(false)

	return s
}

func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
}

func (s *Scheduler) Stop() {
	s.mustStop.Store(true)
	s.wg.Wait()
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	fmt.Println(Msg_SchedulerHasStarted)

	// Time counter.
	// It counts seconds and resets every 24 hours.
	var tc uint = 1
	const SecondsInDay = 86400 // 60*60*24.
	var err error

	for {
		if s.mustStop.Load() {
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
		}

		// Next tick.
		if tc == SecondsInDay {
			tc = 0
		}
		tc++
		time.Sleep(time.Second)
	}

	fmt.Println(Msg_SchedulerHasStopped)
}

func (s *Scheduler) log(v ...any) {
	log.Println(v...)
}
