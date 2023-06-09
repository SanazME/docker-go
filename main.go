package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	fmt.Printf("Running %v as %d\n", os.Args[0:], os.Getpid()) // 0: path 1: command 2: args and params

	// Run itself (this process again)
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
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

	cg()

	// Set hostname for the child process
	syscall.Sethostname([]byte("container"))

	// Change the root of what container can see, we want to have our own version of /proc directory in our container. we use ROOT_FOR_CONTAINER of ubuntu container:
	syscall.Chroot("vagrant/ubuntu-fs") // this directory will be our root
	syscall.Chdir("/")                  // change directory to root
	// we need to mount that directory (/proc) as proc pseudo filesystem so the kernel knows that we’re going to populate that with all the information about these running processes.
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...) // run whatever command is passed in + any params
	// wire up to see os stdin will go to our cmd stdin
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	syscall.Unmount("/proc", 0)
}

// /sys/fs/cgroup/memory
// /sys/fs/cgroup/pids
func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "snz"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "snz/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "snz/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "snz/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
