package flattener

import (
	"fmt"
	"github.com/charmbracelet/log"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Run(parentPath string, excludeSet map[string]struct{}) {
	var filesToMove []string
	var dirsToRemove []string

	fmt.Println()
	log.Infof("Flat: Processing directory: %s", parentPath)

	err := filepath.WalkDir(parentPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Errorf("Accessing path '%s' failed: %v\n", path, err)
			return err
		}

		if path == parentPath {
			return nil
		}

		if d.IsDir() {
			dirsToRemove = append(dirsToRemove, path)
			return nil
		}

		// no need to move top-level files
		if filepath.Dir(path) == parentPath {
			return nil
		}

		fileExt := filepath.Ext(path)
		cleanedFileExt := strings.TrimPrefix(fileExt, ".")
		if _, isExcluded := excludeSet[cleanedFileExt]; isExcluded {
			log.Infof("Excluding %s file, filename: %s", fileExt, filepath.Dir(path))
			return nil
		}

		filesToMove = append(filesToMove, path)
		return nil
	})

	if err != nil {
		log.Errorf("Failed during walk of '%s': %v", parentPath, err)
		return
	}

	// Move all files in filesToMove to parentPath
	for _, oldPath := range filesToMove {
		newPath := filepath.Join(parentPath, filepath.Base(oldPath))

		if _, err := os.Stat(newPath); err == nil {
			log.Warnf("'%s' already exists. Skipping move of '%s'", newPath, oldPath)
			continue
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Errorf("Failed to move '%s' to '%s': %v", oldPath, newPath, err)
		} else {
			log.Infof("Moved '%s' -> '%s'", oldPath, newPath)
		}
	}

	// delete a/folder1/folder2 before deleting a/folder1
	sort.Slice(dirsToRemove, func(i, j int) bool {
		return len(dirsToRemove[i]) > len(dirsToRemove[j])
	})

	for _, dir := range dirsToRemove {
		if err := os.Remove(dir); err == nil {
			log.Infof("Removed empty directory: '%s'", dir)
		}
	}
	log.Infof("Finished processing directory: %s", parentPath)
}
