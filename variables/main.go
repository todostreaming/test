package main

import (
	"fmt"
	"os"
)

func main(){
	for k,v := range os.Environ() {
		fmt.Printf("[%d] = %s\n",k,v)
	}
		
}
