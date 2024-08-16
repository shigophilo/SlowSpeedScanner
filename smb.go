package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/stacktitan/smb/smb"
)

var Host, Input, User, Password, DomainName string
var Delay int
var Out, V bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&DomainName, "domain", " ", "")
	flag.StringVar(&User, "u", "Administrator", "")
	flag.StringVar(&Password, "p", "123456", "")
	flag.BoolVar(&V, "v", true, "")
	flag.BoolVar(&Out, "o", true, "")
	flag.Parse()
	start()
	WG.Wait()
}

func start() {
	if Input != "" {
		Allip = list(Input)
		//fmt.Println(1)
		gosmb(Allip)
	} else if Host != "" {
		if strings.Contains(Host, "/") || strings.Contains(Host, "-") {
			ip := parseIP(Host)
			Allip = ip
			//fmt.Println(2)
			gosmb(Allip)
		} else {
			Allip = append(Allip, Host)
			//fmt.Println(3)
			gosmb(Allip)
		}
		//fmt.Println(Allip)
	}
}

func gosmb(ips []string) {
	for _, v := range ips {
		go tosmb(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func write(str string) {
	ok, err := os.OpenFile("smb.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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

func tosmb(ip string) {
	if V {
		fmt.Println(ip)
	}

	options := smb.Options{
		Host:        ip,
		Port:        445,
		User:        User,
		Domain:      DomainName,
		Workstation: "",
		Password:    Password,
	}
	debug := false
	session, err := smb.NewSession(options, debug)
	if err == nil {
		fmt.Println(ip, 445, User, Password)
		if Out {
			write("[+] SMB:" + ip + ":445" + ":" + User + " " + Password + "\n")
		}
	}
	defer session.Close()

}
