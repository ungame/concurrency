package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var wg sync.WaitGroup

func main() {
	password, err := bcrypt.GenerateFromPassword([]byte("000999888"), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}

	wordLists := getWordLists()

	var lines []string

	for _, wordlist := range wordLists {
		lines = append(lines, getLines(wordlist)...)
	}

	const plus = 100
	total := 0

	for i := 0; i < len(lines); i += plus {
		final := i + plus
		if final > len(lines) {
			final = len(lines)
		}
		wg.Add(1)
		go func(passwords []string) {
			defer wg.Done()
			for _, pwd := range passwords {
				if bcrypt.CompareHashAndPassword(password, []byte(pwd)) == nil {
					fmt.Println("================")
					fmt.Println("PASSWORD:", pwd)
					fmt.Println("================")
					os.Exit(1)
				} else {
					fmt.Printf("%s\r", pwd)
				}
			}
		}(lines[i:final])
		total++
	}
	fmt.Println("go rountines created", total)
	wg.Wait()

}

func getLines(wordlist string) []string {
	var lines []string
	f, err := os.OpenFile(wordlist, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Panicln("Open file failed:", err.Error())
	} else {

		var buf bytes.Buffer
		_, err = buf.ReadFrom(f)
		if err != nil {
			log.Panicln("Read failed failed:", err.Error())
		} else {

			for {
				line, err := buf.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Panicln("Read line failed:", err.Error())
				} else {
					lines = append(lines, strings.TrimSpace(line))
				}
			}

		}

		err = f.Close()
		if err != nil {
			log.Panicln(err)
		}
	}
	return lines
}

func getWordLists() []string {
	var files []string

	err := filepath.Walk("./wordlist", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	})
	if err != nil {
		log.Panicln(err)
	}
	return files
}
