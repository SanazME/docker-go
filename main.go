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
	case "child":
		child()

	default:
		panic("bad command!")

	}

}

// This process creates a namespace
func run() {
	fmt.Printf("Running %v\n", os.Args[0:]) // 0: path 1: command 2: args and params

	// Run itself (this process again)
	cmd := exec.Command("/proc/self/exe", append([]string {"child"}, os.Args[2:]...)...)
	// wire up to see os stdin will go to our cmd stdin
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// to containerize our command, we create a namespace
	// Cloneflags: cloning is what creates a new process in which we run our commands
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWPID, // New mount namespace group with PID starts from 1 in the child process
	}
	cmd.Run() // inside here we're actually getting a clone of a new process and a namespace (it doesn't exit before that)

}

// In run, we clone a new process and in child we actually run that cloned process with a hostname that we set
func child() {
	fmt.Printf("Running %v\n", os.Args[0:]) // 0: path 1: command 2: args and params

	// Set hostname for the child process
	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...) // run whatever command is passed in + any params
	// wire up to see os stdin will go to our cmd stdin
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
