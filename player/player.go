package main

import (
	"github.com/todostreaming/hlsplay"
	"time"
	"runtime"
	"fmt"
)

func main(){
	settings := map[string]string {
		"overscan"		:		"0",
		"x0"			:		"0",
		"y0"			:		"0",
		"x1"			:		"719",
		"y1"			:		"575",
		"vol"			:		"0",	
	}
//	m3u8 := "http:///streamrus/stream.m3u8"
	m3u8 := "http:///radiovida/mobile/playlist.m3u8"
//	m3u8 := "http:///radiovida/livestream/playlist.m3u8"
	
	hls := hlsplay.HLSPlayer(m3u8, "/var/segments/", settings)
	hls.Run()
	time.Sleep(10 * time.Second)
	hls.Volume(true)
	stat := hls.Status()
	fmt.Println(stat.OMXStat)
	time.Sleep(10 * time.Second)
	hls.Volume(false)
	stat = hls.Status()
	fmt.Println(stat.OMXStat)
	for {
		time.Sleep(1 * time.Second)
		stat = hls.Status()
		fmt.Println(stat.OMXStat)
		runtime.Gosched()
	}	
}

