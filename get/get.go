package main

import (
	"github.com/todostreaming/fifo"
	"fmt"
)

func main(){
	cola := fifo.NewQueue()
	fmt.Println("1.-",cola.Len())
	cola.Add(10)
	fmt.Println("2.-",cola.Len())
	cola = fifo.NewQueue()
	fmt.Println("3.-",cola.Len())
}

