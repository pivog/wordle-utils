package main

import (
	"sync"
	"testing"
	"utils"
)

func simulateGame(opener string, answer string, words []string) int {
	if opener == answer {
		return 1
	}
	words = utils.GenerateNewPossible(opener, utils.GetFeedback(opener, answer), words)

	wordInfos := sync.Map{}
	turns := 2
	for turn := 2; turn <= 6; turn++ {
		turns = turn
		bestWord := utils.GenerateBest(&words, &words, &wordInfos, func(s string) {})
		if bestWord == answer {
			break
		}
		words = utils.GenerateNewPossible(bestWord, utils.GetFeedback(bestWord, answer), words)
	}
	return turns
}

func TestBenchTurns(t *testing.T) {
	opener := "tares"
	words := utils.LoadWords()

	avgTurns := float64(0)

	println("started benchmark")
	for _, answer := range words {
		avgTurns += float64(simulateGame(opener, answer, words)) / float64(len(words))
	}
	print("Takes ", avgTurns, " turns on average")
}
