package main

import (
	"fmt"
	"time"

	"github.com/knbr13/copier"
)

func main() {
	p := Person{
		Age: make(chan int, 4),
	}
	p.Age <- 1
	p.Age <- 2

	go func() {
		time.Sleep(time.Second)
		close(p.Age)
	}()

	var p2 Person

	err := copier.DeepCopyStruct(&p2, p)
	if err != nil {
		panic(err)
	}
	// for v := range p.Age {
	// 	fmt.Println("p1:", v)
	// }

	go func() {
		time.Sleep(time.Second)
		close(p2.Age)
	}()

	for v := range p2.Age {
		fmt.Println("p2:", v)
	}
}

type Person struct {
	Age chan int
}

// type Address struct {
// 	Street string
// 	City   string
// 	State  string
// }
