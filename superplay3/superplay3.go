package main

import ("fmt";"os";"strings";"bufio";"os/exec";"time";"strconv";"net/http";"io/ioutil")

var crontab = []string{"04:30:0","07:30:1","08:55:0","13:30:1","17:30:0","23:00:-1"}
var server = "http://10.26.4.121" // :8080 para Dinosol solamente
var media = "C:\\RSPS\\"
var winamp = `C:\Archivos de programa\RSPSL\RadioSupersol.exe`
var playing string
var kk int // indice de la playlist actual

func main(){
	if procsrunning("superplay3") > 1 {
		fmt.Printf("Hay otra instancia de este programa en RAM\n\n")
		os.Exit(1)
	}
	stop() // playing="-1"
	go getcrontab()
	for {
		play := program() // establece kk de la playlist activa ahora mismo
		if playing=="-1"{ // stopped
			if play != "-1" {
				playing=play
				fmt.Printf("[run(1)] - PLAYLIST = %d \tTIME= %s ========>>>>>\n",kk,time.Now().String())
				go doplay(fmt.Sprintf("%splaylist%d.m3u",media,kk),kk)
			}
		}else{ // playing some 0 or 1
			if playing != play {
				stop()
				playing=play
				if play != "-1" {
					fmt.Printf("[run(2)] - PLAYLIST = %d \tTIME= %s ========>>>>>\n",kk,time.Now().String())
					go doplay(fmt.Sprintf("%splaylist%d.m3u",media,kk),kk)
				}
			}
		}		
		time.Sleep(1 * time.Second) // 1 segundo
	}
}

func doplay(file string, index int){
	var tipo string
	hostname,_:=os.Hostname(); hostname=gethostname(hostname)
	fileinfo,err := os.Stat(file)
	if err != nil || fileinfo.Size() == 0 {
		fmt.Printf("El fichero %s no existe\n\n",file)
	}
	for{
		m3u,m3uError := os.OpenFile(file,os.O_RDONLY ,0)
		defer m3u.Close()
		if m3uError == nil { // el fichero existe y se abre
			m3uReader := bufio.NewReader(m3u)
			for {
				linea,errorLinea := m3uReader.ReadString('\n')
				if errorLinea != nil { break; } // final de fichero
				linea = strings.Trim(linea,"\n");linea = strings.Trim(linea,"\r")
				if linea != "" && linea != "#EXTM3U" { 
					if !file_exists(media + linea) {
						continue
					}else{
						dur:=mp3duration(media + linea)
						if dur > 0 {
							cmdn := exec.Command(winamp,media + linea)
							go cmdn.Run()
							//fmt.Println(winamp,media + linea) // --
							fmt.Printf("... Playing %s during %d seconds ...\n",linea,dur)
							if strings.TrimLeft(linea,"20") != linea { tipo="publi" }else{ tipo="music" }
							//fmt.Println(server+"/megafonia2/playing.php?type="+tipo+"&name="+linea+"&tienda="+hostname) // --
							go geturl(server+":8080/megafonia2/playing.php?type="+tipo+"&name="+linea+"&tienda="+hostname, 3)
							time.Sleep(dur * time.Second)
							// debo de pararme ya o seguir ???
							if index != kk || playing == "-1" { fmt.Printf("... Exiting from thread: %d\n",index);return }
						}else{
							fmt.Printf("Skipping %s ...\n",linea)
						}
					}
				}
			}
		}else{
			fmt.Printf("El fichero %s no se ha podido abrir correctamente\n\n",file)
			os.Exit(1)
		}
	}
}

// devuelve la duracion de un fichero mp3 en segundos (0 significa error)
func mp3duration(filename string)(duration time.Duration){
	var bitrate = [16]int{ 0, 32000, 40000, 48000, 56000, 64000, 80000, 96000, 112000, 128000, 160000, 192000, 224000, 256000, 320000, 0 }
	
	var b2 = []byte{0,0}
	var b = []byte{0}
	fr,err := os.Open(filename)
	if err == nil {
		i:=0
		// Vamos a buscar la zona FF FB
		for i=0;!(b2[0] == '\xFF' && b2[1] == '\xFB');i+=2 {
			fr.ReadAt(b2,int64(i))
		}
		fr.ReadAt(b,int64(i))
		fmt.Printf("Byte3: %x\n",b[0])
		br := int(b[0]/16)
		fileinfo,err := os.Stat(filename)
		if err != nil || fileinfo.Size() == 0 {
			duration = 0
		}else{
			duration = time.Duration((8 * fileinfo.Size() / int64(bitrate[br])) + 1)
			fmt.Printf("Bitrate: %d kbps\n",bitrate[br])
		}
	}else{
		duration = 0
	}
	fr.Close()
	return
}

