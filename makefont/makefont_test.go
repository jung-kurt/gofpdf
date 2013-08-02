package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMakefont(t *testing.T) {
	const expect = "Font definition file successfully generated"
	out, err := exec.Command("./makefont", "--dst=../font", "--embed",
		"--enc=../font/cp1252.map", "../font/calligra.ttf").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), expect) {
		t.Fatalf("Unexpected output from makefont")
	}
}
