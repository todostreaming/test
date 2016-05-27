package main

import (
	"github.com/todostreaming/fifo"
	"fmt"
)

func main(){
	ts_queue := fifo.NewQueue()
	fmt.Printf("1)Len=%d\n",ts_queue.Len())
	ts_queue.Add("segment1.ts")
	ts_queue.Add("segment2.ts")
	ts_queue.Add("segment3.ts")

	for q := ts_queue.Next(); q != nil; q = ts_queue.Next() {
		fmt.Println(q)
	}
//	total := ts_queue.Len()
//	for q:=0; q < total; q++ {
//		fmt.Println(ts_queue.Next())
//	}
}

