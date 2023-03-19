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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		if len(args) < 1 {
			fmt.Println("Specify projects you want to watch")
			os.Exit(1)
		}
		dirpaths := args[0:]

		var workTime time.Duration
		now := time.Now()

		isWorking := false

		var watchingDirs []string
		for _, dirpath := range dirpaths {
			gitLsDIrCommand := "cd " + dirpath + "&& git ls-files | sed -e '/^[^\\/]*$/d' -e 's/\\/[^\\/]*$//g' | sort | uniq"
			output, err := exec.Command("sh", "-c", gitLsDIrCommand).CombinedOutput()
			fmt.Printf("CombineOutput: %s, Error: %v\n", output, err)

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
