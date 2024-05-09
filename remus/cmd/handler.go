package main

import (
	"fmt"
	"strings"
	"sync"
)

type Handler struct {
	IsStarted bool
	mu        sync.RWMutex
	// resolver
}

func NewHandler() *Handler {
	return &Handler{
		IsStarted: false,
	}
}

func (h *Handler) handleInput(input string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	args := strings.Split(input, " ")

	switch args[0] {
	case "help":
		printHelp()

	case "start":
		if h.IsStarted {
			fmt.Println("Already started")
			return
		}
		h.IsStarted = true

	case "shutdown":
	case "set":
	case "get":
	case "add":
	case "sub":
	case "mul":
	case "div":
	}
}
