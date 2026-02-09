package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

    fullPath := os.Args[1]

    fmt.Println(fullPath)

	dateFinderRegex := regexp.MustCompile(`\d{4}\.`)

	folderName := folderName(fullPath)
	fmt.Println("folderName:", folderName)

	dateIdx := dateFinderRegex.FindStringIndex(folderName)
	if (dateIdx == nil) {
		displayErrorAndWait("No date found in folder name.", nil)
		return
	}
	fmt.Println("idx:", dateIdx)

	filmWithPeriods := folderName[:dateIdx[0]] + `(` + folderName[dateIdx[0]:dateIdx[1]-1] + `)`
	film := regexp.MustCompile(`\.`).ReplaceAllString(filmWithPeriods, " ")
	fmt.Println("film:", film)
	
	err := findAndRenameFilm(fullPath, film)
	if err != nil {
		displayErrorAndWait("Error renaming film file.", err)
		return
	}

	fmt.Println("Renaming folder: ", fullPath, " to ", filepath.Join(filepath.Dir(fullPath), film))
	err = os.Rename(fullPath, filepath.Join(filepath.Dir(fullPath), film))
	if (err != nil) {
		displayErrorAndWait("Error renaming film folder.", err)
		return
	}
}

func findAndRenameFilm(path string, formattedFilmName string) error {
    c, err := os.ReadDir(path)
    if err != nil {
        return err
    }

	var extensions []string
	extensions = append(extensions, ".mp4", ".mkv", ".avi")

    for _, entry := range c {
		if (entry.IsDir()) {
			continue
		}

		var filmPath string
		name := entry.Name()

		for _, ext := range extensions {
			if strings.HasSuffix(name, ext) {
				filmPath = path
				break
			}
		}

		if (filmPath == "") {
			continue
		}

		oldPath := filepath.Join(path, name)
		newPath := filepath.Join(path, formattedFilmName+filepath.Ext(name))

		fmt.Println("renaming: ", oldPath, " to ", newPath)

		err := os.Rename(oldPath, newPath)
		return err
	}

	return fmt.Errorf("No film found in directory with specified extensions.")
}

func folderName(path string) string {
	pathParts := regexp.MustCompile(`\\`).Split(path, -1)
	return pathParts[len(pathParts)-1]
}

func displayErrorAndWait(message string, err error) {
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println("############ ERROR ############")
	fmt.Println(message)
	fmt.Println(" ")
	
	if (err != nil) {
		fmt.Println("Error: ", err)
		fmt.Println(" ")
	}
	
	fmt.Println("Press any key to close...")
	bufio.NewReader(os.Stdin).ReadBytes('\n') 
}