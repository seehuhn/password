// Some of the code in this file is taken from
// https://code.google.com/p/go/source/browse/?repo=crypto#hg%2Fssh%2Fterminal
// and then has been modified by Jochen Voss.
//
// The original code is distributed under the following license:
//
//	   Copyright 2011 The Go Authors. All rights reserved.
//	   Use of this source code is governed by a BSD-style
//	   license that can be found in the LICENSE file.
//
// All changes to the original code are distributed under the
// following license:
//
//	   Copyright 2013 Jochen Voss. All rights reserved.
//	   Use of this source code is governed by a BSD-style
//	   license that can be found in the LICENSE file.

// +build linux,!appengine darwin

// Package password provides a function to read passwords on the
// command line on Linux and BSD Unix (including MacOS X) systems.
package password

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"unsafe"
)

// Read prints the given prompt to standard output and then reads a
// line of input from standard input with echoing of input disabled.
// This is commonly used for inputting passwords and other sensitive
// data.  The byte slice returned does not include the terminating
// "\n".
func Read(prompt string) ([]byte, error) {
	fd := 0

	_, err := os.Stdout.Write([]byte(prompt))
	if err != nil {
		return nil, err
	}

	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		ioctlReadTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	newState := oldState
	newState.Lflag &^= syscall.ECHO
	newState.Lflag |= syscall.ICANON | syscall.ISIG
	newState.Iflag |= syscall.ICRNL
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
		ioctlWriteTermios, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	// restore terminal after keyboard interrupts
	signalChannel := make(chan os.Signal)
	done := make(chan bool)
	signal.Notify(signalChannel,
		os.Signal(syscall.SIGINT), os.Signal(syscall.SIGTERM))
	go func() {
		select {
		case sig := <-signalChannel:
			signal.Stop(signalChannel)
			syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
				ioctlWriteTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
			p, _ := os.FindProcess(os.Getpid())
			p.Signal(sig)
		case <-done:
		}
	}()
	defer close(done)

	defer func() {
		syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd),
			ioctlWriteTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
		os.Stdout.Write([]byte("\n"))
	}()

	var ret []byte
	var buf [16]byte
	for {
		n, err := syscall.Read(fd, buf[:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			if len(ret) == 0 {
				return nil, io.EOF
			}
			break
		}
		if buf[n-1] == '\n' {
			n--
		}
		ret = append(ret, buf[:n]...)
		if n < len(buf) {
			break
		}
	}
	return ret, nil
}
