package utils

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func LoadWords() []string {
	println("loading words")

	file, _ := os.Open("./" + WordsFilename)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	reader := bufio.NewReader(file)
	var words []string

	for {
		line, _, err := reader.ReadLine()
		if len(line) > 0 {
			words = append(words, string(line))
		}
		if err != nil {
			break
		}
	}
	return words
}

func writeInfos(wordInfos *sync.Map, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		println(err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	stringContent := ""

	(*wordInfos).Range(func(key, value interface{}) bool {
		stringContent += key.(string) + " " + fmt.Sprintf("%f", roundFloat(value.(float64), 2)) + "\n"
		return true
	})

	_, err = file.WriteString(stringContent)
	if err != nil {
		println(err)
		return
	}
}
