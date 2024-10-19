package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type FileStats struct {
	BytesCount int64
	LinesCount int64
	WordCount  int64
	CharCount  int64
}

func getFileStats(reader io.Reader) (FileStats, error) {
	var stats FileStats

	data, err := io.ReadAll(reader)
	if err != nil {
		return stats, err
	}

	stats.BytesCount = int64(len(data))
	stats.WordCount = int64(len(bytes.Fields(data)))
	stats.CharCount = int64(len(bytes.Runes(data)))

	for _, b := range data {
		if b == '\n' {
			stats.LinesCount++
		}
	}

	return stats, nil
}

func displayStats(stats FileStats, filename string, printLines, printWords, printBytes, printChar bool) {
	fmt.Printf("    ")
	if printLines {
		fmt.Printf("%d   ", stats.LinesCount)
	}
	if printWords {
		fmt.Printf("%d   ", stats.WordCount)
	}
	if printBytes {
		fmt.Printf("%d   ", stats.BytesCount)
	}
	if printChar {
		fmt.Printf("%d   ", stats.CharCount)
	}
	fmt.Printf("%s\n", filename)
}

func setupFlags() (printBytes, printLines, printWords, printChar bool) {
	flag.BoolVar(&printBytes, "c", false, "Display the number of bytes in the file.")
	flag.BoolVar(&printLines, "l", false, "Display the number of lines in the file.")
	flag.BoolVar(&printWords, "w", false, "Display the number of words in the file.")
	flag.BoolVar(&printChar, "m", false, "Display the number of characters in the file.")
	flag.Parse()

	// Default behavior: print all stats if no flags are specified
	if !printBytes && !printLines && !printWords && !printChar {
		printBytes, printLines, printWords = true, true, true
	}

	return
}

func getReader(filename string) (io.Reader, string, error) {
	if filename == "" {
		return os.Stdin, "stdin", nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	return file, filename, nil
}

func main() {
	printBytes, printLines, printWords, printChar := setupFlags()

	filename := flag.Arg(0)
	reader, fname, err := getReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if file, ok := reader.(*os.File); ok {
			file.Close()
		}
	}()

	stats, err := getFileStats(reader)
	if err != nil {
		log.Fatal(err)
	}

	displayStats(stats, fname, printLines, printWords, printBytes, printChar)
}
