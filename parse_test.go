package main

import (
	"path/filepath"
	"testing"
)

func TestParseMultiFileDirectory(t *testing.T) {
	var path, _ = filepath.Abs(".")
	var p, _, e = Parse(path)
	if e != nil {
		t.Error(e.Error())
	}
	if len(p.Files) < 2 {
		t.Error(p.Files)
	}
}

func TestParseMissingDirectory(t *testing.T) {
	var path = ""
	var p, _, e = Parse(path)
	if e == nil {
		t.Error(p.Files)
	}
}
