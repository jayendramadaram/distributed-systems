package remus

import "sync"

type VM struct {
	buffer []int64
	mu     sync.RWMutex
}

func NewVM(size int) *VM {
	return &VM{
		buffer: make([]int64, size),
		mu:     sync.RWMutex{},
	}
}

func (vm *VM) Get(index int) int64 {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.buffer[index]
}

func (vm *VM) Set(index int, value int64) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.buffer[index] = value
}

func (vm *VM) Add(index1 int, index2 int, storeIndex int) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.buffer[storeIndex] = vm.buffer[index1] + vm.buffer[index2]
}

func (vm *VM) Sub(index1 int, index2 int, storeIndex int) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.buffer[storeIndex] = vm.buffer[index1] - vm.buffer[index2]
}

func (vm *VM) Mul(index1 int, index2 int, storeIndex int) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.buffer[storeIndex] = vm.buffer[index1] * vm.buffer[index2]
}

// Note do denominator 0 check before calling this function
func (vm *VM) Div(index1 int, index2 int, storeIndex int) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.buffer[storeIndex] = vm.buffer[index1] / vm.buffer[index2]
}
