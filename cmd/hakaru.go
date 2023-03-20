/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

		if len(args) < 1 {
			fmt.Println(`Specify path to projects you want to watch like "hakaru ../path/to/project"`)
			os.Exit(1)
		}
		dirpaths := args[0:]

		for _, dir := range dirpaths {
			if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
				fmt.Printf("Directory %s does not exist\n", dir)
				os.Exit(1)
			}
		}

		var workTime time.Duration
		now := time.Now()

		isWorking := false

		var watchingDirs []string
		for _, dirpath := range dirpaths {
			gitLsDIrCommand := "cd " + dirpath + "&& git ls-files | sed -e '/^[^\\/]*$/d' -e 's/\\/[^\\/]*$//g' | sort | uniq"
			output, err := exec.Command("sh", "-c", gitLsDIrCommand).CombinedOutput()
			if err != nil {
				panic(err)
			}

			dirs := strings.Split(string(output), "\n")

			var fullPathDirs []string
			for _, dir := range dirs {
				fullPathDirs = append(fullPathDirs, dirpath+"/"+dir)
			}
			watchingDirs = append(watchingDirs, fullPathDirs...)
		}

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Has(fsnotify.Write) {
						log.Println("modified file: ", event.Name)

						if isWorking {
							if time.Since(now) <= 5*time.Minute {
								workTime += time.Since(now)
							} else {
								workTime += 5 * time.Minute
							}
						}

						now = time.Now()
						log.Println("work time: ", workTime)

						isWorking = true
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error: ", err)
				}
			}
		}()

		for _, dir := range watchingDirs {
			err = watcher.Add(dir)
			if err != nil {
				log.Fatal(err)
			}
		}

		<-make(chan struct{})
	},
}

func init() {
	rootCmd.AddCommand(hakaruCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hakaruCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hakaruCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
