package service

import (
	"BlackHole/internal/voidengine/contract"
	"BlackHole/internal/voidengine/model"
	"context"
	"testing"
	"time"
)

type trafficDataAccessStub struct {
	traffics []model.NetworkTraffic
}

func (s trafficDataAccessStub) List(context.Context, model.TrafficQuery) ([]model.NetworkTraffic, error) {
	return s.traffics, nil
}

func TestTrafficListReturnsContractTraffic(t *testing.T) {
	timestamp := time.Date(2026, 8, 5, 16, 0, 0, 0, time.UTC)
	service := NewTrafficService(trafficDataAccessStub{
		traffics: []model.NetworkTraffic{{
			ID:              1,
			Timestamp:       timestamp,
			SourceIP:        "10.0.0.1",
			DestinationIP:   "10.0.0.2",
			SourcePort:      1000,
			DestinationPort: 2000,
			Protocol:        "tcp",
			BytesIn:         123,
			BytesOut:        456,
			PacketCount:     7,
			Description:     "demo",
		}},
	})

	traffics, err := service.List(context.Background(), contract.ListNetworkTrafficRequest{})
	if err != nil {
		t.Fatalf("List error = %v", err)
	}
	if len(traffics) != 1 {
		t.Fatalf("len(traffics) = %d, want 1", len(traffics))
	}
	if traffics[0].Timestamp != timestamp {
		t.Fatalf("timestamp = %v, want %v", traffics[0].Timestamp, timestamp)
	}
	if traffics[0].SourceIP != "10.0.0.1" {
		t.Fatalf("source ip = %q, want 10.0.0.1", traffics[0].SourceIP)
	}
	if traffics[0].DestinationIP != "10.0.0.2" {
		t.Fatalf("destination ip = %q, want 10.0.0.2", traffics[0].DestinationIP)
	}
	if traffics[0].SourcePort != 1000 {
		t.Fatalf("source port = %d, want 1000", traffics[0].SourcePort)
	}
	if traffics[0].DestinationPort != 2000 {
		t.Fatalf("destination port = %d, want 2000", traffics[0].DestinationPort)
	}
	if traffics[0].Protocol != "tcp" {
		t.Fatalf("protocol = %q, want tcp", traffics[0].Protocol)
	}
	if traffics[0].BytesIn != 123 {
		t.Fatalf("bytes in = %d, want 123", traffics[0].BytesIn)
	}
	if traffics[0].BytesOut != 456 {
		t.Fatalf("bytes out = %d, want 456", traffics[0].BytesOut)
	}
	if traffics[0].PacketCount != 7 {
		t.Fatalf("packet count = %d, want 7", traffics[0].PacketCount)
	}
	if traffics[0].Description != "demo" {
		t.Fatalf("description = %q, want demo", traffics[0].Description)
	}
}
