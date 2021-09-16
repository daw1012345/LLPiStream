package main

import (
	"github.com/pion/webrtc/v3"
	"log"
	"net"
)

// Wanted to use channels, but they're too slow and cause packet loss??

var track *webrtc.TrackLocalStaticRTP

func Init() {
	ltrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")

	if err != nil {
		log.Panicln(err)
	}

	track = ltrack
}

func StartUDPRTPServer() {
	ln, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 5004,
	})

	if err != nil {
		log.Panicln(err)
	}
	defer ln.Close()

	log.Println("Listening on udp://localhost:8091")
	buf := make([]byte, 1600)

	for {
		n, _, err := ln.ReadFrom(buf)

		if err != nil {
			log.Println(err)
			continue
		}

		track.Write(buf[:n])
	}
}