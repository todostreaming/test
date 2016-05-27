package main

import (
	"net/http"
	"log"
	"bufio"
	"strings"
	"fmt"
	"sync"
)

type M3U8pls struct {
	m3u8							string
	Targetdur						float64
	Mediaseq						int64
	Segment							[]string
	Duration 						[]float64
	mu_pls							sync.Mutex
}

func M3U8playlist(m3u8 string) *M3U8pls {
	m3u := &M3U8pls{}
	m3u.mu_pls.Lock()
	defer m3u.mu_pls.Unlock()
	
	m3u.m3u8 = m3u8

	return m3u
}

func (m *M3U8pls) Parse(){
	m.mu_pls.Lock()
	m.Targetdur = 0.0
	m.Mediaseq = 0
	m.Segment = []string {}
	m.Duration = []float64 {}
	m.mu_pls.Unlock()

	m.analyzem3u8()
}

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
//	m3u8 := "http:///streamrus/stream.m3u8"
	m3u8 := "http:///radiovida/mobile/playlist.m3u8"
	m3u8pls := M3U8playlist(m3u8)
	m3u8pls.Parse()
	fmt.Printf("Targetdur=%.2f\nMediaseq=%d\n",m3u8pls.Targetdur,m3u8pls.Mediaseq)
	fmt.Printf("Segments= %v\nDuration= %v\n",m3u8pls.Segment,m3u8pls.Duration)
}

func (m *M3U8pls) analyzem3u8(){
	substr := ""
	issubstr := false
	m.mu_pls.Lock()
	m3u8 := m.m3u8
	m.mu_pls.Unlock()
	resp, err := http.Get(m3u8)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil{
			break
		}
		line = strings.TrimRight(line,"\n")
		if strings.Contains(line,".m3u8"){
			substr = substream(m3u8, line)
			issubstr = true
			break
		}
		if strings.Contains(line,"#EXT-X-TARGETDURATION:"){
			var targetdur float64
			fmt.Sscanf(line,"#EXT-X-TARGETDURATION:%f",&targetdur)
			m.mu_pls.Lock()
			m.Targetdur = targetdur			
			m.mu_pls.Unlock()
		}
		if strings.Contains(line,"#EXT-X-MEDIA-SEQUENCE:"){
			var mediaseq int64
			fmt.Sscanf(line,"#EXT-X-MEDIA-SEQUENCE:%d",&mediaseq)
			m.mu_pls.Lock()
			m.Mediaseq = mediaseq			
			m.mu_pls.Unlock()
		}
		if strings.Contains(line,"#EXTINF:"){
			var extinf float64
			fmt.Sscanf(line,"#EXTINF:%f,",&extinf)
			m.mu_pls.Lock()
			m.Duration = append(m.Duration,extinf)
			m.mu_pls.Unlock()
		}
		if strings.Contains(line,".ts"){
			m.mu_pls.Lock()
			m.Segment = append(m.Segment,substream(m3u8, line))
			m.mu_pls.Unlock()
		}
		//fmt.Printf("1)=>[%s]<=\n",line)
	}
	resp.Body.Close()
	if issubstr {
		resp, err := http.Get(substr)
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil{
				break
			}
			line = strings.TrimRight(line,"\n")
			if strings.Contains(line,"#EXT-X-TARGETDURATION:"){
				var targetdur float64
				fmt.Sscanf(line,"#EXT-X-TARGETDURATION:%f",&targetdur)
				m.mu_pls.Lock()
				m.Targetdur = targetdur			
				m.mu_pls.Unlock()
			}
			if strings.Contains(line,"#EXT-X-MEDIA-SEQUENCE:"){
				var mediaseq int64
				fmt.Sscanf(line,"#EXT-X-MEDIA-SEQUENCE:%d",&mediaseq)
				m.mu_pls.Lock()
				m.Mediaseq = mediaseq			
				m.mu_pls.Unlock()
			}
			if strings.Contains(line,"#EXTINF:"){
				var extinf float64
				fmt.Sscanf(line,"#EXTINF:%f,",&extinf)
				m.mu_pls.Lock()
				m.Duration = append(m.Duration,extinf)
				m.mu_pls.Unlock()
			}
			if strings.Contains(line,".ts"){
				m.mu_pls.Lock()
				m.Segment = append(m.Segment,substream(substr, line))
				m.mu_pls.Unlock()
			}

			//fmt.Printf("2)=>[%s]<=\n",line)
		}
		resp.Body.Close()
	}
}

func substream(m3u8, sub string) string {
	var substream string
	
	m3u8 = m3u8[7:] // quito http://
	substream = "http://"
	parts := strings.Split(m3u8,"/")
	for _,v := range parts {
		if strings.Contains(v,".m3u8"){
			substream = substream + sub
			break
		}else{
			substream = substream + v + "/"
		}
	}
	
	return substream
}

