package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var Host, Input, Password, Port string
var Delay int
var Out, V bool
var Allip []string

func main() {
	flag.IntVar(&Delay, "d", 0, "")
	flag.StringVar(&Host, "h", "", "")
	flag.StringVar(&Input, "f", "", "")
	flag.StringVar(&Port, "P", "6379", "")
	flag.StringVar(&Password, "p", "", "")
	flag.BoolVar(&V, "v", true, "")
	flag.BoolVar(&Out, "o", true, "")
	flag.Parse()
	start()
}

func start() {
	if Input != "" {
		Allip = list(Input)
		//fmt.Println(1)
		goredis(Allip)
	} else if Host != "" {
		if strings.Contains(Host, "/") || strings.Contains(Host, "-") {
			ip := parseIP(Host)
			Allip = ip
			//fmt.Println(2)
			goredis(Allip)
		} else {
			Allip = append(Allip, Host)
			//fmt.Println(3)
			goredis(Allip)
		}
		//fmt.Println(Allip)
	}
}

func goredis(ips []string) {
	for _, v := range ips {
		toredis(v)
		if Delay != 0 {
			time.Sleep(time.Second * time.Duration(Delay))
		}
	}
}

func write(str string) {
	ok, err := os.OpenFile("redis.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
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

func toredis(ip string) {
	if V {
		fmt.Println(ip)
	}
	// 创建一个新的 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     ip + ":" + Port,
		Password: Password,
		DB:       0, // 默认DB 0
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		//log.Println(err)
		return
	} else {
		if Password == "" {
			Password = "unauthorized"
		}
		fmt.Println(ip, Port, Password)
		if Out {
			write("[+] Redis:" + ip + ":" + Port + ":" + " " + Password + "\n")
		}
	}
}
