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

	_ "github.com/go-sql-driver/mysql"
)

var Host, Input, User, Password, Port string
var Delay int
var Out, V bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&User, "u", "root", "")
	flag.StringVar(&Port, "P", "3306", "")
	flag.StringVar(&Password, "p", "root", "")
	flag.BoolVar(&V, "v", true, "")
	flag.BoolVar(&Out, "o", true, "")
	flag.Parse()
	start()
}

func start() {
	if Input != "" {
		Allip = list(Input)
		//fmt.Println(1)
		gomysql(Allip)
	} else if Host != "" {
		if strings.Contains(Host, "/") || strings.Contains(Host, "-") {
			ip := parseIP(Host)
			Allip = ip
			//fmt.Println(2)
			gomysql(Allip)
		} else {
			Allip = append(Allip, Host)
			//fmt.Println(3)
			gomysql(Allip)
		}
		//fmt.Println(Allip)
	}
}

func gomysql(ips []string) {
	for _, v := range ips {
		toMysql(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func write(str string) {
	ok, err := os.OpenFile("mysql.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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

// func Expand(s string) ([]string, error) {
// 	ip, ipnet, err := net.ParseCIDR(s)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nb, _ := ipnet.Mask.Size()
// 	if nb == 32 {
// 		return []string{ip.String()}, nil
// 	}
// 	nHosts := 1<<(32-nb) - 2
// 	ips := make([]string, nHosts)

// 	// 初始化 IP 地址
// 	ipBytes := ip.To4()
// 	if ipBytes == nil {
// 		return nil, fmt.Errorf("Invalid IPv4 address")
// 	}

// 	// 更新 IP 地址
// 	for n := 0; n < nHosts; n++ {
// 		ipBytes[3]++
// 		for i := 3; i >= 0; i-- {
// 			if ipBytes[i] > 255 {
// 				ipBytes[i] = 1
// 				if i > 0 {
// 					ipBytes[i-1]++
// 				}
// 			}
// 		}
// 		ips[n] = net.IP(ipBytes).String()
// 	}
// 	fmt.Println("ip:", ips)
// 	return ips, nil
// }

func toMysql(ip string) {
	if V {
		fmt.Println(ip)
	}
	db, err := sql.Open("mysql", User+":"+Password+"@tcp("+ip+":"+Port+")/"+"mysql")
	if err != nil {
		fmt.Println(ip, err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err == nil {
		fmt.Println(ip, Port, User, Password)
		if Out {
			write("[+] mysql:" + ip + ":" + Port + ":" + User + " " + Password + "\n")
		}
	}
}
