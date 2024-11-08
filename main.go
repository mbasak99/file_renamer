package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

func renameFile(files []fs.DirEntry, oldName string, newName string, c chan string, dir string) {
	defer close(c)

	for _, file := range files {
		if strings.Contains(file.Name(), oldName) {
			fmt.Println("Old: " + file.Name())
			newFileName := dir + strings.Replace(file.Name(), oldName, newName, 1)

			err := os.Rename(dir+file.Name(), newFileName)
			if err != nil {
				log.Fatal(err)
			}

			c <- newFileName
		}
	}
}

func main() {
	// Flags for program
	dir := flag.String("dir", "", "The directory where the files are located. NO DEFAULT.")
	mode := flag.Int("mode", 1, "Choose between single-threaded (-1) or multi-threaded mode (1). Default: 1 (multi-threaded).")
	oldKeyword := flag.String("old", "", "The word that will be targeted for replacement. NO DEFAULT.")
	newKeyword := flag.String("new", "", "The word that will be used to replace the old keyword. NO DEFAULT.")

	flag.Parse()

	fmt.Printf("%s\n", *dir)
	fmt.Printf("%d\n", *mode)
	fmt.Printf("%s\n", *oldKeyword)
	fmt.Printf("%s\n", *newKeyword)

	// error check
	if len(*dir) == 0 {
		log.Fatalln("Please enter a directory")
	}

	if len(*oldKeyword) == 0 {
		log.Fatalln("Please enter a keyword to replace.")
	}

	if len(*newKeyword) == 0 {
		log.Fatalln("Please enter a keyword that will be used to replace with.")
	}

	if !(*mode == 1 || *mode == -1) {
		log.Fatalln("Please enter a valid mode.")
	}

	files, err := os.ReadDir(*dir)

	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	// multi-thread
	if *mode == 1 {
		c := make(chan string, len(files))

		go renameFile(files, *oldKeyword, *newKeyword, c, *dir)

		for file := range c {
			fmt.Println("New File Name: " + file)
		}
	} else if *mode == -1 { // single-thread
		for _, file := range files {
			if strings.Contains(file.Name(), *oldKeyword) {
				fmt.Println("Old: " + file.Name())
				newFileName := *dir + strings.Replace(file.Name(), *oldKeyword, *newKeyword, 1)
				err := os.Rename(*dir+file.Name(), newFileName)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("New File Name: " + newFileName)
			}
		}
	}

	log.Printf("Time taken: %s", time.Since(start))
}
