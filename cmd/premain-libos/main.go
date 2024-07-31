// Copyright (c) Edgeless Systems GmbH.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"

	//"syscall"

	marblePremain "github.com/edgelesssys/marblerun/marble/premain"
	"github.com/fatih/color"
	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
)

type count32 int32

func (c *count32) add(amount int32) int32 {
	return atomic.AddInt32((*int32)(c), amount)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func (c *count32) reset() {
	atomic.StoreInt32((*int32)(c), 0)
}

// libOS constants for specific checks.
// Use 1000 as a starting point for distinction.
const (
	gramine = iota + 1000
	occlum
)

func exit(format string, args ...interface{}) {
	// Print error message in red and append newline
	// then exit with error code 1
	msg := fmt.Sprintf("Error: %s\n", format)
	_, _ = color.New(color.FgRed).Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	log.SetPrefix("[PreMain] ")

	// Automatically detect libOS based on uname
	libOS, err := detectLibOS()
	if err != nil {
		exit("failed to detect libOS: %s", err)
	}

	// Use filesystem from libOS
	hostfs := afero.NewOsFs()

	var service string
	// Use different execution flows depending on libOS
	switch libOS {
	case gramine:
		log.Println("detected libOS: Gramine")

		// Gramine: Get service to launch before MarbleRun's premain
		service, err = prepareGramine(hostfs)
		if err != nil {
			exit("activating Gramine Marble failed: %s", err)
		}

	case occlum:
		log.Println("detected libOS: Occlum")

		// Occlum: Get entrypoint from MarbleRun manifest, adjust it to Occlum's quirks
		service, err = prepareOcclum(hostfs)
		if err != nil {
			exit("activating Occlum Marble failed: %s", err)
		}
	}

	if service[0] != '/' {
		service = "./" + service
	}
	// Test StartOperation RPC Request
	// Send a GRPC request through unsafe channel gRPC

	// Launch service
	log.Println("Launching service: ", service)
	//if err := unix.Exec(service, os.Args, os.Environ()); err != nil {
	//	exit("%s", err)
	//}
	//id, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Printf("Fork %d\n", id)
	//if id == 0 {
	//	log.Println("EXECCC")
	//	if err := unix.Exec(service, os.Args, os.Environ()); err != nil {
	//		exit("%s", err)
	//	}
	//}
	//log.Println("ODA")

	cmd := exec.Command(service, os.Args[1:]...)
	cmd.Env = os.Environ()
	for _, element := range cmd.Env {
		log.Println(element)
	}
	log.Printf("PATH: %s\n", os.Getenv("PATH"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the command: %v\n", err)
		os.Exit(1)
	}
	log.Println("Successfully started service. Waiting for it to finish...")
	cmd.Wait()
}

func detectLibOS() (int, error) {
	utsname := unix.Utsname{}
	if err := unix.Uname(&utsname); err != nil {
		return 0, err
	}

	// Clean utsname
	sysname := strings.ReplaceAll(string(utsname.Sysname[:]), "\x00", "")
	release := strings.ReplaceAll(string(utsname.Release[:]), "\x00", "")
	version := strings.ReplaceAll(string(utsname.Version[:]), "\x00", "")
	machine := strings.ReplaceAll(string(utsname.Machine[:]), "\x00", "")

	// Occlum detection
	// Taken from: https://github.com/occlum/occlum/blob/master/src/libos/src/misc/uname.rs
	if sysname == "Occlum" {
		return occlum, nil
	}

	// Gramine detection
	// This looks like a general Linux kernel name, making it harder to detect... But it's unlikely someone is running SGX code on Linux 3.10.0.
	// Taken from: https://github.com/gramineproject/gramine/blob/c83ec08f10cdbb3a258d18b118dd95602a55abc9/libos/src/sys/libos_uname.c
	if sysname == "Linux" && release == "3.10.0" && version == "1" && machine == "x86_64" {
		return gramine, nil
	}

	return 0, errors.New("cannot detect libOS")
}

func prepareGramine(hostfs afero.Fs) (string, error) {
	// Save the passed argument which is our service to spawn
	service := os.Args[0]

	// Run MarbleRun premain
	if err := marblePremain.PreMainEx(marblePremain.GramineQuoteIssuer{}, marblePremain.GramineActivate, hostfs, hostfs); err != nil {
		return "", err
	}

	return service, nil
}

func prepareOcclum(hostfs afero.Fs) (string, error) {
	// Run MarbleRun premain
	if err := marblePremain.PreMainEx(marblePremain.OcclumQuoteIssuer{}, marblePremain.ActivateRPC, hostfs, hostfs); err != nil {
		return "", err
	}

	// Check if the entrypoint defined in os.Args[0] actually exists
	if _, err := os.Stat(os.Args[0]); os.IsNotExist(err) {
		color.Red("ERROR: The entrypoint does not seem to exist: '$%s'", os.Args[0])
		color.Red("Please make sure that you define a valid entrypoint in your manifest (for example: /bin/hello_world).")
		return "", errors.New("invalid entrypoint definition in argv[0]")
	}

	// Modify os.Args[0] / argv[0] to only hold the program name, not the whole path, but keep it as service so we can correctly spawn the application.
	service := os.Args[0]
	os.Args[0] = filepath.Base(os.Args[0])

	return service, nil
}
