package processor

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/vault-thirteen/auxie/SAAC"
	sma "github.com/vault-thirteen/auxie/SMA"
	"go.uber.org/atomic"
	"go.uber.org/multierr"
)

const (
	ErrAlreadyStarted         = "already started"
	ErrNotStarted             = "not started"
	ErrNumbersChannelIsClosed = "numbers channel is closed"
)

type Processor struct {
	componentsCount      int
	numbersChannel       chan uint
	printerSleepInterval time.Duration
	mustReadWholeQueue   bool

	saac *saac.Calculator
	sma  *sma.Calculator

	lastSAAValue string
	lastSMAValue sma.ValueType

	lock        sync.Mutex
	isWorking   atomic.Bool
	mustClose   chan bool
	processorWG *sync.WaitGroup
	printerWG   *sync.WaitGroup
}

func New(
	numbersChannel chan uint,
	printerSleepInterval time.Duration,
	smaWindowSize int,
	mustReadWholeQueue bool,
) (p *Processor, err error) {
	p = &Processor{
		componentsCount:      2,
		numbersChannel:       numbersChannel,
		printerSleepInterval: printerSleepInterval,
		mustReadWholeQueue:   mustReadWholeQueue,
	}

	p.saac = saac.New()

	p.sma, err = sma.New(smaWindowSize)
	if err != nil {
		return nil, err
	}

	p.isWorking.Store(false)
	p.mustClose = make(chan bool, p.componentsCount)
	p.processorWG = new(sync.WaitGroup)
	p.printerWG = new(sync.WaitGroup)

	return p, nil
}

func (p *Processor) Start() (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.isWorking.Load() != false {
		return errors.New(ErrAlreadyStarted)
	}

	// Start the processor.
	p.processorWG.Add(1)
	go p.processData()

	// Start the printer.
	p.printerWG.Add(1)
	go p.printData()

	p.isWorking.Store(true)

	return nil
}

func (p *Processor) Stop() (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.isWorking.Load() != true {
		return errors.New(ErrNotStarted)
	}

	for i := 1; i <= p.componentsCount; i++ {
		p.mustClose <- true
	}

	// Wait for the processor to stop.
	p.processorWG.Wait()

	// Wait for the printer to stop.
	p.printerWG.Wait()

	p.isWorking.Store(false)

	return nil
}

func (p *Processor) processData() {
	defer p.processorWG.Done()

	var err error
	var number uint

mainLoop:
	for {
		select {
		case quitSignal := <-p.mustClose:
			// Stop working.
			log.Printf("processor: a quit signal <%v> was received.\r\n", quitSignal)

			if quitSignal == true {
				break mainLoop
			}

		default:
			number, err = p.receiveItem()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second) // Cool-down period.
				continue
			}

			err = p.processItem(number)
			if err != nil {
				log.Println(err)
			}
		}
	}

	// Process all the items rest in the queue.
	// This behaviour is optional and may be disabled.
	if p.mustReadWholeQueue {
		for number = range p.numbersChannel {
			err = p.processItem(number)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (p *Processor) receiveItem() (number uint, err error) {
	var ok bool

	// Receive the number.
	number, ok = <-p.numbersChannel
	if !ok {
		return 0, errors.New(ErrNumbersChannelIsClosed)
	}

	log.Printf("a new number is received: %v.\r\n", number) // Debug.

	return number, nil
}

func (p *Processor) processItem(number uint) (err error) {
	var saacErr error
	p.lastSAAValue, saacErr = p.saac.AddItemAndGetAverage(
		strconv.FormatUint(uint64(number), 10),
	)

	var smaErr error
	p.lastSMAValue, smaErr = p.sma.AddItemAndGetSMA(sma.ValueType(number))
	if (smaErr != nil) &&
		(smaErr.Error() == sma.ErrDataSetIsCold) {
		log.Println(smaErr)
		smaErr = nil
	}

	err = multierr.Combine(saacErr, smaErr)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) printData() {
	defer p.printerWG.Done()

mainLoop:
	for {
		select {
		case quitSignal := <-p.mustClose:
			// Stop working.
			log.Printf("printer: a quit signal <%v> was received.\r\n", quitSignal)

			if quitSignal == true {
				break mainLoop
			}

		default:
			time.Sleep(p.printerSleepInterval)
			p.printItem()
		}
	}
}

func (p *Processor) printItem() {
	fmt.Printf("Simple average is %v.\r\n", p.lastSAAValue)
	fmt.Printf("Simple moving average is %v.\r\n", p.lastSMAValue)
}
