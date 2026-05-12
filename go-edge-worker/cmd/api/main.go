package main

import (
	"context"
	"encoding/binary"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type SeismicData struct {
	DeviceID  string  `json:"device_id"`
	Timestamp int64   `json:"timestamp"`
	AccX      float32 `json:"acc_x"`
	AccY      float32 `json:"acc_y"`
	AccZ      float32 `json:"acc_z"`
	PGA       float32 `json:"pga"`
	STALTA    float32 `json:"sta_lta"`
	IsTrigger bool    `json:"is_trigger"`
}

const PacketSize = 45

var bufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 1024)
		return &b
	},
}

func parseSeismicData(data []byte) (SeismicData, error) {
	if len(data) < PacketSize {
		return SeismicData{}, fmt.Errorf(
			"packet too short: got=%d expected=%d",
			len(data),
			PacketSize,
		)
	}

	// UUID = first 16 bytes
	u, err := uuid.FromBytes(data[0:16])
	if err != nil {
		return SeismicData{}, err
	}

	return SeismicData{
		DeviceID:  u.String(),
		Timestamp: int64(binary.LittleEndian.Uint64(data[16:24])),

		AccX: math.Float32frombits(
			binary.LittleEndian.Uint32(data[24:28]),
		),

		AccY: math.Float32frombits(
			binary.LittleEndian.Uint32(data[28:32]),
		),

		AccZ: math.Float32frombits(
			binary.LittleEndian.Uint32(data[32:36]),
		),

		PGA: math.Float32frombits(
			binary.LittleEndian.Uint32(data[36:40]),
		),

		STALTA: math.Float32frombits(
			binary.LittleEndian.Uint32(data[40:44]),
		),

		IsTrigger: data[44] == 1,
	}, nil
}

func response(conn *net.UDPConn, addr *net.UDPAddr, payload SeismicData) {
	respMsg := fmt.Sprintf(
		"ACK: Device %s received at %d",
		payload.DeviceID,
		time.Now().Unix(),
	)

	_, err := conn.WriteToUDP([]byte(respMsg), addr)
	if err != nil {
		fmt.Println("response error:", err)
	}
}

func main() {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "user-events",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,

		// lebih cocok untuk high throughput
		Async: true,
	}
	defer writer.Close()

	addr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		fmt.Println("resolve error:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("listen error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP Server listening on :9999")

	for {
		bufPtr := bufferPool.Get().(*[]byte)
		buf := *bufPtr

		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("read error:", err)

			bufferPool.Put(bufPtr)
			continue
		}

		// copy exact packet size
		packet := make([]byte, n)
		copy(packet, buf[:n])

		// buffer langsung dikembalikan
		bufferPool.Put(bufPtr)

		go func(addr *net.UDPAddr, rawData []byte) {
			seismic, err := parseSeismicData(rawData)
			if err != nil {
				fmt.Println("parse error:", err)
				return
			}

			response(conn, addr, seismic)

			kafkaPayload, err := json.Marshal(seismic)
			if err != nil {
				fmt.Println("json error:", err)
				return
			}

			err = writer.WriteMessages(
				context.Background(),
				kafka.Message{
					Key:   []byte(seismic.DeviceID),
					Value: kafkaPayload,
				},
			)

			if err != nil {
				fmt.Println("kafka error:", err)
				return
			}

			fmt.Printf(
				"[KAFKA] device=%s pga=%.2f\n",
				seismic.DeviceID,
				seismic.PGA,
			)

		}(remoteAddr, packet)
	}
}
