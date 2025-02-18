package main

import (
	"fmt"
	"github.com/kr/pty"
	"io"
	"os"
	"syscall"
	"os/exec"
	"os/signal"
)

func main() {
	args := os.Args[1:]

	// Calling the srcds_linux executable in the same directory
	cmd := exec.Command("./srcds_linux", args...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = append(cmd.Env, "LD_LIBRARY_PATH=.:bin:"+os.Getenv("LD_LIBRARY_PATH"))

	// Redirecting the SIGINT or SIGKILL signal to the srcds_linux
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			err := cmd.Process.Signal(sig)
			if err != nil {
				fmt.Printf("[sourcelogger] couldn't redirect signal %v to srcds_linux\n", sig)
			}
		}
	}()

	// Starting the pseudo terminal for catching the stdout of gmod
	file, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("[sourcelogger] could't start the srcds_linux executable")
		panic(err)
	}

    // Redirect stdin -> gmod
    go func() { _, _ = io.Copy(file, os.Stdin) }()

	// Redirecting the output of gmod to the stdout
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		fmt.Println("[sourcelogger] finished with copy error (nothing to worry about)")
	} else {
		fmt.Println("[sourcelogger] finished")
	}

}
