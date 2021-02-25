package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	searchDir := flag.String("dir", ".", "The directory to search")
	flag.Parse()

	rootDir, err := filepath.Abs(*searchDir) // Turn possible relative directories into absolute paths
	if err != nil {
		log.Fatal("Error turning search directory in absolute path: " + err.Error())
	}

	filesWithNewlines, err := findHardcodedNewlines(rootDir)
	if err != nil {
		log.Fatal("Error finding hardcoded newlines: " + err.Error())
	}

	if len(filesWithNewlines) == 0 {
		fmt.Println("No hardcoded newlines have been found")
		os.Exit(0)
	}

	fmt.Println("Found hardcoded newlines in", len(filesWithNewlines), "files")

	for _, fileResult := range filesWithNewlines {
		for _, line := range fileResult.NewLines {
			fmt.Println(fileResult.Filename+": linefeed on line ", line)
		}
	}
	os.Exit(1)
}

func findHardcodedNewlines(rootDir string) ([]FileResults, error) {
	files, err := getJavaSourceFiles(rootDir)
	if err != nil {
		return nil, err
	}

	results := make([]FileResults, 0)

	for _, file := range files {
		newLines, err := findHardcodedNewlinesInFile(file)
		if err != nil {
			return nil, err
		}

		if len(newLines) == 0 {
			continue
		}

		results = append(results, FileResults{
			Filename: file[len(rootDir):],
			NewLines: newLines,
		})
	}

	return results, nil
}

func getJavaSourceFiles(rootDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.Name() == ".git" {
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) != ".java" {
			return nil
		}

		files = append(files, path)

		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

type FileResults struct {
	Filename string
	NewLines []int
}

func findHardcodedNewlinesInFile(filepath string) ([]int, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("Failed opening file: " + err.Error())
	}
	defer file.Close()

	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("Failed reading data from file: " + err.Error())
	}

	return findHardcodedNewlineInData(filedata), nil
}

func findHardcodedNewlineInData(filedata []byte) []int {
	linefeeds := []int{}

	lineNumber := 1

	for i := 0; i < len(filedata); i++ {

		if filedata[i] == '\\' {
			if i+1 == len(filedata) {
				continue
			}

			if filedata[i+1] != 'n' {
				continue
			}

			// 			if i >= 2 && filedata[i-1] == 'r' && filedata[i-2] == '\\' {
			// 				continue
			// 			}

			linefeeds = append(linefeeds, lineNumber)
		}

		if filedata[i] == '\n' {
			lineNumber++
		}
	}
	return linefeeds
}
