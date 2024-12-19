package utils

import (
	"math"
	"strings"
)

// A feedback is a combination of green, yellow and gray letters
// Green is represented with +
// Yellow is represented with =
// Gray is represented with -

// Example:
// If the word to check is "audio"(a common starter) and the answer is "autos", the feedback will be "++--="
// For easier visualization:
// audio
// autos
// ++--=

// GetFeedback Generate the feedback from a word and a possible answer using the wordle rules
func GetFeedback(word string, answer string) string {
	out := [5]string{"-", "-", "-", "-", "-"}
	lettersUsed := ""
	for i := 0; i < 5; i++ {
		if word[i] == answer[i] {
			lettersUsed += string(word[i])
			out[i] = "+"
		}
	}
	for i := 0; i < 5; i++ {
		if out[i] == "+" {
			continue
		}
		if strings.Contains(lettersUsed, string(word[i])) {
			continue
		}
		if strings.Contains(answer, string(word[i])) && string(word[i]) != "+" {
			out[i] = "="
			lettersUsed += string(word[i])
		}
	}
	return strings.Join(out[0:], "")
}

// Generates all possible feedbacks with base 3 counter behavior
func getAllFeedbacks() (map[string]float64, []string) {

	// Get the next item in the list, wrapping back to index 0
	getNextChar := func(characters string, element string) string {
		currentIndex := strings.Index(characters, string(element[0]))
		nextIndex := (currentIndex + 1) % len(characters)
		return characters[nextIndex : nextIndex+1]
	}

	// a map which will map a feedback to the frequency of a feedback
	feedbacks := map[string]float64{}
	var feedbacksIter []string
	letters := "-=+"
	buf := []string{"-", "-", "-", "-", "-"}
	// an ugly loop that works
	for i1 := 0; i1 < 3; i1++ {
		for i2 := 0; i2 < 3; i2++ {
			for i3 := 0; i3 < 3; i3++ {
				for i4 := 0; i4 < 3; i4++ {
					for i5 := 0; i5 < 3; i5++ {
						feedbacks[strings.Join(buf, "")] = 0
						feedbacksIter = append(feedbacksIter, strings.Join(buf, ""))
						buf[4] = getNextChar(letters, buf[4])
					}
					buf[3] = getNextChar(letters, buf[3])
				}
				buf[2] = getNextChar(letters, buf[2])
			}
			buf[1] = getNextChar(letters, buf[1])
		}
		buf[0] = getNextChar(letters, buf[0])
	}
	return feedbacks, feedbacksIter
}

// get the expected information a word will give measured
// in bits of information as defined in information theory
func getExpectedInfo(word string,
	possibles []string,
	feedbacks map[string]float64,
	feedbacksCounter []string) float64 {

	// Make a new map to populate with frequencies
	// It is important to deep copy not to mess with the original map because
	out := map[string]float64{}
	for k, v := range feedbacks {
		out[k] = v
	}

	expectedWordInfo := float64(0)

	// Check every word that could be the answer and populate frequencies of each feedback
	for _, possibleAnswer := range possibles {
		out[GetFeedback(word, possibleAnswer)] = out[GetFeedback(word, possibleAnswer)] + 1
	}

	// Calculate the information of each feedback from the frequencies
	for _, feedback := range feedbacksCounter {
		if out[feedback] == 0 {
			continue
		}
		expectedWordInfo += math.Log2(float64(len(possibles))/out[feedback]) / float64(len(possibles))
		// An alternative way to calculate the information, the values will not be calculated in bits,
		// but they would have the same comparative properties as bits
		// The advantage would be replacing logarithms with multiplication
		//expectedWordInfo *= out[feedback]
	}

	return expectedWordInfo
}
