package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

var Host, Input, User, Password, Port string
var Delay int
var Out, V bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&User, "u", "system", "")
	flag.StringVar(&Port, "P", "1521", "")
	flag.StringVar(&Password, "p", "system", "")
	flag.BoolVar(&V, "v", true, "")
	flag.BoolVar(&Out, "o", true, "")
	flag.Parse()
	start()
}

func start() {
	if Input != "" {
		Allip = list(Input)
		//fmt.Println(1)
		gooracle(Allip)
	} else if Host != "" {
		if strings.Contains(Host, "/") || strings.Contains(Host, "-") {
			ip := parseIP(Host)
			Allip = ip
			//fmt.Println(2)
			gooracle(Allip)
		} else {
			Allip = append(Allip, Host)
			//fmt.Println(3)
			gooracle(Allip)
		}
		//fmt.Println(Allip)
	}
}

func gooracle(ips []string) {
	for _, v := range ips {
		tooracle(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func write(str string) {
	ok, err := os.OpenFile("oracle.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer ok.Close()
	_, err = ok.Write([]byte(str))
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func list(urlfile string) []string {
	var url_list []string
	url_file, err := os.Open(urlfile)
	if err != nil {
		fmt.Println("Can't open urlfile:", err)
		return nil
	}
	defer url_file.Close()

	reader_Url := bufio.NewReader(url_file)
	for {
		url, err := reader_Url.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading from file:", err)
			break
		}
		url = strings.TrimSpace(url) // 移除行尾的换行符
		url_list = append(url_list, url)
	}
	return url_list
}

func tooracle(ip string) {
	if V {
		fmt.Println(ip)
	}
	//oracle://sys:sys1@218.154.228.65:1521/ORCL
	server := "oracle://" + User + ":" + Password + "@" + ip + ":" + Port + "/ORCL"
	conn, err := sql.Open("oracle", server)
	if err != nil {
		//fmt.Println("Can't open the driver: ", err)
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			//fmt.Println("Can't close connection: ", err)
			return
		}
	}()

	err = conn.Ping()
	if err != nil {
		//fmt.Println("Can't ping connection: ", err)
		return
	} else {
		fmt.Println(ip, Port, User, Password)
		if Out {
			write("[+] oracle:" + ip + ":" + Port + ":" + User + " " + Password + "\n")
		}
	}
}
