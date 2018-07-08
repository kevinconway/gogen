package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGOPATHUsesEnvWhenPresent(t *testing.T) {
	var original, restore = os.LookupEnv("GOPATH")
	defer func() {
		if restore {
			_ = os.Setenv("GOPATH", original)
		}
	}()
	_ = os.Setenv("GOPATH", "/test")
	var p = GOPATH()
	if p != "/test" {
		t.Error(p)
	}
}

func TestGOPATHDefault(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux only") // TODO: Add per GOOS build flagged files for this
	}
	var original, restore = os.LookupEnv("GOPATH")
	defer func() {
		if restore {
			_ = os.Setenv("GOPATH", original)
		}
	}()
	_ = os.Unsetenv("GOPATH")
	var p = GOPATH()
	if p != filepath.Join(os.Getenv("HOME"), "go") {
		t.Error(p)
	}
}
