package mapper

type Work struct {
	Action    string
	Count     func(text string) int
	CountChan chan int
	List      func(text string) []string
	ListChan  chan []string
}

func Run(text string, queryMap map[string]Work) {
	for _, v := range queryMap {
		if v.Action == "count" {
			v.CountChan <- v.Count(text)
		} else {
			v.ListChan <- v.List(text)
		}
	}
}
