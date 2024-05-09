package remus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type command []string

type State struct {
	VmState []int64   `json:"vmState"`
	Cmds    []command `json:"cmds"`
}

type Resolver struct {
	StateBuffer        []State
	currentPrimaryHost int
	// node client
}

func New(primaryHost int) *Resolver {
	return &Resolver{
		currentPrimaryHost: primaryHost,
		StateBuffer:        make([]State, 0),
	}
}

// gin server
// - register <-- backup / primary host
// - heartbeat <-- primary host
// - promote --> backup host
// - buffer --> backup host
func (r *Resolver) Run() {

	s := gin.Default()

	// todo: check origin
	// do we need to make sure it is concurrent safe according to remus?
	s.POST("/heartbeat", func(c *gin.Context) {
		req := State{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}

		r.StateBuffer = append(r.StateBuffer, req)
	})

	s.GET("/register", func(c *gin.Context) {
		// add new replica to set
	})

	go func() {
		if err := s.Run(fmt.Sprintf(":8000")); err != nil {
			panic(err)
		}
	}()

	// 50ms timer
	timer_50ms := time.NewTicker(50 * time.Millisecond)
	timer_100ms := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-timer_50ms.C:
			// post buffer to all slave nodes

		case <-timer_100ms.C:
			// elect new slave as primary
		}
	}
}

type ResolverClient struct {
	port string
}

func NewResolverClient(port string) *ResolverClient {
	return &ResolverClient{
		port: port,
	}
}

func (r *ResolverClient) PostBuffer(state []State) {
	if err := SendRequest("POST", state, nil); err != nil {
		// todo: handle error | use proper logger
		panic(err)
	}
}

func (r *ResolverClient) Register(nodeIndex int) {

}

func SendRequest(method string, params, respType interface{}) error {

	jsonData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// url := fmt.Sprintf("http://%s:%s%s", c.Ip, c.Port, c.Path)
	req, err := http.NewRequest(method, "http://localhost:8000", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send request: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&respType)
}
