package main

import ("fmt";"os/exec";"bufio";"strings";"net/http";"io/ioutil";"os";"time";"strconv";"math/rand")

var settings map[string]string = make(map[string]string)
var publi0 map[int]string = make(map[int]string)
var publi1 map[int]string = make(map[int]string)
var spacing0,spacing1,license0,license1,general0,general1,special0,special1 string
var all = "dir\\*BEAT* dir\\*BLUES* dir\\*COUNTRY* dir\\*DANCE* dir\\*FLAMENC* dir\\*LATIN* dir\\*OLDIES* dir\\*POPESP* dir\\*POPINT* dir\\*RAPINT* dir\\*REGGAE* dir\\*ROCKESP* dir\\*ROCKINT* dir\\*SALSA* dir\\*SAMBA* "
var mid = "dir\\*BLUES* dir\\*CHILLOUT* dir\\*FLAMENC* dir\\*LATIN* dir\\*POPESP* dir\\*POPINT* dir\\*RAPINT* "
var mediadir = "C:\\RSPS" //  "musica" ; "/home/pymedia/audio"
var server = "http://10.26.4.121" //
var settingsFile = "C:\\RSPS\\settings.reg" // "/home/pymedia/local/bin/settings.reg"
// usage: updatefiles -d (no borra dicheros)
var delete, help bool

func main(){
	go func(){
		time.Sleep(58*time.Minute)
		cmdp := exec.Command("cmd","/c","taskkill /F /IM wget.exe") // taskkill /F /IM wget.exe winamp.exe
		cmdp.Run()
		fmt.Println("[main(0)] - Ejecucion muy larga - interrumpida a los 59 minutos")
		os.Exit(1)		
	}()
	if len(os.Args) > 1 {
		if os.Args[1] == "-d" { delete = true }
		if os.Args[1] == "-h" { help = true }
	}
	if help {
		fmt.Printf("Usage: %s [options]\nOptions:\n",os.Args[0])
		fmt.Printf("\t-h Show this help\n\t-d Don't delete files\n")
		os.Exit(0)
	}
	ping()
	fileinfo,err := os.Stat(settingsFile)
	if err != nil || fileinfo.Size() == 0 {
		fmt.Println("No existe el fichero: "+settingsFile)
		settings["last"]=""
	}else{
		loadSettings(settingsFile)
	}
	y,m,d := time.Now().Date()
	if settings["last"] == fmt.Sprintf("%d-%s-%d",d,m.String(),y) {
		fmt.Println("[main(1)] - Listas ya bajadas. No update")
		i,j:=0,0
		for k,v := range settings{
			if strings.Contains(k,"pub0") {
				publi0[i]=v; i++
			}else if strings.Contains(k,"pub1") {
				publi1[j]=v; j++
			}
		}
	}else{
		fmt.Println("[main(2)] - Listas aun sin bajar. Updating !!!")
		getplaylists()
		doplay ("0" ,mediadir+"\\playlist1.m3u")
		doplay ("1" ,mediadir+"\\playlist2.m3u")
		doplay ("0" ,mediadir+"\\playlist3.m3u")
		doplay ("1" ,mediadir+"\\playlist4.m3u")
		doplay ("0" ,mediadir+"\\playlist5.m3u")
		for k,v := range publi0{
			settings["pub0"+strconv.Itoa(k)]=v
		}
		for k,v := range publi1{
			settings["pub1"+strconv.Itoa(k)]=v
		}
		saveSettings(settingsFile)
	}
	getfiles() // este proceso puede durar mas de 1 hora a veces (interrumpible)
}

func ping(){
	var url string
	hostname,_:=os.Hostname()
	url = server + ":8080/megafonia2/ping.php?tienda="+gethostname(hostname)  // server + ":8080/megafonia2/playlist.php?tienda="+hostname
	fmt.Printf("[ping(1)] - I'm alive & kicking: \n")
	geturl(url, 3)
}

// funcion que carga los settings de un fichero *.reg
func loadSettings (filename string){
	fr,err := os.Open(filename)
	defer fr.Close()
	if err == nil{
		reader := bufio.NewReader(fr)
		for{
			linea,rerr := reader.ReadString('\n')
			if rerr != nil { break; }
			linea = strings.TrimRight(linea,"\n")
			item := strings.Split(linea,"=")
			if len(item) == 2 {
				settings[item[0]]=item[1]
			}
		}
	}
}

