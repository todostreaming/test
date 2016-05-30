package main

import (
	"fmt"
	"github.com/todostreaming/m3u8pls"
	"time"
)

/*
  HTTP/1.1 200 OK
  Date: Fri, 27 May 2016 07:14:18 GMT
  Content-Type: video/MP2T
  Accept-Ranges: bytes
  Server: FlashCom/3.5.7
  Cache-Control: no-cache
  Content-Length: 876080
*/

func main(){
//	m3u8 := "http://190.7.56.118/streamrus/stream.m3u8"
	m3u8 := "http://orion.comelson.es/radiovida/mobile/playlist.m3u8"
//	m3u8 := "http://pablo002.todostreaming.es/radiovida/livestream/playlist.m3u8"
	m3u8pls := m3u8pls.M3U8playlist(m3u8)
	m3u8pls.Parse()
	if m3u8pls.Ok { 
		fmt.Printf("Targetdur=%.2f\nMediaseq=%d\n",m3u8pls.Targetdur,m3u8pls.Mediaseq)
		fmt.Printf("Segments= %v\nDuration= %v\n",m3u8pls.Segment,m3u8pls.Duration)
	}else{
		fmt.Println("No es accesible la URL M3U8")
	}
	time.Sleep(time.Duration(m3u8pls.Targetdur) * time.Second)
	m3u8pls.Parse()
	if m3u8pls.Ok { 
		fmt.Printf("Targetdur=%.2f\nMediaseq=%d\n",m3u8pls.Targetdur,m3u8pls.Mediaseq)
		fmt.Printf("Segments= %v\nDuration= %v\n",m3u8pls.Segment,m3u8pls.Duration)
	}else{
		fmt.Println("No es accesible la URL M3U8")
	}
}
