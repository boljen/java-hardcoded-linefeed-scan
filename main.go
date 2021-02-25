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

	rootDir, err := filepath.Abs(*searchDir)
	if err != nil {
		log.Fatal("Error turning search directory in absolute path: " + err.Error())
	}

	filesWithLF, err := findHardcodedLF(rootDir);
	if err != nil {
		log.Fatal(err)
	}

	if len(filesWithLF) == 0 {
		fmt.Println("No hardcoded linefeeds have been found")
		os.Exit(0)
	}

	fmt.Println("Found linefeeds in", len(filesWithLF), " files")

	for _, fileResult := range filesWithLF {
		for _, line := range fileResult.LineFeeds {
			fmt.Println(fileResult.Filename + ": linefeed on line ", line)
		}
	}
	os.Exit(-1)
}

func findHardcodedLF(rootDir string) ([]FileResults, error) {
	files, err := getJavaSourceFiles(rootDir)
	if err != nil {
		return nil, err
	}

	results := []FileResults{}

	for _, file := range files {
		result, err := findLineFeedsInFile(file)
		if err != nil {
			return nil, err
		}

		result.Filename = result.Filename[len(rootDir):]

		if len(result.LineFeeds) > 0 {
			results = append(results, *result)
		}
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
	LineFeeds []int
}

func findLineFeedsInFile(filepath string) (*FileResults, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("Failed opening file for line feeds: " + err.Error())
	}
	defer file.Close()

	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("Failed reading data for line feeds: " + err.Error())
	}

	linefeeds := findLinefeedInData(filedata)
	return &FileResults{
		Filename: filepath,
		LineFeeds: linefeeds,
	}, nil
}

func findLinefeedInData(filedata []byte) ([]int) {
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
