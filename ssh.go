package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var Host, Input, User, Password, Port string
var Delay int
var Out bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&User, "u", "root", "")
	flag.StringVar(&Port, "P", "22", "")
	flag.StringVar(&Password, "p", "root", "")
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
		sshConn(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func sshConn(ip string) {
	addr := ip + ":" + Port
	// 创建 SSH 配置
	config := &ssh.ClientConfig{
		User: User,
		Auth: []ssh.AuthMethod{
			ssh.Password(Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 连接到 SSH 服务器
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("Failed to authenticate for", ip, ":", err)
		return
	} else {
		fmt.Println(addr, User, Password)
		if Out {
			//[+] SSH:127.0.0.1:22:root 123456
			write("[+] SSH:" + addr + ":" + User + " " + Password + "\n")
		}
	}

	defer conn.Close()
}

func write(str string) {
	ok, err := os.OpenFile("ssh.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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
