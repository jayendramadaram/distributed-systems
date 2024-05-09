package main

import (
	"fmt"
	"mapreduce"
	"mapreduce/mapper"
	"os"
	"regexp"
	"strings"
)

func printHelp() {
	fmt.Println("Please provide an even number of arguments")
	fmt.Println("Usage \n  go run cmd/main.go input.file <query> \n\n\t <query> space separated instructions key=\"pattern\" \n\t where pattern is space separated commands")
	fmt.Println("\n Example: \n  go run cmd/main.go input.file url=\"/^(https?://)?([da-z.-]+).([a-z.]{2,6})([/w.-]*)*/?$/ count\" pincode=\"^[1-9][0-9]{5}$ list\" ")
}

func regex_util_finder(re *regexp.Regexp) func(text string) []string {
	return func(text string) []string {
		return re.FindAllString(text, -1)
	}
}

func regex_util_count(re *regexp.Regexp) func(text string) int {
	return func(text string) int {
		return len(re.FindAllString(text, -1))
	}
}

func exact_match_finder(match string) func(text string) int {
	return func(text string) int {
		return strings.Count(text, match)
	}
}

func main() {
	// The first value is the path to the executable
	args := os.Args[1:]

	if len(args) < 2 {
		printHelp()
		os.Exit(1)
	}

	queryMap := make(map[string]mapper.Work)

	input_file := args[0]
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			printHelp()
			os.Exit(1)
		}

		command := strings.Split(parts[1], " ")
		if len(command) != 2 {
			printHelp()
			os.Exit(1)
		}

		instruction, action := command[0], command[1]
		fmt.Println("instruction: ", instruction, " action: ", action)

		if instruction == "exact" {
			queryMap[parts[0]] = mapper.Work{
				Action:    "count",
				Count:     exact_match_finder(parts[0]),
				CountChan: make(chan int),
			}

		} else {
			re := regexp.MustCompile(instruction)
			switch action {
			case "count":
				queryMap[parts[0]] = mapper.Work{
					Action:    "count",
					Count:     regex_util_count(re),
					CountChan: make(chan int),
				}
			case "list":
				queryMap[parts[0]] = mapper.Work{
					Action:   "list",
					List:     regex_util_finder(re),
					ListChan: make(chan []string),
				}
			}
		}
	}

	mapreduce.Execute(input_file, queryMap)
}
