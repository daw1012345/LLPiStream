package main

import (
	"github.com/pion/webrtc/v3"
	"log"
	"net"
)

// Wanted to use channels, but they're too slow and cause packet loss??

var track *webrtc.TrackLocalStaticRTP

func InitWebRTCTrack(mime string, sid string) {
	switch mime {
	case "video/VP8":
		break
	case "video/H264":
		break
	default:
		log.Panicln("Unsupported codec selected!")
	}

	ltrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: mime}, "video", sid)

	if err != nil {
		log.Panicln(err)
	}

	track = ltrack
}

func RunUDPRTPServer(ip string, port int) {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	if err != nil {
		log.Panicln(err)
	}

	defer ln.Close()

	log.Printf("Listening for RTP packets on udp://%s:%d", ip, port)

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
