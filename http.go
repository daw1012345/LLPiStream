package main

import (
	"encoding/json"
	"github.com/pion/webrtc/v3"
	"log"
	"net/http"
)

var iceServers []string

func InitWebRTCPApi(icesrvs []string) {
	iceServers = icesrvs
}

func OnNewPeer(w http.ResponseWriter, r *http.Request) {
	// Only allow one HTTP method to call this API
	if r.Method != "POST" {
		log.Println("Invalid HTTP method!")
		http.Error(w, "Invalid HTTP method!", http.StatusBadRequest)
		return
	}

	// Check that the offer sent to us is not malformed and that we can continue
	var offer webrtc.SessionDescription
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		http.Error(w, "Bad session description!", http.StatusBadRequest)
		return
	}

	// Create a peer and add the video track to it so
	peerConnection := MakePeer()

	// Initialise the object with the offer we got from the peer (to negotiate
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Panicln(err)
	}

	// Apply negotiation algo to decide on communication/what tracks get sent
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Initialise peer object with our reply to the offer presented by peer
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Wait for us to check whether/how we can communicate with peer (via ICE)
	<-gatherComplete

	// Send our reply to peer as via HTTP
	response, err := json.Marshal(peerConnection.LocalDescription())
	if err != nil {
		log.Panicln(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		log.Panicln(err)
	}
}

func MakePeer() *webrtc.PeerConnection {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: iceServers,
			},
		},
	})

	if err != nil {
		log.Panicln(err)
	}

	rtpSender, err := peerConnection.AddTrack(track)
	if err != nil {
		log.Panicln(err)
	}

	go InfiniteRTCPRead(rtpSender)

	return peerConnection
}

func InfiniteRTCPRead(client *webrtc.RTPSender) {
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, rtcpErr := client.Read(rtcpBuf); rtcpErr != nil {
			return
		}
	}
}
