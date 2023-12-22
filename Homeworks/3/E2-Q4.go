package main

import (
	"fmt"
	"sync"
	"time"
)

// number of philosophers
var numPhil int = 5

// philosopher structure
type Philosopher struct {
	ID                  int
	LeftFork, RightFork *sync.Mutex
}

// Thinking process
func (p Philosopher) think() {
	fmt.Printf("Philosopher %d is thinking...\n", p.ID)
	time.Sleep(time.Duration(p.ID+1) * time.Second)
}

// eating process
func (p Philosopher) eat() {
	p.LeftFork.Lock()
	p.RightFork.Lock()

	fmt.Printf("Philosopher %d is eating...\n", p.ID)
	time.Sleep(time.Duration(p.ID+1) * time.Second)

	p.RightFork.Unlock()
	p.LeftFork.Unlock()
}

// Dining
func (p Philosopher) dine(wg *sync.WaitGroup) {
	defer wg.Done()

	//each philosopher eat at most 3 three times
	for i := 0; i < 3; i++ {
		p.think()
		p.eat()
	}
}

func main() {
	//create forks
	forks := make([]*sync.Mutex, numPhil)
	for i := 0; i < numPhil; i++ {
		forks[i] = &sync.Mutex{}
	}

	//create philosophers
	Philosophers := make([]Philosopher, numPhil)
	for i := 0; i < numPhil; i++ {
		leftFork := forks[i]
		rightFork := forks[(i+1)%numPhil]

		Philosophers[i] = Philosopher{
			ID:        i + 1,
			LeftFork:  leftFork,
			RightFork: rightFork,
		}
	}

	var wg sync.WaitGroup
	//5 goroutines
	wg.Add(numPhil)

	//start dining
	for _, philosopher := range Philosophers {
		go philosopher.dine(&wg)
	}

	//ensure that all goroutines finish their works
	wg.Wait()
	fmt.Println("Dining philosophers have finished.")
}
