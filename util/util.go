package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DirectoriesWithAbsolutePath(relativePaths []string) []string {
	absolutePaths := []string{}
	for _, p := range relativePaths {
		path, _ := filepath.Abs(p)
		absolutePaths = append(absolutePaths, path)
	}
	return absolutePaths
}

func GetParentDirPath(parentDirPaths []string, dir string) string {
	for _, parentDir := range parentDirPaths {
		relPath, err := filepath.Rel(parentDir, dir)
		if err != nil {
			panic(err)
		}

		if relPath != "." && relPath != ".." && filepath.Base(relPath) != relPath {
			// dir has parent directory
			return parentDir
		}
	}
	fmt.Printf("%v has no watching parent directory. Terminate watcher.", dir)
	os.Exit(1)

	return ""
}

func GetWatchingDirs(dirPaths []string) []string {
	var watchingDirs []string
	for _, dirPath := range dirPaths {
		disableEscapeMultiByteCharsCommand := "cd " + dirPath + "&& git config core.quotepath false"
		output, err := exec.Command("sh", "-c", disableEscapeMultiByteCharsCommand).CombinedOutput()
		if err != nil {
			panic(err)
		}

		gitLsDirCommand := "cd " + dirPath + "&& git ls-files | sed -e '/^[^\\/]*$/d' -e 's/\\/[^\\/]*$//g' | sort | uniq"
		output, err = exec.Command("sh", "-c", gitLsDirCommand).CombinedOutput()
		if err != nil {
			panic(err)
		}

		dirs := strings.Split(string(output), "\n")

		var fullPathDirs []string
		for _, dir := range dirs {
			fullPathDirs = append(fullPathDirs, filepath.Join(dirPath, dir))
		}
		watchingDirs = append(watchingDirs, fullPathDirs...)
	}

	return watchingDirs
}