// funcion que reintenta hasta count veces bajarse la url en un string res completo (sin fallos)
func geturl(url string, count int)(res string){
	i:=0; e:=0;
	var resp *http.Response;
	var err error;
	
	if count < 1 {count = 1}
	for(i < count){
		i++
		fmt.Printf("[geturl(1)] - GET URL: %s\n",url)
		resp, err = http.Get(url)
		if err != nil {
			e++
			fmt.Printf("[geturl(2)] - Error GETTING URL: %s\n",url)
			continue
		}else {
			break
		}
	}
	if e >= count {
		fmt.Printf("[geturl(3)] - IMPOSSIBLE GETTING URL(SERVER DOWN): %s\n",url)
		res = ""	
	}else{
		body, _ := ioutil.ReadAll(resp.Body)
		res = string(body)
		resp.Body.Close()
		if strings.Contains(res,"</") || strings.Contains(res,"404") {
			fmt.Printf("[geturl(5)] - NO EXISTE URL(SERVER UP): %s\n",url) 
			res = "" 
		}else{
			fmt.Printf("[geturl(4)] - URL OK %s\n",url)
		}
	}
	return res
}

func file_exists(dir string) bool {
    info, err := os.Stat(dir)
    if err != nil {
        return false
    } else if info.IsDir() {
        return false
    }
    return true
}

// Hostname a prueba de lerdos de Dinosol
func gethostname(host string)(string){
	r:=strings.NewReplacer("BK","")
	trozo:=strings.Split(r.Replace(strings.ToUpper(host)),".")
	return trozo[0]  
}

// devuelve en Windows la cantidad de procesos ejecutandose con el nombre name
func procsrunning(name string) (num int) {
	out,_:=exec.Command("cmd","/c","tasklist").CombinedOutput()
	lines := strings.Split(string(out),"\n")
	for _,v:=range lines {
		if strings.Contains(v,name+".exe") { num++ }
	}
	return 
}

// funcion que cada hora actualiza el crontab del contenido del server interno
// /usr/bin/wget -q -O /dev/stdout -T 5 -t 3 http://localhost/megafonia2/crontab.php
func getcrontab(){
	for{
		url := server + "/megafonia2/crontab.php"
		fmt.Printf("[getcrontab(1)] - Actualizando CRONTAB\n")
		res := geturl(url, 3)
		if res != "" {
			lines := strings.Split(res,"<br>")
			i := 0
			fmt.Printf("[getcrontab(2)] - CRONTAB Actualizado OK\n")
			for _,line := range lines {
				//fmt.Println(line)
				if line == "" { continue }
				crontab[i]=line
				i++
			}
		}
		time.Sleep(60 * time.Minute) // 1 hora	
	}
}

// funcion que devuelve la playlist que corresponde Now()
// devuelve 0 (general),1 (suave) o -1 (stop) de la var. crontab al momento actual
func program() (playstatus string){
	playstatus = "-1"
	h,m,s := time.Now().Clock()
	seconds := h*3600 + m*60 + s
	for k,v := range crontab {
		seg := strings.Split(v,":")
		hh,_ := strconv.Atoi(seg[0])
		mm,_ := strconv.Atoi(seg[1])
		if (hh*3600+mm*60) <= seconds {
			playstatus = seg[2]
			continue
		}else{
			if playstatus == "-1" && k != 0 { playstatus = seg[2] }
			kk = k
			break
		}
	}
	return
}

// funcion que para la reproduccion actual por completo
func stop(){
	fmt.Printf("[stop(1)] - STOP Playing ...\n")
	playing="-1"
	exec.Command("cmd","/c","taskkill /F /IM overplay.exe").Run()
	exec.Command("cmd","/c","taskkill /F /IM RadioSupersol.exe").Run()
	time.Sleep(1 * time.Second)
}
