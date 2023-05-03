/*
Copyright © 2023 Seiya Miyata <odradek38@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// hakaruCmd represents the hakaru command
var hakaruCmd = &cobra.Command{
	Use:   "hakaru",
	Short: "Watch your work time on each projects",
	Long: `Watch your work time on each projects. You can specify directory paths to your GitHub projects, and then
	a watcher process would measure work time on each projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		if len(args) < 1 {
			fmt.Println(`Specify path to projects you want to watch like "hakaru /path/to/project1 /path/to/project2"`)
			os.Exit(1)
		}
		dirPaths := args[0:]

		validateDirectories(dirPaths)

		wt := make(map[string]time.Duration)
		for _, path := range dirPaths {
			wt[path] = 0 * time.Minute
		}
		workTime := workTime{
			dirPaths: directoriesWithAbsolutePath(dirPaths),
			workTime: wt,
		}

		now := time.Now()

		isWorking := false

		watchingDirs := getWatchingDirs(dirPaths)

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Has(fsnotify.Write) {
						log.Println("modified file: ", event.Name)

						parentDirPath := getParentDirPath(directoriesWithAbsolutePath(watchingDirs), event.Name)

						if isWorking {
							if time.Since(now) <= 5*time.Minute {
								workTime.workTime[parentDirPath] += time.Since(now)
							} else {
								workTime.workTime[parentDirPath] += 5 * time.Minute
							}
						}

						now = time.Now()
						log.Println("work time: ", workTime.workTime[parentDirPath])

						isWorking = true
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error: ", err)
				case s := <-sigs:
					log.Println("Signal accepted:", s)
					log.Println("Directories is", strings.Join(directoriesWithAbsolutePath(dirPaths), ", "))
					log.Println("Total work time is", totalWorkTime(workTime.workTime).String())
					for _, p := range workTime.dirPaths {
						log.Println("Directory:", p, "Work time:", workTime.workTime[p].String())
					}
					os.Exit(1)
				}
			}
		}()

		addWatchingDirsToWatcher(directoriesWithAbsolutePath(watchingDirs), watcher)

		<-make(chan struct{})
	},
}

func init() {
	rootCmd.AddCommand(hakaruCmd)
}

type workTime struct {
	dirPaths []string
	workTime map[string]time.Duration
}

func validateDirectories(dirPaths []string) {
	for _, dir := range dirPaths {
		if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
			fmt.Printf("Directory %s does not exist\n", dir)
			os.Exit(1)
		}
	}
}

func directoriesWithAbsolutePath(relativePaths []string) []string {
	absolutePaths := []string{}
	for _, p := range relativePaths {
		path, _ := filepath.Abs(p)
		absolutePaths = append(absolutePaths, path)
	}
	return absolutePaths
}

func getWatchingDirs(dirPaths []string) []string {
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

func addWatchingDirsToWatcher(watchingDirs []string, watcher *fsnotify.Watcher) {
	for _, dir := range watchingDirs {
		err := watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func totalWorkTime(workTimes map[string]time.Duration) time.Duration {
	var total time.Duration
	for _, t := range workTimes {
		total += t
	}
	return total
}

func getParentDirPath(parentDirPaths []string, dir string) string {
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
