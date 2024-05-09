package mapreduce

import (
	"mapreduce/mapper"
	"mapreduce/reduce"
	"os"
	"sync"
)

// split input text stream into M chunks
// M == WorkerCount

// N == Query Output Elements

// pass M input chunks to M workers
// along with generic looking functions to mapper(Workers) to work on text
// Workers send processed data into N channels which are received by Reducers

// Mapping and Shuffle would be done
// reducers aggregate the data and wait for all Workers to complete once complete reducers dump data and close channels

// life Cycle
// - Split
// - Assign Workers and work
// - process and Shuffle
// - reduce
// - dump to output.txt

func splitStringEqualParts(str string, n int) []string {
	if n <= 0 {
		return nil
	}
	// Calculate the size of each piece with ceiling division
	partSize := (len(str) + n - 1) / n
	parts := make([]string, n)
	var start int
	for i := range parts {
		end := start + partSize
		if end > len(str) {
			end = len(str)
		}
		parts[i] = str[start:end]
		start = end
	}
	return parts
}

func Execute(inputFile string, queryMap map[string]mapper.Work) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	reducers := make([]reduce.Reducer, 0)

	rdcr_wg := sync.WaitGroup{}
	rdcr_wg.Add(len(queryMap))
	for k, v := range queryMap {
		var reducer reduce.Reducer
		if v.Action == "count" {
			reducer = reduce.New(k, v.CountChan, "int")
		} else {
			reducer = reduce.New(k, v.ListChan, "stringSlice")
		}

		reducers = append(reducers, reducer)

		go func() {
			defer rdcr_wg.Done()
			reducer.Start()
		}()
	}

	wg := sync.WaitGroup{}
	for _, text := range splitStringEqualParts(string(data), WorkerCount) {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			mapper.Run(text, queryMap)
		}(text)
	}
	wg.Wait()

	for _, reducer := range reducers {
		reducer.Stop()
	}

	rdcr_wg.Wait()

	return nil
}
