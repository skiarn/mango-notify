package io

import (
	"fmt"
	"os"
	"time"
)

//OnFileChange detects if the file is changed.
func OnFileChange(file *string, doneChan chan bool) error {
	var err error
	go func(doneChan chan bool) {
		defer func() {
			doneChan <- true
		}()
		err = watchFile(*file)
		fmt.Println("File has been changed")
	}(doneChan)
	return err
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
