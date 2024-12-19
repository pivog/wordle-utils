package utils

import (
	"fmt"
	"github.com/nathan-fiscaletti/consolesize-go"
	"math"
	"strings"
	"time"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func rightPad(str string) string {
	cols, _ := consolesize.GetConsoleSize()
	if cols < len(str) {
		return str
	}
	return str + strings.Repeat(" ", cols-len(str))
}

func FancyPrint(str string) {
	print("\r\033[k" + rightPad(str))
}

func formatTime(t time.Duration) string {
	if t.Seconds() >= 60 {
		return fmt.Sprintf("%fm %fs", math.Floor(t.Minutes()), math.Floor(t.Seconds()))
	}
	return fmt.Sprintf("%d.%ds", int(math.Floor(t.Seconds())), t.Milliseconds()%1000)
}

// in milliseconds
func getEstimatedTime(startTime time.Time, count int, possibleCount int) time.Duration {
	return time.Duration(float64(time.Since(startTime).Nanoseconds()) / float64(count+1) * float64(possibleCount-count-1))
}
