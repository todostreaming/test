package main

import (
	"fmt"
)

func main() {
	var a,b,c byte
	
	a = 0x20
	b = 0x78 // symmetrical keybyte
	c = a ^ b
	fmt.Printf("Encode:%X XOR [%X] = %X\n",a,b,c)
	
	fmt.Printf("Decode:%X XOR [%X] = %X\n",c,b, c^b)
}
