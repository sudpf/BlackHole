package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"time"
)

// LogMessage 定义日志消息的结构体
type LogMessage struct {
	Timestamp       string `json:"timestamp"`        // 时间戳，使用 RFC3339 格式
	SourceIP        string `json:"source_ip"`        // 源 IP 地址
	DestinationIP   string `json:"destination_ip"`   // 目标 IP 地址
	SourcePort      int64  `json:"source_port"`      // 源端口
	DestinationPort int64  `json:"destination_port"` // 目标端口
	Protocol        string `json:"protocol"`         // 协议类型（如 TCP、UDP 等）
	BytesIn         int64  `json:"bytes_in"`         // 输入字节数
	BytesOut        int64  `json:"bytes_out"`        // 输出字节数
	PacketCount     int64  `json:"packet_count"`     // 数据包数量
	Description     string `json:"description"`      // 描述信息
}

// sendLog 通过 Unix Datagram Socket 发送日志消息
func sendLog(socketPath string, logMessage LogMessage) error {
	// 设置 Unix Socket 地址
	addr := net.UnixAddr{Name: socketPath, Net: "unixgram"}

	// 使用 Unixgram 协议连接 Unix Socket
	conn, err := net.DialUnix("unixgram", nil, &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 将日志消息编码为 JSON
	messageBytes, err := json.Marshal(logMessage)
	if err != nil {
		return err
	}

	log.Println(string(messageBytes))
	// 发送 JSON 格式的日志消息
	_, err = conn.Write(messageBytes)
	return err
}

func main() {
	// Unix Socket 地址（与 Syslog 服务配置一致）
	socketPath := "/tmp/suricata_unix_sock"

	// 检查 Socket 文件是否存在
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		log.Fatalf("Socket file does not exist: %v", socketPath)
	}

	// 创建一个示例日志消息
	logMessage := LogMessage{
		Timestamp:       time.Now().Format(time.RFC3339), // 使用高精度的时间格式
		SourceIP:        "192.168.1.100",
		DestinationIP:   "192.168.1.1",
		SourcePort:      12345,
		DestinationPort: 80,
		Protocol:        "TCP",
		BytesIn:         1024,
		BytesOut:        2048,
		PacketCount:     30,
		Description:     "Example log message for network traffic",
	}

	// 发送日志消息
	if err := sendLog(socketPath, logMessage); err != nil {
		log.Fatalf("Failed to send log: %v", err)
	}

}
