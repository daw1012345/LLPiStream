package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

func main() {
	haddr := flag.String("haddr", "127.0.0.1:8080", "The host/port to run the API server on")
	rtphost := flag.String("rtphost", "127.0.0.1", "The port to listen on for RTP packets")
	rtpport := flag.Int("rtpport", 8082, "The port to listen on for RTP packets")
	trackName := flag.String("track", "default", "The WebRTC name of the track")
	iceServersStr := flag.String("ice", "stun:stun.l.google.com:19302", "Comma-separated list of ICE servers")
	codec := flag.String("codec", "video/VP8", "Codec to report to WebRTC client ('video/VP8' or 'video/H264')")
	serveFiles := flag.Bool("static", true, "Should the static files be served by the app")
	webroot := flag.String("webroot", "./static/", "Location of files to serve")
	apibase := flag.String("apibase", "/webrtc/apiv1/offer", "Location of files to serve")

	flag.Parse()
	iceServers := strings.Split(*iceServersStr, ",")

	InitWebRTCTrack(*codec, *trackName)
	InitWebRTCPApi(iceServers)

	go RunUDPRTPServer(*rtphost, *rtpport)

	if *serveFiles {
		http.Handle("/", http.FileServer(http.Dir(*webroot)))
	}

	http.HandleFunc(*apibase, OnNewPeer)
	log.Panicln(http.ListenAndServe(*haddr, nil))
}
