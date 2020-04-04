package interrupt

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	funcs   = make(map[string]func())
	running bool
	mutex   = &sync.Mutex{}
)

// Add .
func Add(name string, function func()) {
	mutex.Lock()
	defer mutex.Unlock()
	if !running {
		start()
	}
	funcs[name] = function
}

// Remove .
func Remove(name string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(funcs, name)
}

func start() {
	interruptSig := make(chan os.Signal)
	signal.Notify(interruptSig, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-interruptSig
		fmt.Printf("Receive signal: %v, exiting...", sig)
		mutex.Lock()
		for _, function := range funcs {
			function()
		}
		os.Exit(1)
	}()
	running = true
}
