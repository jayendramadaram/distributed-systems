package reduce

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type reducer[T any] struct {
	label       string
	pipe        chan T
	cache       interface{}
	channelType string
	quit        chan bool
}

type Reducer interface {
	Label() string
	Start() error
	Stop()
	Cache() interface{}
}

func New[T any](label string, pipe chan T, channelType string) Reducer {

	reducer := &reducer[T]{
		label:       label,
		pipe:        pipe,
		quit:        make(chan bool),
		channelType: channelType,
	}

	switch channelType {
	case "stringSlice":
		reducer.cache = make([]string, 0)
	case "int":
		reducer.cache = 0
	}
	return reducer
}

func (r *reducer[T]) Start() error {
	var aggregateFunc func(V interface{})

	switch r.channelType {
	case "stringSlice":
		aggregateFunc = func(v interface{}) {
			r.cache = append(r.cache.([]string), v.([]string)...)
		}

	case "int":
		aggregateFunc = func(v interface{}) {
			r.cache = r.cache.(int) + v.(int)
		}

	default:
		return fmt.Errorf("unknown channel type %s", r.channelType)
	}

	for {
		select {
		case <-r.quit:
			if err := r.writeToFile(r.cache); err != nil {
				panic(err)
			}
			return nil
		case v := <-r.pipe:
			aggregateFunc(v)
		}
	}
}

func (r *reducer[T]) writeToFile(data interface{}) error {
	// Get the current date for the filename
	now := time.Now().Format("2006-01-02")
	fileName := "results/output_" + r.label + "_" + now + ".txt"

	// Check the type of the data using reflect
	typeOfData := reflect.TypeOf(data)
	var content string

	switch r.channelType {
	case "stringSlice":
		stringSlice := data.([]string)
		content = strings.Join(stringSlice, "\n")
	case "int":
		intVal := data.(int)
		content = fmt.Sprintf("%d", intVal) // Convert int to string
	default:
		return fmt.Errorf("unsupported data type: %v", typeOfData)
	}

	// Write the content to the file
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Printf("Successfully wrote content to %s\n", fileName)
	return nil
}

func (r *reducer[T]) Stop() {
	r.quit <- true
}

func (r *reducer[T]) Label() string {
	return r.label
}

func (r *reducer[T]) Cache() interface{} {
	return r.cache
}
