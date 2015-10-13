package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func matchTail(str, tailStr string) (match bool, headStr string) {
	sln := len(str)
	ln := len(tailStr)
	if sln > ln {
		match = str[sln-ln:] == tailStr
		if match {
			headStr = str[:sln-ln]
		}
	}
	return
}

func matchHead(str, headStr string) (match bool, tailStr string) {
	ln := len(headStr)
	if len(str) > ln {
		match = str[:ln] == headStr
		if match {
			tailStr = str[ln:]
		}
	}
	return
}

func main() {
	var err error
	var ok bool
	var showStr, name string
	err = filepath.Walk("pdf/reference", func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			name = filepath.Base(path)
			ok, name = matchTail(name, ".pdf")
			if ok {
				name = strings.Replace(name, "_", " ", -1)
				ok, showStr = matchHead(name, "Fpdf ")
				if ok {
					fmt.Printf("[%s](%s)\n", showStr, path)
				} else {
					ok, showStr = matchHead(name, "contrib ")
					if ok {
						fmt.Printf("[%s](%s)\n", showStr, path)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
