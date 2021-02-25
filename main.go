package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const root = "C:\\Users\\chris\\dev\\repos\\github.com\\switchfully\\java-feb-2021"

func main() {

	searchDir := root


	var files []string

	err := filepath.Walk(searchDir, visit(&files))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		checkForLinuxNewline(file)
	}

}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.Name() == ".git" {
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) != ".java" {
			return nil
		}

		*files = append(*files, path)

		return nil
	}
}

func checkForLinuxNewline(filepath string) {

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("failed opening file", err)
	}
	defer file.Close()

	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("failed reading file")
	}

	lineNumber := 1

	for i := 0; i < len(filedata); i++ {

		if filedata[i] == '\\' {
			if i+1 == len(filedata) {
				continue
			}

			if filedata[i+1] != 'n' {
				continue
			}

			if i >= 2 && filedata[i-1] == 'r' && filedata[i-2] == '\\' {
				continue
			}

			fmt.Println(filepath[len(root):], " @ line ", lineNumber)

		}

		if filedata[i] == '\n' {
			lineNumber++
		}
	}

}
