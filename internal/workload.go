package internal

import (
	"errors"
	"fmt"
)

type Workload struct {
	Name     string
	CPU      int
	Memory   int // Memory in MB
	Priority int
}

func (w *Workload) Validate() error {
	if w.Name == "" {
		return errors.New("workload name cannot be empty")
	}
	if w.CPU <= 0 {
		return errors.New("CPU must be greater than 0")
	}
	if w.Memory <= 0 {
		return errors.New("memory must be greater than 0")
	}
	if w.Priority <= 0 {
		return errors.New("priority must be greater than 0")
	}
	return nil
}

func (w *Workload) String() string {
	return fmt.Sprintf("Workload{Name: %s, CPU: %d, Memory: %dMB, Priority: %d}", 
		w.Name, w.CPU, w.Memory, w.Priority)
}