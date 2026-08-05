package contract

import "time"

type ListNetworkTrafficRequest struct {
	ListQuery
}

type ListNetworkTrafficResponse struct {
	Timestamp       time.Time `json:"timestamp"`
	SourceIP        string    `json:"sourceIp"`
	DestinationIP   string    `json:"destinationIp"`
	SourcePort      int       `json:"sourcePort"`
	DestinationPort int       `json:"destinationPort"`
	Protocol        string    `json:"protocol"`
	BytesIn         int64     `json:"bytesIn"`
	BytesOut        int64     `json:"bytesOut"`
	PacketCount     int       `json:"packetCount"`
	Description     string    `json:"description"`
}
