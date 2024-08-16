package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var Host, Input, User, Password, Port string
var Delay int
var Out, V bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&User, "u", "postgres", "")
	flag.StringVar(&Port, "P", "5432", "")
	flag.StringVar(&Password, "p", "postgres", "")
	flag.BoolVar(&V, "v", true, "")
	flag.BoolVar(&Out, "o", true, "")
	flag.Parse()
	start()
}

func start() {
	if Input != "" {
		Allip = list(Input)
		//fmt.Println(1)
		gopostgresql(Allip)
	} else if Host != "" {
		if strings.Contains(Host, "/") || strings.Contains(Host, "-") {
			ip := parseIP(Host)
			Allip = ip
			//fmt.Println(2)
			gopostgresql(Allip)
		} else {
			Allip = append(Allip, Host)
			//fmt.Println(3)
			gopostgresql(Allip)
		}
		//fmt.Println(Allip)
	}
}

func gopostgresql(ips []string) {
	for _, v := range ips {
		topostgresql(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func write(str string) {
	ok, err := os.OpenFile("postgresql.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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

func topostgresql(ip string) {
	if V {
		fmt.Println(ip)
	}
	connStr := "postgres://" + User + ":" + Password + "@" + ip + ":" + Port + "/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening connection: ", err)
	}
	defer db.Close()

	// 检查连接
	err = db.Ping()
	if err != nil {
		//log.Fatal("Error connecting to the database: ", err)
		return
	}
	if err != nil {
		//fmt.Println("Can't ping connection: ", err)
		return
	} else {
		fmt.Println(ip, Port, User, Password)
		if Out {
			write("[+] postgresql:" + ip + ":" + Port + ":" + User + " " + Password + "\n")
		}
	}
}
