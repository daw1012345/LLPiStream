const defaultRtcConfig = new RTCPeerConnection({iceServers: [{urls: "stun:stun.l.google.com:19302"}]})

export function makeLivestream(videoElement, rtcConfig = defaultRtcConfig, endpoint = "/webrtc/apiv1/offer", replyHandler = null, errorHandler = null) {
    let pc = new RTCPeerConnection(rtcConfig)

    pc.ontrack = function (event) {
        var el = videoElement
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true
        el.muted = true
    }

    pc.addTransceiver('video', {'direction': 'recvonly'})
    pc.createOffer()
        .then(offer => {
            pc.setLocalDescription(offer)
            return fetch(endpoint, {
                method: 'post',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(offer)
            })
        })
        .then(res => res.json())
        .then(res => pc.setRemoteDescription(new RTCSessionDescription(res)))
        .then(replyHandler)
        .catch(err => errorHandler(err))
}