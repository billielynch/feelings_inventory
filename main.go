package main

/*
TODO:
- better return (things, err) and handling
- clean up the returns
- look at what are conventions for cases in various things and comments, defacto style internal, the way they implement them.
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"strings"
)

const dataLocation = "./feelings_data.json"

type feelingsGroup struct {
	Feelings []string `json:"feelingsList"`
	Name     string   `json:"groupName"`
}

type resultsGroup struct {
	FeelingsHad []string
	Name        string
	Saturation  float32
}

func main() {
	slog.Info("Started Feeling Inventory")
	feelingsGroups, _ := loadData(dataLocation)
	slog.Debug("Collecting User Inputs")
	results := collectFeels(feelingsGroups)
	slog.Debug("Displaying Results")
	displayShortResults(results)
	slog.Info("Ended Feeling Inventory")
}

// Loading the json data that has the feelings to prompt with
func loadData(absolutePath string) (feelingGroups []feelingsGroup, err error) {

	slog.Debug(absolutePath)

	jsonFile, err := os.Open(absolutePath)
	if err != nil {
		slog.Error(err.Error())
		return feelingGroups, err
	}
	slog.Debug("Opened datafile")
	defer jsonFile.Close() // what is the scope of defer

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		slog.Error(err.Error())
		return feelingGroups, err
	}

	json.Unmarshal(byteValue, &feelingGroups) // what does `&` do here == reference and dereference operator &/*
	slog.Debug("Loaded Data")
	return feelingGroups, err
}

func promptUser(question string) (felt bool) {
	/*
		Prints the given question string with a question mark to the terminal
		prompts for user input in the terminal and accepts it when a newline is
		entered.

		Given any input returns true, if no input is given returns false.
	*/

	fmt.Printf("%s? ", question)
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		slog.Error("An error occured while reading input", err)
		return false
	}

	userInput = strings.TrimSuffix(userInput, "\n")
	if userInput != "" {
		return true
	}
	return false
}

func generateRandomArrayOfIndexes(maxNumber int) []int {
	/*
		A helper for a wacky (?) way to iterate randomly though a slice exactly once
		generates an array to be used as the indexes.
	*/

	indexArray := make([]int, 0, maxNumber)
	for i := 0; i < maxNumber; i++ {
		indexArray[i] = i
	}

	rand.Shuffle(maxNumber, func(i, j int) {
		indexArray[i], indexArray[j] = indexArray[j], indexArray[i]
	})
	// figure out wth this is actually doing, this is about not supporting generics outright, looking at how sort
	// they do support generics now, there might be a non-experimental version of this
	return indexArray
}

// Prompts user for input and records result
func collectFeels(feelingsGroups []feelingsGroup) (results []resultsGroup) {
	numGroups := len(feelingsGroups)
	randomisedIndexes := generateRandomArrayOfIndexes(numGroups)

	for i := 0; i < numGroups; i++ {
		j := randomisedIndexes[i]
		result := promptForGroup(feelingsGroups[j])
		results = append(results, result)
	}

	return results
}

func promptForGroup(group feelingsGroup) (result resultsGroup) {
	/*
		Implements the prompting and result gathering for each group of emotions
	*/

	var feelingsHad []string
	numFeels := len(group.Feelings)
	randomisedIndexes := generateRandomArrayOfIndexes(numFeels)

	for i := 0; i < numFeels; i++ {
		j := randomisedIndexes[i]
		emotion := group.Feelings[j]
		result := promptUser(emotion)
		if result == true {
			feelingsHad = append(feelingsHad, emotion)
		}
	}

	result = resultsGroup{
		FeelingsHad: feelingsHad,
		Name:        group.Name,
		Saturation:  float32(len(feelingsHad)) / float32(len(group.Feelings)),
	}

	return result
}

func displayVerboseResults(results []resultsGroup) {
	/*
		Prints a verbose version of each given results group in the slice
	*/

	for i := range results {
		resultsGroup := results[i]
		fmt.Println("----------------")
		fmt.Println(strings.ToUpper(resultsGroup.Name))
		saturationFormatted := fmt.Sprintf("saturation : %.2f", resultsGroup.Saturation)
		fmt.Println(saturationFormatted)
		feelingsHad := resultsGroup.FeelingsHad
		for j := range feelingsHad {
			emotion := feelingsHad[j]
			fmt.Println(emotion)
		}
	}
}

func displayShortResults(results []resultsGroup) {
	/*
		Prints a short version of each given results group in the slice
	*/
	for i := range results {
		resultsGroup := results[i]
		saturationFormatted := fmt.Sprintf("%s %.2f", resultsGroup.Name, resultsGroup.Saturation)
		fmt.Println(saturationFormatted)
	}

}
