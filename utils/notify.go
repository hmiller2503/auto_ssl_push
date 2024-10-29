package utils

import (
	"autoUCDNcert/config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
)

func SendNotify(message string) {
	if config.AppConfig.NotifyType == "email" {
		sendEmail(message)
	} else if config.AppConfig.NotifyType == "wx_webhook" {
		sendWeChatWebhook(message)
	}
}

func sendEmail(content string) {
	smtpHost := config.AppConfig.EmailHost
	smtpPort := config.AppConfig.EmailPort

	from := config.AppConfig.EmailUsername
	password := config.AppConfig.EmailPwd
	to := []string{config.AppConfig.EmailReceiver}

	subject := "自动同步SSL至UCLOUD程序通知"
	body := fmt.Sprint("自动同步SSL至UCLOUD程序通知: \n", content)

	// 连接到 SMTP 服务器
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, &tls.Config{
		InsecureSkipVerify: true, // 测试时可以设置为true，生产环境请设为false
		ServerName:         smtpHost,
	})
	if err != nil {
		fmt.Println("连接到SMTP服务器失败:", err)
		return
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println("创建SMTP客户端失败:", err)
		return
	}

	// 使用PlainAuth进行认证
	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err = client.Auth(auth); err != nil {
		fmt.Println("认证失败:", err)
		return
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		fmt.Println("设置发件人失败:", err)
		return
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			fmt.Println("设置收件人失败:", err)
			return
		}
	}

	// 获取写入接口
	writer, err := client.Data()
	if err != nil {
		fmt.Println("获取写入接口失败:", err)
		return
	}

	// 邮件消息内容，确保包含From, To和Subject头信息
	message := []byte("From: " + from + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// 写入邮件内容
	_, err = writer.Write(message)
	if err != nil {
		fmt.Println("写入邮件内容失败:", err)
		return
	}

	// 关闭写入接口
	err = writer.Close()
	if err != nil {
		fmt.Println("关闭写入接口失败:", err)
		return
	}

	// 发送邮件完成
	client.Quit()
	fmt.Println("邮件发送成功!")
}

type WeChatMessage struct {
	MsgType string      `json:"msgtype"`
	Text    MessageText `json:"text"`
}

type MessageText struct {
	Content string `json:"content"`
}

func sendWeChatWebhook(message string) error {
	webhookURL := config.AppConfig.WxWebhookURL
	// 构建消息内容
	msg := WeChatMessage{
		MsgType: "text",
		Text:    MessageText{Content: fmt.Sprint("自动同步SSL至UCLOUD程序通知: \n", message)},
	}

	// 将消息编码为 JSON
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("编码消息失败: %v", err)
	}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(msgBytes))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	fmt.Println("消息发送成功!")
	return nil
}
