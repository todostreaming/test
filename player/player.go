package main

import (
	"bufio"
	"os/exec"
	"strings"
	"fmt"
	"runtime"
	"time"
)

//var mediareader *bufio.Reader
var mediawriter *bufio.Writer
var omx_exe, avconv_exe *exec.Cmd


var settings = make(map[string]string)

func main(){
	settings = map[string]string{
		"overscan"		:		"1",
		"x0"			:		"0",
		"y0"			:		"0",
		"x1"			:		"719",
		"y1"			:		"575",
		"vol"			:		"1",	
	}
	
	go restamper()
	go player()
	time.Sleep(10 * time.Second)
	mediawriter.WriteByte('+'); mediawriter.Flush();
	time.Sleep(10 * time.Second)
	mediawriter.WriteByte('+'); mediawriter.Flush();
	for {
		//time.Sleep(5 * time.Second)
		runtime.Gosched()
	}	
}

func player(){
	
	cmdline := "/usr/bin/omxplayer -s -o both --no-osd -b /tmp/fifo2"
	for{
		omx_exe = exec.Command("/bin/sh","-c",cmdline)
		stderrRead,_ := omx_exe.StderrPipe()
		mediareader := bufio.NewReader(stderrRead)
		stdinWrite,_ := omx_exe.StdinPipe()
		mediawriter = bufio.NewWriter(stdinWrite)
		fmt.Println(cmdline)
		omx_exe.Start()
		for {
			line,err := mediareader.ReadString('\n')
			if err != nil {
				fmt.Println("Salimos de omxplayer")
				break;
			}
			line=strings.TrimRight(line,"\n")
			if strings.Contains(line,"Comenzando...") {
				fmt.Println("OMXPlayer Ready...")
			}
			if strings.Contains(line,"Time:") {
				fmt.Println("[omx]",line)
			}
			runtime.Gosched()
		
		}
		killall("omxplayer omxplayer.bin dbus-daemon")
		if omx_exe.Process != nil { // si el proceso encoder_exe existe en memoria
			omx_exe.Wait()
		}
		
		
	}	
}

func restamper(){
	cmdline := "/usr/bin/ffmpeg -y -f mpegts -re -i /tmp/fifo1 -f mpegts -acodec copy -vcodec copy /tmp/fifo2"
	for {
		avconv_exe = exec.Command("/bin/sh","-c",cmdline)
		stderrRead,_ := avconv_exe.StderrPipe()
		mediareader := bufio.NewReader(stderrRead)
		fmt.Println(cmdline)
		avconv_exe.Start()
		for {
			line,err := mediareader.ReadString('\n')
			if err != nil {
				fmt.Println("Salimos de avconv")
				break;
			}
			line=strings.TrimRight(line,"\n")
			if strings.Contains(line,"built") {
				fmt.Println("AVConv Ready...")
			}
			if strings.Contains(line,"frame=") {
				fmt.Println("[ffmpeg]",line)
			}
			runtime.Gosched()
		
		}
		killall("ffmpeg")
		if avconv_exe.Process != nil { // si el proceso encoder_exe existe en memoria
			avconv_exe.Wait()
		}
				
		
	}	
}

 
// killall("bmdcapture avconv")
func killall(list string){
	prog := strings.Fields(list)
	for _,v := range prog {
		exec.Command("/bin/sh","-c","kill -9 `ps -A|awk '/"+v+"/{print $1}'`").Run()
	}
}


