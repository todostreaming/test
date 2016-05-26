package main

import (
	"time"
	"fmt"
)

func main(){
	fmt.Println("Empiezo a contar")
	c := Tick( 1 * time.Second ) // c es un channel de tipo Time
	// recibe del channel a v hasta que el canal se cierre con close(chan)
	// v <- c
	for v :=range c { // recibe
		fmt.Printf("%02d:%02d:%02d\r",v.Hour(),v.Minute(),v.Second())
	}
}

// el valor devuelto es un canal time.Time que envia el Time actual
// pero aunque lo envía por el canal, como valor devuelto es visto desde fuera
// de la función y por tanto es un canal por el que se reciben cosas en main
// por tanto es un receive-only channel  <-chan que puede ser leido por range (v <- c)
func Tick(t time.Duration) <-chan time.Time { // desde main es un receive-only
	c := make(chan time.Time) // bidireccional

	go func(){
		for {
			time.Sleep(t)
			c <- time.Now() // 1 (send-only)
		}
	}()
	
	return c // parado hasta 1
}
