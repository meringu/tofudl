// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package mockmirror

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
)

func buildFake(t *testing.T) []byte {
	_, filename, _, _ := runtime.Caller(1)
	fakeDir := path.Join(path.Dir(filename), "fake")
	if err := os.MkdirAll(fakeDir, 0755); err != nil {
		t.Fatalf("Failed to create fake dir (%v)", err)
	}
	binaryPath := path.Join(fakeDir, "fake")
	if contents, err := os.ReadFile(binaryPath); err == nil {
		return contents
	}

	dir := path.Join(os.TempDir(), "fake")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	if err := os.WriteFile(path.Join(dir, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(dir, "main.go"), []byte(code), 0644); err != nil {
		t.Fatal()
	}

	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", "fake")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() != 0 {
			t.Fatalf("Build failed (exit code %d)", exitErr.ExitCode())
		} else {
			t.Fatalf("Build failed (%v)", err)
		}
	}

	contents, err := os.ReadFile(path.Join(dir, "fake"))
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(binaryPath, contents, 0755); err != nil {
		t.Fatalf("Failed to create fake binary at %s (%v)", binaryPath, err)
	}
	return contents
}

var code = `package main

func main() {
	print("Hello world!")
}
`

var gomod = `module fake

go 1.21`
