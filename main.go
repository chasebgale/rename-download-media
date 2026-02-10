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

	// Note: install.reg initializes the context menu command to call this program with the path of the folder that was right-clicked on as the first argument
    fullPath, err := getFullPathFromArgs()
    if err != nil {
        displayErrorAndWait("Error getting full path.", err)
        return
    }
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
	err = findAndRenameFilm(fullPath, filmName)
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

func getFullPathFromArgs() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("No path argument provided.")
	}

	if (!filepath.IsAbs(os.Args[1])) {
		return "", fmt.Errorf(`Provided path "%s" is not an absolute path.`, os.Args[1])
	}

	if (os.Args[1] == "") {
		return "", fmt.Errorf("Provided path is empty.")
	}

	return os.Args[1], nil
}

func findAndRenameFilm(path string, formattedFilmName string) error {
    c, err := os.ReadDir(path)
    if err != nil {
        return err
    }

	var extensions []string
	extensions = append(extensions, ".mp4", ".mkv", ".avi")

	var mediaFiles []os.DirEntry

    for _, entry := range c {
		if (entry.IsDir()) {
			continue
		}

		for _, ext := range extensions {
			if strings.HasSuffix(entry.Name(), ext) {
				mediaFiles = append(mediaFiles, entry)
				break
			}
		}
	}

	if (len(mediaFiles) == 0) {
		return fmt.Errorf("No film found in directory with specified extensions.")
	}

	// If there are multiple media files, only rename the largest one
	var largestFile os.DirEntry
	var largestFileSize int64 = 0

	for _, entry := range mediaFiles {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.Size() > largestFileSize {
			largestFileSize = info.Size()
			largestFile = entry
		}
	}

	oldPath := filepath.Join(path, largestFile.Name())
	newPath := filepath.Join(path, formattedFilmName+filepath.Ext(largestFile.Name()))

	fmt.Println("renaming: ", oldPath, " to ", newPath)

	err = os.Rename(oldPath, newPath)
	return err
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