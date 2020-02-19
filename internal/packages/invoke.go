// Copyright 2020 Manlio Perillo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The implementation of invokeGo is based on
// github.com/perillo/cmdgo/internal/invoke, but simplified to match the
// requirements of goprint.

package packages

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// gocmd is the go command to use.
var gocmd = "go"

// attr holds the attributes that will be applied to the cmd/go command.
type attr struct {
	// Env specifies the environment of the cmd/go command.
	// Each entry is of the form "key=value".
	// If Env is nil, the cmd/go command uses the current process's
	// environment.
	Env []string

	// Dir specifies the working directory of the cmd/go command.
	// If Dir is the empty string, the cmd/go command runs in the calling
	// process's current directory.
	Dir string
}

// invokeGo invokes a cmd/go command.
//
// If the cmd/go command returns a non 0 exit status, invokeGo will return a
// nil io.Reader and an error.
//
// If the cmd/go command returns a 0 exit status, invokeGo will return the
// stdout content as an io.Reader and a nil error.
//
// The child process stderr will be redirected to the parent process stderr.
func invokeGo(verb string, argv []string, attr *attr) (io.Reader, error) {
	argv = append([]string{verb}, argv...)
	stdout := new(bytes.Buffer)

	cmd := exec.Command(gocmd, argv...)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	if attr != nil {
		cmd.Dir = attr.Dir
		cmd.Env = attr.Env
	}

	if err := cmd.Run(); err != nil {
		argv := strings.Trim(fmt.Sprint(argv), "[]")

		return nil, fmt.Errorf("%s %s: %v", gocmd, argv, err)
	}

	return stdout, nil
}
