/*
 * Copyright 2021 Filip Varga
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package run

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func exec1(cmd exec.Cmd) error {
	err := cmd.Run()
	if err == nil {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		es := ws.ExitStatus()
		if es != 0 {
			err = fmt.Errorf("command failed, exit status %d is non zero", es)
		}
	}
	return err
}

func exec2(cmd exec.Cmd) ([]byte, error) {
	output, err := cmd.CombinedOutput()
	if err == nil {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		es := ws.ExitStatus()
		if es != 0 {
			err = fmt.Errorf("command failed, exit status %d is non zero", es)
		}
	}
	return output, err
}

func Exec1(command string, args ...string) error {
	return exec1(*exec.Command(command, args...))
}

func Exec2(command string, args ...string) ([]byte, error) {
	return exec2(*exec.Command(command, args...))
}

func Exec3(quiet bool, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	if quiet == false {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return exec1(*cmd)
}

/* vim: set ts=2: */
