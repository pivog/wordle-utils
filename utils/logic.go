package utils

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Custom types to work with goroutines and channels
type syncedVar[T any] struct {
	sync.RWMutex
	str any
}

func Get(variable *syncedVar[any]) any {
	variable.RLock()
	defer variable.RUnlock()
	return variable.str
}

func Set(variable *syncedVar[any], value any) {
	variable.Lock()
	defer variable.Unlock()
	variable.str = value
}

type infoAndWord struct {
	word string
	info float64
}

// GenerateNewPossible Generate the reduced list of possible answers after getting feedback
// from a word (simulating a possibility or playing a word in wordle)
func GenerateNewPossible(word string, feedback string, possibles []string) []string {
	var out []string
	for _, element := range possibles {
		if GetFeedback(word, element) == feedback {
			out = append(out, element)
		}
	}
	return out
}

// GenerateBest Generate the best possible word
func GenerateBest(possibles *[]string, allWords *[]string, wordInfos *sync.Map, printer func(string)) string {
	if len(*possibles) == 1 {
		return (*possibles)[0]
	}

	// If there are two words there is 50% chance for both
	if len(*possibles) == 2 {
		return (*possibles)[0]
	}

	feedbacksMap, feedbacksCounter := getAllFeedbacks()

	bestInfo := &syncedVar[float64]{}
	Set((*syncedVar[any])(bestInfo), float64(0))

	bestWord := &syncedVar[string]{}
	Set((*syncedVar[any])(bestWord), "")

	doneCount := &syncedVar[int]{}
	Set((*syncedVar[any])(doneCount), 0)

	wg := sync.WaitGroup{}

	// Channel to take in processed words and expected information
	infoInChan := make(chan infoAndWord)
	go func() {
		for {
			infoIn := <-infoInChan
			if infoIn.info > Get((*syncedVar[any])(bestInfo)).(float64) {
				Set((*syncedVar[any])(bestInfo), infoIn.info)
				Set((*syncedVar[any])(bestWord), infoIn.word)
			}
		}
	}()

	startTime := time.Now()
	// Start a goroutine for each word to see the expected information
	for _, currentTry := range *allWords {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wordExpectedInfo := getExpectedInfo(currentTry, *possibles, feedbacksMap, feedbacksCounter)
			infoInChan <- infoAndWord{
				word: currentTry,
				info: wordExpectedInfo,
			}
			// Make a local copy not to access a synced variable more than needed
			_doneCount := Get((*syncedVar[any])(doneCount)).(int)
			progressString := "Checked    " +
				fmt.Sprintf("%d", Get((*syncedVar[any])(doneCount)).(int)+1) +
				". " +
				currentTry +
				"    " +
				strings.Trim(fmt.Sprintf("%f", 100*roundFloat(float64(_doneCount)/float64(len(*possibles)), 3)), "0") +
				"%     " +
				strings.Trim(fmt.Sprintf("%f", roundFloat(time.Since(startTime).Seconds(), 3)), "0") +
				"s    " +
				strings.Trim(fmt.Sprintf("%f", roundFloat(float64(_doneCount+1)/time.Since(startTime).Seconds(), 3)), "0") +
				"wps    " +
				formatTime(getEstimatedTime(startTime, _doneCount, len(*possibles)))

			// Print progress
			printer(progressString)
			// Store the expected information
			(*wordInfos).LoadOrStore(currentTry, wordExpectedInfo)
			Set((*syncedVar[any])(doneCount), Get((*syncedVar[any])(doneCount)).(int)+1)
		}()
	}
	wg.Wait()
	return Get((*syncedVar[any])(bestWord)).(string)
}
