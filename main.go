package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker 			run image <cmd> <params>
// go run main.go	run 	  <cmd> <params>
func main() {
	switch os.Args[1] {
	case "run":
		run()

	default:
		panic("bad command!")

	}

}

func run() {
	fmt.Printf("Running %v\n", os.Args[0:]) // 0: path 1: command 2: args and params

	cmd := exec.Command(os.Args[2], os.Args[3:]...) // run whatever command is passed in + any params
	// wire up to see os stdin will go to our cmd stdin
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// to containerize our command, we create a namespace
	// Cloneflags: cloning is what creates a new process in which we run our commands
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS, // New mount namespace group
	}

	cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
