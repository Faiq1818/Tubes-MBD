package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type SeismicData struct {
	DeviceID  uint64  `json:"device_id"`
	Timestamp int64   `json:"timestamp"`
	AccX      float32 `json:"acc_x"`
	AccY      float32 `json:"acc_y"`
	AccZ      float32 `json:"acc_z"`
	PGA       float32 `json:"pga"`
	STALTA    float32 `json:"sta_lta"`
	IsTrigger bool    `json:"is_trigger"`
}

var bufferPool = sync.Pool{
	New: func() any {
		fmt.Println("Pool making new buffer on heap")

		b := make([]byte, 1024)
		return &b
	},
}

func parseSeismicData(data []byte) (SeismicData, error) {
	if len(data) < 37 {
		return SeismicData{}, fmt.Errorf("packet too short")
	}

	return SeismicData{
		DeviceID:  binary.LittleEndian.Uint64(data[0:8]),
		Timestamp: int64(binary.LittleEndian.Uint64(data[8:16])),
		AccX:      math.Float32frombits(binary.LittleEndian.Uint32(data[16:20])),
		AccY:      math.Float32frombits(binary.LittleEndian.Uint32(data[20:24])),
		AccZ:      math.Float32frombits(binary.LittleEndian.Uint32(data[24:28])),
		PGA:       math.Float32frombits(binary.LittleEndian.Uint32(data[28:32])),
		STALTA:    math.Float32frombits(binary.LittleEndian.Uint32(data[32:36])),
		IsTrigger: data[36] == 1,
	}, nil
}

func response(conn *net.UDPConn, addr *net.UDPAddr, payload SeismicData) {
	respMsg := fmt.Sprintf("ACK: Device %d received at %d", payload.DeviceID, time.Now().Unix())

	fmt.Println("Goroutine executed")

	_, err := conn.WriteToUDP([]byte(respMsg), addr)
	if err != nil {
		fmt.Println("Error sending back data:", err)
	}
}

func main() {
	// kafka connection setup
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "user-events",
		Balancer:               &kafka.LeastBytes{},
		Async:                  false,
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()

	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening UDP:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP Server running on  port :9999...")

	for {
		bufPtr := bufferPool.Get().(*[]byte)
		buf := *bufPtr

		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading data:", err)
			continue
		}

		// go response(conn, remoteAddr, buf[:n])

		go func(addr *net.UDPAddr, rawData []byte, originalBuf *[]byte) {
			defer bufferPool.Put(originalBuf)

			seismic, err := parseSeismicData(rawData)
			if err != nil {
				fmt.Println("Parse error:", err)
				return
			}

			response(conn, addr, seismic)
			kafkaPayload, _ := json.Marshal(seismic)

			// kafka send
			err = writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(fmt.Sprintf("%d", seismic.DeviceID)),
					Value: kafkaPayload,
				},
			)
			if err != nil {
				fmt.Println("Error writing to Kafka:", err)
			} else {
				fmt.Println("Message successfully persisted to Kafka")
			}

		}(remoteAddr, append([]byte(nil), buf[:n]...), bufPtr)

	}

}
