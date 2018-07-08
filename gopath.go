package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// GOPATH returns the active GOPATH settings for the runtime.
func GOPATH() string {
	var p = os.Getenv("GOPATH")
	if p == "" {
		return defaultGOPATH()
	}
	return p
}

func defaultGOPATH() string {
	var env = "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	}
	if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		var def = filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			return ""
		}
		return def
	}
	return ""
}
