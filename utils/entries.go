package utils

import (
	"fmt"
	"sync"
)

// Solver Interactively play wordle
func Solver() {
	words := LoadWords()
	//allWords := words
	bestWord := "tares"
	wordInfos := sync.Map{}
	input := ""
	print("What is the feedback of ", bestWord, ": ")
	_, err := fmt.Scan(&input)
	if err != nil {
		panic(err)
	}
	words = GenerateNewPossible(bestWord, input, words)
	for len(words) != 1 {
		println("\n", len(words), " words left")
		// Easy mode in wordle, but doesn't work yet
		//bestWord = generateBest(&words, &allWords, &wordInfos)
		print(len(words))
		bestWord = GenerateBest(&words, &words, &wordInfos, FancyPrint)
		print("\nWhat is the feedback of ", bestWord, ": ")
		_, err := fmt.Scan(&input)
		if err != nil {
			panic(err)
		}
		// Handles cases when the input is not a 5 character string(not a valid feedback)
		// Current behavior is to remove the best word from the list
		// The intention is to use this if wordle doesn't accept the word
		if len(input) != 5 {
			var temp []string
			for _, word := range words {
				if word != bestWord {
					temp = append(temp, word)
				}
				words = temp
			}
			continue
		}
		words = GenerateNewPossible(bestWord, input, words)
	}
	println("The word is", words[0])
}

// GenerateBestAndWriteFromStart Generate the best opener and write the output to infos.txt
func GenerateBestAndWriteFromStart() {
	words := LoadWords()
	fmt.Printf("Loaded %d words", len(words))
	wordInfos := sync.Map{}
	for _, word := range words {
		wordInfos.Store(word, float64(0))
	}
	fmt.Println("\nThe best opener is: ", GenerateBest(&words, &words, &wordInfos, FancyPrint))
	writeInfos(&(wordInfos), infosFilename)
}
