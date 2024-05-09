package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// listen for commands from stdin
// Possible Commands

/*
	start <number of nodes `n`> <primary host index `p`> <vm space>
		- p <= n ??
		- if already running! eat five star do nothing.
	shutdown <index> | all
		- if not running!
	vm Actions
		- Set <index> <num>
		- Get <index>
		- (add | sub | mul | div) <index1> <index2> <storeIndex>
*/

// start resolver
// start backup servers which register as backup server with resolver
// start primary server which registers as primary server with

func printHelp() {
	fmt.Println(`
Help:
		start <number of nodes> <primary host index > <vm space>
		- p <= n ??
		- if already running! eat five star do nothing.
	shutdown <index> | all
		- if not running!
	vm Actions
		- Set <index> <num>
		- Get <index>
		- (add | sub | mul | div) <index1> <index2> <storeIndex>

	`)
}

func main() {
	handler := NewHandler()

	scanner := bufio.NewScanner(os.Stdin)
	printHelp()

	print("> ")
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		handler.handleInput(input)
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}
