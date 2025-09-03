package terminal

import (
	"bufio"
	"fmt"
	"golang-web-core/logging"
	"golang-web-core/service"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
)

type Terminal struct {
	services []service.Service
	shutdown bool
	logger   logging.Logger
}

func NewTerminal(services []service.Service) *Terminal {
	logger := logging.NewLogger()
	logger.ServiceName = "Terminal"

	return &Terminal{
		services: services,
		shutdown: false,
		logger:   logger,
	}
}

func (t *Terminal) Start() {
	t.checkServices()
	t.startServices()
	t.catchSignals()
	t.startTerminal()
}

func (t *Terminal) checkServices() {
	names := []string{}
	for _, s := range t.services {
		if slices.Contains(names, s.Name()) {
			t.logger.Errorf("Found more than one service with the same name: %v", s.Name())
			os.Exit(1)
		}

		names = append(names, s.Name())
	}
}

func (t *Terminal) startServices() {
	wg := sync.WaitGroup{}

	for _, s := range t.services {
		wg.Add(1)
		go func() {
			s.Start()
			wg.Done()
		}()
	}

	wg.Wait()
}

func (t *Terminal) catchSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		<-signals
		t.logger.Info("Received shutdown signal, shutting down")
		t.handleShutdown()
	}()
}

func (t *Terminal) startTerminal() {
	reader := bufio.NewReader(os.Stdin)

	for {
		if t.shutdown {
			break
		}

		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		t.processCommand(command)
	}

	t.logger.Info("Terminal shutdown")
}

func (t *Terminal) getServiceIndex(serviceName string) int {
	for i, s := range t.services {
		if s.Name() == serviceName {
			return i
		}
	}
	return -1
}

func (t *Terminal) forceQuit() {
	os.Exit(0)
}
