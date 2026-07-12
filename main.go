package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"time"
)

func getLocalIPAddresses() string {
	var result string

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
			if len(ipStr) > 0 && ipStr[0] != ':' {
				result += "- " + ipStr + "\n"
			}
		}
	}

	if result == "" {
		return "未找到有效的本地网络 IP 地址"
	}

	return result
}

func sendMailTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, skipVerify bool) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerify,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if err = c.Auth(auth); err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "未知设备"
	}
	return hostname
}

func main() {
	cfg := LoadConfig()

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	hostname := getHostname()
	ipInfo := getLocalIPAddresses()
	subject := fmt.Sprintf("[%s] 本地 IP 地址信息 %s", hostname, now.Format("2006-01-02 15:04"))
	body := fmt.Sprintf("设备名称：%s\n\n本地联网 IP 地址信息：\n\n%s\n\n发送时间：%s\n来自 Go 程序", hostname, ipInfo, now.Format(time.RFC1123))

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.FromEmail,
		cfg.ToEmail,
		subject,
		body,
	)

	fmt.Println("连接 SMTP...")
	auth := smtp.PlainAuth("", cfg.FromEmail, cfg.AuthCode, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

	if cfg.TLSSkipVerify {
		fmt.Println("⚠️ 警告：TLS 证书校验已跳过")
	}

	err := sendMailTLS(addr, auth, cfg.FromEmail, []string{cfg.ToEmail}, []byte(msg), cfg.TLSSkipVerify)
	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		return
	}

	fmt.Println("✅ 已发送")
}
