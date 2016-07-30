package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var buf []byte
	var pos int
	var b byte
	var err error
	buf, err = ioutil.ReadAll(os.Stdin)
	if err == nil {
		for _, b = range buf {
			fmt.Printf("0x%02X, ", b)
			pos++
			if pos >= 16 {
				fmt.Println("")
				pos = 0
			}
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
