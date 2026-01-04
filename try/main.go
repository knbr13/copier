package main

import (
	// "github.com/knbr13/copier"
	"fmt"

	"github.com/jinzhu/copier"
)

func main() {
	// p := Person{
	// 	Age: make(chan int, 4),
	// }
	// p.Age <- 1
	// p.Age <- 2

	// go func() {
	// 	time.Sleep(time.Second)
	// 	close(p.Age)
	// }()

	// var p2 Person

	i, j := 12, 18

	err := copier.Copy(&i, j)
	if err != nil {
		panic(err)
	}
	fmt.Println(i, j)
	// for v := range p.Age {
	// 	fmt.Println("p1:", v)
	// }

	// go func() {
	// 	time.Sleep(time.Second)
	// 	close(p2.Age)
	// }()

	// for v := range p2.Age {
	// 	fmt.Println("p2:", v)
	// }
}

type Person struct {
	Age chan int
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
}
