package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("must provide at least 2 args, the profile to transition to, and the program to run (with any args")
	}
	desiredProfile := os.Args[1]
	prog := os.Args[2]
	args := os.Args[3:]

	// TODO: we should probably be extra safe like runc and verify that we are
	//       writing to something in procfs
	// see https://github.com/opencontainers/runc/commit/d463f6485b809b5ea738f84e05ff5b456058a184
	// and CVE-2019-16884
	f, err := os.OpenFile("/proc/self/attr/exec", os.O_WRONLY, 0)
	if err != nil {
		fmt.Println("could not open process exec attr to transition apparmor profile:", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "exec %s", desiredProfile)
	if err != nil {
		fmt.Printf("could not set process exec attr to %s: %v\n", desiredProfile, err)
		os.Exit(1)
	}

	// exec the program with the new profile
	// can't use exec/cmd here because we need to actually exec, i.e. not fork
	// - but if this code needs to be a new process and we need to keep running,
	// we could fork, then in the child process set the proc exec attr, then
	// exec there
	if err := syscall.Exec(prog, args, os.Environ()); err != nil {
		fmt.Println("failed to re-exec:", err)
		os.Exit(1)
	}

	// should be impossible to reach here
}
