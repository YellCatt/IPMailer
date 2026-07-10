package main

import (
	"bytes"
	"fmt"
	"net"
	"net/smtp"
	"os/exec"
	"strings"
	"time"
)

func getLocalIPAddresses() string {
	var result strings.Builder

	ifaces, err := net.Interfaces()
	if err != nil {
		return "获取网络接口信息失败: " + err.Error()
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ipStr := ip.String()
			if strings.Contains(ipStr, ":") {
				continue
			}

			result.WriteString("- " + ipStr + "\n")
		}
	}

	if result.Len() == 0 {
		return getLocalIPAddressesFallback()
	}

	return result.String()
}

func getLocalIPAddressesFallback() string {
	commands := []string{
		"ip addr show",
		"ifconfig",
		"cat /proc/net/fib_trie",
	}

	for _, cmd := range commands {
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			continue
		}

		var c *exec.Cmd
		if len(parts) > 1 {
			c = exec.Command(parts[0], parts[1:]...)
		} else {
			c = exec.Command(parts[0])
		}

		var stderr bytes.Buffer
		c.Stderr = &stderr
		output, err := c.Output()
		if err != nil {
			continue
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "inet ") && !strings.Contains(line, "127.0.0.1") && !strings.Contains(line, "::1") {
				parts := strings.Fields(line)
				for _, part := range parts {
					if strings.Contains(part, ".") && !strings.Contains(part, ":") {
						ip := strings.Split(part, "/")[0]
						if ip != "127.0.0.1" {
							return "- " + ip + "\n"
						}
					}
				}
			}
		}
	}

	return "未找到有效的本地网络 IP 地址"
}

func main() {
	fromEmail := "768305875@qq.com"
	toEmail := "768305875@qq.com"
	authCode := "gpfruabgjebubdad"

	ipInfo := getLocalIPAddresses()
	subject := fmt.Sprintf("本地 IP 地址信息 %s", time.Now().Format("2006-01-02 15:04"))
	body := fmt.Sprintf("本地联网 IP 地址信息：\n\n%s\n\n发送时间：%s\n来自 Go 程序", ipInfo, time.Now().Format(time.RFC1123))

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		fromEmail,
		toEmail,
		subject,
		body,
	)

	fmt.Println("连接 SMTP...")
	auth := smtp.PlainAuth("", fromEmail, authCode, "smtp.qq.com")

	err := smtp.SendMail("smtp.qq.com:465", auth, fromEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		fmt.Printf("发送失败: %v\n", err)
		return
	}

	fmt.Println("✅ 已发送")
}