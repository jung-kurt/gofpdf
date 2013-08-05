package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMakefont(t *testing.T) {
	var out []byte
	var err error
	const expect = "Font definition file successfully generated"
	// Make sure makefont utility has been built before generating font definition file
	err = exec.Command("go", "build").Run()
	if err != nil {
		t.Fatal(err)
	}
	out, err = exec.Command("./makefont", "--dst=../font", "--embed",
		"--enc=../font/cp1252.map", "../font/calligra.ttf").CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), expect) {
		t.Fatalf("Unexpected output from makefont")
	}
}
