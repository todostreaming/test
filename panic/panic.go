package main

import (
	"os/exec"
	"strconv"
	"strings"
	"fmt"
	"os"
	"time"
)

// /usr/local/bin/panic 2>> /usr/local/bin/panic.log  (program will only write to stdout and error files, and stderr in panics will go to panic.log
// crontab -e (root)   * * * * * /usr/local/bin/panic 2>> /usr/local/bin/panic.log

var panico map[string]string

func main(){
	if procsrunning("panic") > 1 {
		fmt.Printf("Hay otra instancia de este programa en RAM\n\n")
		os.Exit(1)
	}
	v := time.Now()
	fmt.Printf("%02d:%02d:%02d\n",v.Hour(),v.Minute(),v.Second())
	fmt.Println("Hello...")
	panico["pepe"]="tio"
	fmt.Println(panico["pepe"])
}

// devuelve en Linux la cantidad de procesos ejecutandose con el nombre name
func procsrunning(name string) (int) {
	exe := fmt.Sprintf("/usr/bin/pgrep %s | /usr/bin/wc -l",name)
	out,_:=exec.Command("/bin/sh","-c",exe).CombinedOutput()
	num,_ := strconv.Atoi(strings.TrimRight(string(out),"\n"))
	return num
}
