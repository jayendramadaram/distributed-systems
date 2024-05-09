package remus

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Role int

const (
	Primary Role = iota
	Replica
)

type Node struct {
	vm    *VM
	role  Role
	index int
	// resolver Client
}

func NewNode(vmSize int, role Role, index int) *Node {
	return &Node{
		vm:    NewVM(vmSize),
		role:  role,
		index: index,
	}
}

// - buffer <--
// - promote <--
// - state + heartbeat -->
// - register -->
func (node *Node) Run() {
	r := gin.Default()

	// todo: check origin
	r.POST("/buffer", func(c *gin.Context) {
		// do computations
	})

	r.GET("/promote", func(c *gin.Context) {
		// update role
		// keep sending latest states
	})

	if err := r.Run(fmt.Sprintf(":800%d", node.index)); err != nil {
		panic(err)
	}

}
