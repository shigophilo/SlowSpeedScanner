package main

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// func main() {
// 	fmt.Println(parseIP("172.16.1.250/12"))
// }

func parseIP(host string) []string {
	var ips []string

	if strings.Contains(host, "-") {
		return parseIP1(host)
	} else if strings.Contains(host, "/8") {
		return parseIP8(host)
	} else if strings.Contains(host, "/12") {
		return parseIP12(host)
	} else if strings.Contains(host, "/16") {
		return parseIP16(host)
	} else if strings.Contains(host, "/24") {
		return parseIP24(host)
	} else {
		return append(ips, host)
	}
}
func parseIP8(host string) []string {
	var ips []string
	ip := getIP(host)
	tmpip := strings.Split(ip, ".")
	for e := 0; e <= 255; e++ {
		for i := 0; i <= 255; i++ {
			for n := 1; n < 255; n++ {
				tip := tmpip[0] + "." + strconv.Itoa(e) + "." + strconv.Itoa(i) + "." + strconv.Itoa(n)
				ips = append(ips, tip)
			}
		}
	}
	return ips
}
func parseIP12(host string) []string {
	var ips []string
	ip := getIP(host)
	tmpip := strings.Split(ip, ".")
	for e := 16; e <= 32; e++ {
		for i := 0; i <= 255; i++ {
			for n := 1; n < 255; n++ {
				tip := tmpip[0] + "." + strconv.Itoa(e) + "." + strconv.Itoa(i) + "." + strconv.Itoa(n)
				ips = append(ips, tip)
			}
		}
	}
	return ips
}
func parseIP16(host string) []string {
	var ips []string
	ip := getIP(host)
	tmpip := strings.Split(ip, ".")
	for i := 0; i <= 255; i++ {
		for n := 1; n < 255; n++ {
			tip := tmpip[0] + "." + tmpip[1] + "." + strconv.Itoa(i) + "." + strconv.Itoa(n)
			ips = append(ips, tip)
		}
	}
	return ips
}
func parseIP24(host string) []string {
	var ips []string
	ip := getIP(host)
	tmpip := strings.Split(ip, ".")
	for i := 1; i < 255; i++ {
		tip := tmpip[0] + "." + tmpip[1] + "." + tmpip[2] + "." + strconv.Itoa(i)
		ips = append(ips, tip)
	}
	Shuffle(ips)
	return ips
}

// Shuffle 利用Fisher-Yates算法对切片中的元素进行随机排序。
func Shuffle(slice []string) []string {
	rand.Seed(time.Now().UnixNano()) // 使用当前时间作为随机种子

	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)                   // 生成一个0到i之间的随机索引
		slice[i], slice[j] = slice[j], slice[i] // 交换元素
	}
	return slice
}

func getIP(host string) string {
	ip := strings.Split(host, "/")
	return ip[0]
}

// 解析ip段:
//
//	192.168.111.1-255
//	192.168.111.1-192.168.112.255
func parseIP1(ip string) []string {
	IPRange := strings.Split(ip, "-")
	testIP := net.ParseIP(IPRange[0])
	var AllIP []string
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			AllIP = append(AllIP, ip)
		}
	}
	return AllIP
}
