package main

import (
	"fmt"
	"github.com/todostreaming/fifo"
)

type Segment struct {
	Name string
	Dur  float64
}

func main() {
	ts_queue := fifo.NewQueue()
	ts_queue.Add(&Segment{"stream1.ts", 10.56})
	ts_queue.Add(&Segment{"stream2.ts", 9.02})
	ts_queue.Add(&Segment{"stream3.ts", 11.00})

	for q := ts_queue.Next(); q != nil; q = ts_queue.Next() {
		p := q.(*Segment)
		fmt.Printf("TS: %s \t %.2f\n", p.Name, p.Dur)
	}

	/*
		fmt.Printf("1)Len=%d\n",ts_queue.Len())
		ts_queue.Add("segment1.ts")
		ts_queue.Add("segment2.ts")
		ts_queue.Add("segment3.ts")

		for q := ts_queue.Next(); q != nil; q = ts_queue.Next() {
			fmt.Println(q)
		}
	*/
	//	total := ts_queue.Len()
	//	for q:=0; q < total; q++ {
	//		fmt.Println(ts_queue.Next())
	//	}
}