// funcion que guarda los settings en un fichero *.reg
func saveSettings (filename string) {
	fw,err := os.Create(filename)
	if err == nil {
		writer := bufio.NewWriter(fw)
		for k,v := range settings {
			writer.WriteString(k+"="+v+"\n")
		}
		writer.Flush()
	}
	fw.Close()
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

// funcion que reproduce la playlist que se le pasa, o espera hasta que esta este preparada
// General (play=0) ; Suave (play=1) ; Nada (play=x)
func doplay (play ,filename string){
	var pub map[int]string = make(map[int]string)
	var audiogap, pubgap int
	var music string
	pub = map[int]string{}
	// establecimiento de variables necesarias
	if play == "0" { // audiogap, pubgap, music, publi[]
		gap:=strings.Split(spacing0,"."); audiogap,_ = strconv.Atoi(gap[0]); pubgap,_ = strconv.Atoi(gap[1]);
		r:=strings.NewReplacer("dir", mediadir)
		music = r.Replace(general0 + special0) // filtro de busqueda ls -l
		for k,v := range publi0 {
			pub[k]=v
		}
	}else if play == "1" { // audiogap, pubgap, music, publi[]
		gap:=strings.Split(spacing1,"."); audiogap,_ = strconv.Atoi(gap[0]); pubgap,_ = strconv.Atoi(gap[1]);
		r:=strings.NewReplacer("dir", mediadir)
		music = r.Replace(general1 + special1) // filtro de busqueda ls -l
		for k,v := range publi1 {
			pub[k]=v
		}
	}else {
		return
	}
	// ya podemos comenzar la creacion de la playlist conforme a las reglas audiogap, pubgap, music, publi[]
	audio := make(map[int]string);
	audio = map[int]string{}
	
	{
		fmt.Printf("[doplay(2)] - Recargamos la musica accesible en el disco duro para la playlist: %s\n",play)
		fmt.Println("dir /B "+music)
		cmd := exec.Command("cmd","/c","dir /B "+music) // dir /B 20*.mp3
		// comienza la ejecucion del pipe
		i := 0;
		stdoutRead,_ := cmd.StdoutPipe()
		reader := bufio.NewReader(stdoutRead)
		cmd.Start()
		for{
			line,err := reader.ReadString('\n')
			if err != nil {
				break;
			}
			//fmt.Printf("%s",line)
			audio[i]=strings.TrimRight(strings.TrimRight(line,"\n"),"\r")
			i++
		}
		cmd.Wait()
		fmt.Printf("[doplay(3)] - Shuffle de MUSIC: %d canciones\n",i)
		rand.Seed(time.Now().UnixNano())
		shuffle := rand.Perm(len(audio)) // shuffle[0-n)
		i=0
		// elegimos la publicidad a reproducir
		fmt.Printf("[doplay(4)] - Preparamos la publicidad de HOY: %d anuncios\n",len(pub))
		// creamos la playlist mezclando audio + pub con el gap correspondiente
		a,p := 0,0
		fw,_ := os.Create(filename)
		writer := bufio.NewWriter(fw)
		writer.WriteString("#EXTM3U\r\n")
		for _,v:= range shuffle {
			a++
			writer.WriteString(audio[v]+"\r\n")
			if a == audiogap {
				for p<pubgap {
					if len(pub)==0 { break }
					i++;p++
					writer.WriteString(pub[i%len(pub)]+"\r\n")
				}
				a,p = 0,0
			}
		}
		writer.Flush()
		fw.Close()
	}
}

// baja una sola vez al dia reintando cada hora hasta q lo consigue la playlist, si lo consigue la parsea
// y graba: spacing0/1,license0/1,general0/1,special0/1,publi0/1[]
func getplaylists(){
	var url string
	hostname,_:=os.Hostname()
	url = server + ":8080/megafonia2/playlist.php?tienda="+gethostname(hostname)  // server + ":8080/megafonia2/playlist.php?tienda="+hostname
	fmt.Printf("[getplaylists(1)] - Actualizando PLAYLISTS: \n")
	res := geturl(url, 3)
	if res == "" {
		fmt.Printf("[getplaylists(2)] - DEFAULT PLAYLISTS: \n")
		// ponemos por defecto
		publi0 = map[int]string{}
		publi1 = map[int]string{}
		spacing0="2.1";spacing1="2.1";
		license0="0";license1="0";
		general0=all;general1=mid;
		special0="";special1="";
	}else{
		y,m,d := time.Now().Date()
		settings["last"] = fmt.Sprintf("%d-%s-%d",d,m.String(),y)
		fmt.Printf("[getplaylists(3)] - Parseando PLAYLISTS actualizada: \n")
		parseplaylist(res)
	}
}

// parsea la playlist recien bajada y graba: spacing0/1,license0/1,general0/1,special0/1,publi0/1[]
func parseplaylist(res string){
	var block int
	i,j := 0,0
	publi0 = map[int]string{}
	publi1 = map[int]string{}
	lines := strings.Split(res,"<br>")
	for _,line := range lines {
		if line == "" { continue }
		switch line {
			case "#musica0":
				block=0
			case "#musica1":
				block=1
			case "#publi0":
				block=2
			case "#publi1":
				block=3
			default:
				switch block {
					case 0:
						if strings.Contains(line,"spacing") {
							spacing0 = strings.Trim(line,"spacing=")
						}else if strings.Contains(line,"license") {
							license0 = strings.Trim(line,"license=")
						}else if strings.Contains(line,"general") {
							general0 = strings.Trim(line,"general=")
						}else if strings.Contains(line,"special") {
							special0 = strings.Trim(line,"special=")
						}
					case 1:
						if strings.Contains(line,"spacing") {
							spacing1 = strings.Trim(line,"spacing=")
						}else if strings.Contains(line,"license") {
							license1 = strings.Trim(line,"license=")
						}else if strings.Contains(line,"general") {
							general1 = strings.Trim(line,"general=")
						}else if strings.Contains(line,"special") {
							special1 = strings.Trim(line,"special=")
						}
					case 2:
						publi0[i] = line
						i++
					case 3:
						publi1[j] = line
						j++
				}
		}
	}
	// vamos a extraer los afijos de las listas general y special
	url := server + "/megafonia2/lists/"
	if general0 == "ALL" {
		general0 = all
	}else if general0 == "MID" {
		general0 = mid
	}else {
		res = geturl(url+general0, 3)
		lines := strings.Split(res," ")
		general0 = ""
		for _,line := range lines {
			if line != "" && line != " " && line != "\n" { general0 += "dir\\*"+line+"* " }
		}
	}
	
	if general1 == "ALL" {
		general1 = all
	}else if general1 == "MID" {
		general1 = mid
	}else {
		res = geturl(url+general1, 3)
		lines := strings.Split(res," ")
		general1 = ""
		for _,line := range lines {
			if line != "" && line != " " && line != "\n" { general1 += "dir\\*"+line+"* " }
		}
	}

	if special0 != ""{res = geturl(url+special0, 3)}else {res=""}
	lines = strings.Split(res," ")
	if res == "" && special0 != "" {
		special0 = "dir\\*"+special0+"* "
	}else{
		special0 = ""
		for _,line := range lines {
			if line != "" && line != " " && line != "\n" { special0 += "dir\\*"+line+"* " }
		}
	}

	if special0 != ""{res = geturl(url+special0, 3)}else {res=""}
	lines = strings.Split(res," ")
	if res == "" && special1 != "" {
		special1 = "dir\\*"+special1+"* "
	}else{
		special1 = ""
		for _,line := range lines {
			if line != "" && line != " " && line != "\n" { special1 += "dir\\*"+line+"* " }
		}
	}
}

// funcion que actualiza por clonado el contenido del servidor interno: musica,sgae,pub
// borrando lo viejo y bajando lo nuevo (revisa cada 1 hora, actualiza primero toda la publi, pero no baja mas de 4 canciones a la hora)
func getfiles(){
	// a = ficheros que tenemos aqui ; b = ficheros que hay en el server interno
	var a map[string]bool
	var b map[string]bool
	{
		fmt.Printf("[getfiles(1)] - Actualizando el contenido del disco duro\n")
		// 5- ahora tratamos la publicidad
		a = map[string]bool{}  // vaciamos 'a' de nuevo, ahora llevara la pub del disco duro
		b = map[string]bool{}  // vaciamos 'b' de nuevo, ahora llevara la pub programada para HOY (playlists)
		fmt.Println("dir /B "+mediadir+"\\20*.mp3")
		cmd := exec.Command("cmd","/c","dir /B C:\\RSPS\\20*.mp3") // dir /B 20*.mp3
		// comienza la ejecucion del pipe
		stdoutRead,_ := cmd.StdoutPipe()
		reader := bufio.NewReader(stdoutRead)
		cmd.Start()
		for{
			line,err := reader.ReadString('\n')
			if err != nil {
				break;
			}
			name:=strings.TrimRight(strings.TrimRight(line,"\n"),"\r")
			a[name]=true
		}
		cmd.Wait()
		for _,v := range publi0 {
			b[v]=true
		}
		for _,v := range publi1 {
			b[v]=true
		}
		if len(b)!=0 {
			// borramos lo que sobra
			if !delete {
				e := diffsets(a,b)
				for k,_ := range e {
					fmt.Printf("[getfiles(5)] - del /Q %s\\%s\n",mediadir,k)
					cmdn := exec.Command("cmd","/c","del /Q C:\\RSPS\\"+k)
					cmdn.Run()
				}
			}
			// bajamos lo que necesitamos
			f:= diffsets(b,a)
			for k,_ := range f {
				// 	wget -T 5 -t 3 -P downloads/ http://xxxxx
				fmt.Printf("[getfiles(6)] - c:\\wget.exe -T 5 -t 3 -P c:\\RSPS %s/megafonia2/media/pub/%s\n",server,k)
				cmdn := exec.Command("cmd","/c","c:\\wget.exe -T 5 -t 3 -P c:\\RSPS "+server + "/megafonia2/media/pub/"+k)
				cmdn.Run()
			}
		}
		// 1-bajamos la list de fichero del server interno
		url := server + "/megafonia2/files2.php" // se baja sin SGAE
		res := geturl(url, 3)
		if res != "" { 
			a = map[string]bool{}  // vaciamos 'a' de nuevo
			b = map[string]bool{}  // vaciamos 'b' de nuevo
			// 2-vemos los ficheros que tenemos nosotros
			fmt.Println("dir /B "+mediadir+"\\20*.mp3")
			cmd = exec.Command("cmd","/c","dir /B C:\\RSPS\\*.mp3") // dir /B *.mp3
			// comienza la ejecucion del pipe
			stdoutRead,_ = cmd.StdoutPipe()
			reader = bufio.NewReader(stdoutRead)
			cmd.Start()
			for{
				line,err := reader.ReadString('\n')
				if err != nil {
					break;
				}
				name:=strings.TrimRight(strings.TrimRight(line,"\n"),"\r")
				a[name]=true
			}
			cmd.Wait()
			cmd := exec.Command("cmd","/c","dir /B C:\\RSPS\\20*.mp3") // dir /B 20*.mp3
			// comienza la ejecucion del pipe
			stdoutRead,_ := cmd.StdoutPipe()
			reader := bufio.NewReader(stdoutRead)
			cmd.Start()
			for{
				line,err := reader.ReadString('\n')
				if err != nil {
					break;
				}
				name:=strings.TrimRight(strings.TrimRight(line,"\n"),"\r")
				b[name]=true
			}
			cmd.Wait()
			c := diffsets(a,b)
			b = map[string]bool{}  // vaciamos 'b' de nuevo
			b=parsefiles(res)
			// 3- vemos los que hay que borrar (4 cada hora max)
			max := 0
			if !delete {
				cc := diffsets(c,b)
				for k,_ := range cc {
					fmt.Printf("[getfiles(2)] - del /Q %s\\%s\n",mediadir,k) // del /Q nombre.mp3
					cmdn := exec.Command("cmd","/c","del /Q C:\\RSPS\\"+k) // dir /B 20*.mp3
					cmdn.Run()
					max++
					if max > 3 { break }
				}
			}
			// 4- vemos los que hay que bajar (4 cada hora max)
			max = 0
			d:= diffsets(b,c)
			for k,_ := range d {
				// 	wget -T 5 -t 3 -P downloads/ http://xxxxx
				fmt.Printf("[getfiles(6)] - c:\\wget.exe -T 5 -t 3 -P C:\\RSPS %s/megafonia2/media/musica/%s\n",server,k)
				cmdn := exec.Command("cmd","/c","c:\\wget.exe -T 5 -t 3 -P C:\\RSPS "+server + "/megafonia2/media/musica/"+k)
				cmdn.Run()
				max++
				if max > 3 { break }
			}
		}
	}
}

// Operacion Teoria de Conjuntos : C = A - B
func diffsets(a map[string]bool, b map[string]bool)(c map[string]bool){
	c  = map[string]bool{}
	for k,_ := range a {
		if !b[k] {
			c[k]=true
		}
	}
	return
}

// funcion que pasea la respuesta del server interno sobre el contenido de ficheros en su disco
func parsefiles(res string)(c map[string]bool){
	c  = map[string]bool{}
	var block int
	lines := strings.Split(res,"<br>")
	for _,line := range lines {
		//fmt.Println(line)
		if line == "" { continue }
		switch line {
			case "#sgae":
				block=0
			case "#musica":
				block=1
			case "#pub":
				block=2
			default:
				switch block {
					case 0: // sgae
						c[line]=true
					case 1: // musica
						c[line]=true
					case 2: // pub
						c[line]=true
				}
		}
	}
	return
}

// Hostname a prueba de lerdos de Dinosol
func gethostname(host string)(string){
	r:=strings.NewReplacer("BK","")
	trozo:=strings.Split(r.Replace(strings.ToUpper(host)),".")
	return trozo[0]  
}
