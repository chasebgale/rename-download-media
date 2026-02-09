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

	// install.reg initializes the context menu command to call this program with the path of the folder that was right-clicked on as the first argument
    fullPath := os.Args[1]
    fmt.Println("path: ", fullPath)

	folderName := folderName(fullPath)
	fmt.Println("folder: ", folderName)

	dateIdx := dateIndexFromFolderName(folderName)
	if (dateIdx == nil) {
		displayErrorAndWait("No date found in folder name.", nil)
		return
	}

	// Format the preferred film name as "Film Name (Year)" by using the date found in the folder name and removing any periods (This is the format Jellyfin expects)
	// TODO: Perhaps scan original name for dashes and periods and replace whichever is more common, as some rare files use dashes over periods
	filmWithPeriods := folderName[:dateIdx[0]] + `(` + folderName[dateIdx[0]:dateIdx[1]-1] + `)`
	filmName := regexp.MustCompile(`\.`).ReplaceAllString(filmWithPeriods, " ")
	fmt.Println("film:", filmName)
	
	// Rename the film file with our formatted name
	// TODO: Upate func to only rename largest media file found to account for sample files, etc
	err := findAndRenameFilm(fullPath, filmName)
	if err != nil {
		displayErrorAndWait("Error renaming film file.", err)
		return
	}

	// Rename the folder with our formatted name
	fmt.Println("Renaming folder: ", fullPath, " to ", filepath.Join(filepath.Dir(fullPath), filmName))
	err = os.Rename(fullPath, filepath.Join(filepath.Dir(fullPath), filmName))
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

func dateIndexFromFolderName(folderName string) []int {
	dateFinderRegex := regexp.MustCompile(`\d{4}\.`)
	return dateFinderRegex.FindStringIndex(folderName)
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