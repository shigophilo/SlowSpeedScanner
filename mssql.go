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

	_ "github.com/denisenkom/go-mssqldb"
)

var Host, Input, User, Password, Port string
var Delay int
var Out bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&User, "u", "sa", "")
	flag.StringVar(&Port, "P", "1433", "")
	flag.StringVar(&Password, "p", "sa", "")
	flag.BoolVar(&Out, "o", true, "")
	flag.Parse()
	start()
}

func start() {
	if Input != "" {
		Allip = list(Input)
		//fmt.Println(1)
		gossh(Allip)
	} else if Host != "" {
		if strings.Contains(Host, "/") {
			ip := parseIP(Host)
			Allip = ip
			//fmt.Println(2)
			gossh(Allip)
		} else {
			Allip = append(Allip, Host)
			//fmt.Println(3)
			gossh(Allip)
		}
		//fmt.Println(Allip)
	}
}

func gossh(ips []string) {
	for _, v := range ips {
		toMssql(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func write(str string) {
	ok, err := os.OpenFile("mssql.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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

func toMssql(ip string) {
	var dbinfo string
	// 打开数据库连接
	if Port == "1433" {
		dbinfo = "server=" + ip + ";user id=" + User + ";password=" + Password + ";database=master"
	} else {
		dbinfo = "server=" + ip + ":" + Port + ";user id=" + User + ";password=" + Password + ";database=master"
	}
	//fmt.Println(dbinfo)
	// "server=localhost;user id=sa;password=your_password;database=my_database"
	db, err := sql.Open("sqlserver", dbinfo)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		//fmt.Printf("Connection failed: %v\n", err)
		return
	} else {
		fmt.Println(ip, Port, User, Password)
		if Out {
			//[+] mssql 127.0.0.1:1433:sa pass@123
			write("[+] mssql:" + ip + ":" + Port + ":" + User + " " + Password + "\n")
		}
	}
}
