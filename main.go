package main

import "utils"

func main() {
	if utils.SovlerMode {
		// if the program is in solver mode then it interactively solves the wordle game for the user
		utils.Solver()
	} else {
		// analyses for the best opener
		utils.GenerateBestAndWriteFromStart()
	}
}
