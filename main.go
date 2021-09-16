
package main

import (
	"net/http"
)

func main() {
	Init()

	go StartUDPRTPServer()

	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/api/start", OnNewPeer)
	panic(http.ListenAndServe(":8080", nil))
}