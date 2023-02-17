package generator

import (
	"errors"
	"log"
	"sync"
	"time"

	rnd "github.com/vault-thirteen/auxie/random"
	"go.uber.org/atomic"
)

const (
	ErrAlreadyStarted = "already started"
	ErrNotStarted     = "not started"
)

type Generator struct {
	componentsCount int

	// Minimum & maximum values for delay before next number generation.
	// These values are created randomly within the range specified in the
	// constructor.
	delayBeforeNextGenerationMin time.Duration
	delayBeforeNextGenerationMax time.Duration

	// Minimum & maximum values for generated random numbers.
	// These values are created randomly within the range specified in the
	// constructor.
	generatedNumberMin uint
	generatedNumberMax uint

	// Channel for generated random numbers.
	output chan uint

	lock        sync.Mutex
	isWorking   atomic.Bool
	mustClose   chan bool
	generatorWG *sync.WaitGroup
}

func New(
	// Range within which a random delay will be created, milliseconds.
	delayBeforeNextGenerationMinRangeBorderMs uint,
	delayBeforeNextGenerationMaxRangeBorderMs uint,

	// Range within which random numbers will be generated.
	generatedNumberMinRangeBorder uint,
	generatedNumberMaxRangeBorder uint,

	outputChannelSize int,
) (g *Generator, err error) {
	g = &Generator{
		componentsCount: 1,
	}

	// Create random delays.
	{
		var delayAMs uint
		delayAMs, err = rnd.Uint(delayBeforeNextGenerationMinRangeBorderMs, delayBeforeNextGenerationMaxRangeBorderMs)
		if err != nil {
			return nil, err
		}

		var delayBMs uint
		delayBMs, err = rnd.Uint(delayBeforeNextGenerationMinRangeBorderMs, delayBeforeNextGenerationMaxRangeBorderMs)
		if err != nil {
			return nil, err
		}

		if delayAMs < delayBMs {
			g.delayBeforeNextGenerationMin = time.Millisecond * time.Duration(delayAMs)
			g.delayBeforeNextGenerationMax = time.Millisecond * time.Duration(delayBMs)
		} else {
			g.delayBeforeNextGenerationMin = time.Millisecond * time.Duration(delayBMs)
			g.delayBeforeNextGenerationMax = time.Millisecond * time.Duration(delayAMs)
		}
	}

	// Create random limits for random numbers.
	{
		var limitA uint
		limitA, err = rnd.Uint(generatedNumberMinRangeBorder, generatedNumberMaxRangeBorder)
		if err != nil {
			return nil, err
		}

		var limitB uint
		limitB, err = rnd.Uint(generatedNumberMinRangeBorder, generatedNumberMaxRangeBorder)
		if err != nil {
			return nil, err
		}

		if limitA < limitB {
			g.generatedNumberMin = limitA
			g.generatedNumberMax = limitB
		} else {
			g.generatedNumberMin = limitB
			g.generatedNumberMax = limitA
		}
	}

	g.output = make(chan uint, outputChannelSize)
	g.isWorking.Store(false)
	g.mustClose = make(chan bool, g.componentsCount)
	g.generatorWG = new(sync.WaitGroup)

	return g, nil
}

func (g *Generator) Start() (err error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.isWorking.Load() != false {
		return errors.New(ErrAlreadyStarted)
	}

	// Start the generator.
	g.generatorWG.Add(1)
	go g.generateNumbers()

	g.isWorking.Store(true)

	return nil
}

func (g *Generator) Stop() (err error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.isWorking.Load() != true {
		return errors.New(ErrNotStarted)
	}

	for i := 1; i <= g.componentsCount; i++ {
		g.mustClose <- true
	}

	// Wait for the generator to stop.
	g.generatorWG.Wait()

	g.isWorking.Store(false)

	return nil
}

func (g *Generator) generateNumbers() {
	defer g.generatorWG.Done()

	var err error

mainLoop:
	for {
		select {
		case quitSignal := <-g.mustClose:
			// Stop working.
			log.Printf("generator: a quit signal <%v> was received.\r\n", quitSignal)

			if quitSignal == true {
				break mainLoop
			}

		default:
			err = g.work()
			if err != nil {
				log.Println(err)
			}
		}
	}

	close(g.output)
}

func (g *Generator) work() (err error) {
	var randomNumber uint
	var sleepTime time.Duration

	// Generate the next random number.
	randomNumber, err = g.generateNumber()
	if err != nil {
		return err
	}

	// Prepare to sleep.
	sleepTime, err = g.generateSleepTime()
	if err != nil {
		return err
	}
	//log.Printf("sleepTime=%v.\r\n", sleepTime.String()) // Debug.

	// Publish the generated number.
	g.output <- randomNumber

	// Wait.
	time.Sleep(sleepTime)

	return nil
}

func (g *Generator) generateNumber() (number uint, err error) {
	return rnd.Uint(g.generatedNumberMin, g.generatedNumberMax)
}

func (g *Generator) generateSleepTime() (time.Duration, error) {
	delayMinMs := uint(g.delayBeforeNextGenerationMin / time.Millisecond)
	delayMaxMs := uint(g.delayBeforeNextGenerationMax / time.Millisecond)

	var sleepTimeMs uint
	var err error
	sleepTimeMs, err = rnd.Uint(delayMinMs, delayMaxMs)
	if err != nil {
		return 0, err
	}

	return time.Duration(sleepTimeMs) * time.Millisecond, nil
}

func (g *Generator) GetOutputChannel() chan uint {
	return g.output
}
