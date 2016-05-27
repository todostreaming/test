package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
)

func main(){
	m3u8 := "http:///streamrus/stream.m3u8"
//	m3u8 := "http:///radiovida/mobile/playlist.m3u8"
	resp, err := http.Get(m3u8)
	if err != nil {
		fmt.Println("Cannot connect")
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("No M3u8")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read")
		return
	}
	fmt.Println(string(body))
}
