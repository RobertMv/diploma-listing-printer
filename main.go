package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ignored []string

func main() {
	// getting reader of console input
	scanner := bufio.NewScanner(os.Stdin)

	// getting files of directory
	fmt.Print("Enter absolute path to input directory: ")
	scanner.Scan()
	inputDir := scanner.Text()
	files := getDirFiles(inputDir)

	// print all of their names to console
	printDir(files)

	// get ignored of files to read
	fmt.Print("Enter list of directories you want to ignore (e.g. .git .idea config) or press ENTER to continue: ")
	scanner.Scan()
	constraintString := scanner.Text()
	ignored = strings.Split(constraintString, " ")

	fmt.Print("Enter absolute path to output directory: ")
	scanner.Scan()
	outputDir := scanner.Text()

	// prepare full listing
	printListing(outputDir, files)

	fmt.Println("\n\nPress the Enter Key to exit")
	fmt.Scanln()
}

func printDir(files []string) {
	fmt.Printf("Selected directory contains %d files:\n", len(files))
	for _, file := range files {
		fmt.Println("\t" + file)
	}
}

func printListing(outDirectory string, files []string) {
	// create out file
	out, err := os.Create(outDirectory + "/listing.txt")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	defer out.Close()
	var scanner *bufio.Scanner
	writer := bufio.NewWriter(out)
	printed := 0

	// iterate through slice of fileInfos
	for _, file := range files {
		if !isInIgnored(file) {
			// open each source file
			f, err := os.Open(file)
			if err != nil {
				fmt.Printf("%s", err.Error())
				return
			}
			scanner = bufio.NewScanner(f)

			dirs := strings.Split(f.Name(), "\\")
			writeLineInFile(writer, dirs[len(dirs)-1])
			// scan this source file
			for scanner.Scan() {
				// write its content to out file
				writeLineInFile(writer, scanner.Text())
			}
			writeLineInFile(writer, "\n\n")
			err = f.Close()
			if err != nil {
				fmt.Printf("%s", err.Error())
				return
			}
			printed++
		}
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	fmt.Printf("Printed %d files", printed)
}

func writeLineInFile(writer *bufio.Writer, text string) {
	_, err := fmt.Fprintln(writer, text)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func getDirFiles(inputDir string) []string {
	var files []string
	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// if not a directory and extension is appropriate
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	return files
}

func isInIgnored(path string) bool {
	for _, el := range ignored {
		if strings.Contains(path, el) {
			return true
		}
	}
	return false
}
