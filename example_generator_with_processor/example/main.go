package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kr/pretty"
	"github.com/vault-thirteen/junk/example_generator_with_processor/generator"
	"github.com/vault-thirteen/junk/example_generator_with_processor/processor"
)

func main() {
	var (
		delayMinMs           uint = 1000
		delayMaxMs           uint = 5000
		numMin               uint = 1
		numMax               uint = 100
		outputChannelSize         = 1024
		smaWindowSize             = 3
		printerSleepInterval      = time.Second * 10
		mustReadWholeQueue        = true
	)

	var (
		gen *generator.Generator
		prc *processor.Processor
		err error
	)
	gen, prc, err = startAllServices(
		delayMinMs,
		delayMaxMs,
		numMin,
		numMax,
		outputChannelSize,
		smaWindowSize,
		printerSleepInterval,
		mustReadWholeQueue,
	)
	mustBeNoError(err)

	waitForQuitSignal()

	err = stopAllServices(gen, prc)
	mustBeNoError(err)
}

func startAllServices(
	// Generator Settings.
	delayMinMs uint,
	delayMaxMs uint,
	numMin uint,
	numMax uint,
	outputChannelSize int,

	// Processor Settings.
	smaWindowSize int,
	printerSleepInterval time.Duration,
	mustReadWholeQueue bool,
) (
	gen *generator.Generator,
	prc *processor.Processor,
	err error,
) {
	gen, err = startServiceOne(
		delayMinMs,
		delayMaxMs,
		numMin,
		numMax,
		outputChannelSize,
	)
	if err != nil {
		return nil, nil, err
	}

	prc, err = startServiceTwo(
		gen.GetOutputChannel(),
		smaWindowSize,
		printerSleepInterval,
		mustReadWholeQueue,
	)
	if err != nil {
		return nil, nil, err
	}

	return gen, prc, nil
}

func startServiceOne(
	delayMinMs uint,
	delayMaxMs uint,
	numMin uint,
	numMax uint,
	outputChannelSize int,
) (gen *generator.Generator, err error) {
	gen, err = generator.New(
		delayMinMs,
		delayMaxMs,
		numMin,
		numMax,
		outputChannelSize,
	)
	if err != nil {
		return nil, err
	}

	_, _ = pretty.Println(gen)

	err = gen.Start()
	if err != nil {
		return nil, err
	}

	return gen, nil
}

func startServiceTwo(
	numbersChannel chan uint,
	smaWindowSize int,
	printerSleepInterval time.Duration,
	mustReadWholeQueue bool,
) (p *processor.Processor, err error) {
	p, err = processor.New(
		numbersChannel,
		printerSleepInterval,
		smaWindowSize,
		mustReadWholeQueue,
	)
	if err != nil {
		return nil, err
	}

	err = p.Start()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func waitForQuitSignal() {
	quitSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalsChannel, os.Interrupt)

	sig := <-quitSignalsChannel
	log.Printf("a signal <%v> was received.\r\n", sig)

	return
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func stopAllServices(
	gen *generator.Generator,
	prc *processor.Processor,
) (err error) {
	err = gen.Stop()
	if err != nil {
		return err
	}

	err = prc.Stop()
	if err != nil {
		return err
	}

	return nil
}
