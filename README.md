# L(ow)L(atency)PiStream

Low-profile low-latency WebRTC streaming project.

## Rationale
I was working on a Raspberry Pi-based project and needed a very basic low latency video stream.
All solutions I came across had far too many dependencies or features which made integrating them into my project too difficult.

This project is meant to:
* be extremely simple
* provide a baseline for basic projects 
* demonstrate a basic WebRTC setup


## Installation
### From source
1. `> git clone https://github.com/daw1012345/LLPiStream`
2. `> cd LLPiStream`
3. `> GOBIN=/usr/local/bin go install`
### Via `go install` into `/usr/local/bin`
1. `> GOBIN=/usr/local/bin go install https://github.com/daw1012345/LLPiStream@latest`

## Usage
```> go run llpistream --help
Usage of llpistream:
    -apibase string
        Location of files to serve (default "/webrtc/apiv1/offer")
    -codec string
        Codec to report to WebRTC client ('video/VP8' or 'video/H264') (default "video/VP8")
    -haddr string
        The host/port to run the API server on (default "127.0.0.1:8080")
    -ice string
        Comma-separated list of ICE servers (default "stun:stun.l.google.com:19302")
    -rtphost string
        The port to listen on for RTP packets (default "127.0.0.1")
    -rtpport int
        The port to listen on for RTP packets (default 8082)
    -static
        Should the static files be served by the app (default true)
    -track string
        The WebRTC name of the track (default "default")
    -webroot string
        Location of files to serve (default "./static/")
```
## Examples
### VP8 stream
#### Requirements
* `gstreamer`
* `gstreamer-plugins-good` (for rtvp8pay)
* `gstreamer-plugins-base` (for videotestsrc)
#### Instructions
* `llpistream &`
* `gst-launch-1.0 videotestsrc ! "video/x-raw,height=720,width=1280,framerate=30/1" ! queue ! videoscale ! videoconvert ! vp8enc error-resilient=partitions keyframe-max-dist=10 deadline=1 ! rtpvp8pay ! udpsink host=127.0.0.1 port=8082`
* Open `localhost:8080` in a browser

### H264 stream
#### Requirements
* `gstreamer`
* `gstreamer-plugins-good` (for rtph264pay)
* `gstreamer-plugins-ugly` (for x264enc)
* `gstreamer-plugins-base` (for videotestsrc)
* 
#### Instructions
* `> llpistream --codec video/H264 &`
* `> gst-launch-1.0 videotestsrc ! "video/x-raw,height=720,width=1280,framerate=30/1" ! queue ! videoconvert ! videoscale ! queue ! x264enc tune=zerolatency speed-preset=6 ! "video/x-h264,profile=constrained-baseline,width=1280,height=720,stream-format=byte-stream" ! rtph264pay ! udpsink host=127.0.0.1 port=8082`
* Open `localhost:8080` in a browser

### VP8 vs H264

The codec you choose depends on your hardware. Check whether your hardware supports hardware encoding of either codec, or benchmark both by trying the examples.

## Embedding in your own project
It's best to not use the provided `index.html` file (it's just an example, and probably doesn't do what you want it to do anyways).

All you have to do is serve the `livestream.js` file with a web server (Nginx or Apache) and pass requests to the WebRTC endpoints to `llpistream`. 
Then import the `livestream.js` module in your javascript and use the provided `makeLivestream()` function to turn your `<video>` tags into livestreams!

Is `llpistream` missing any features you want? Do you need authentication? Fork it and add them!

## Acknowledgments
Uses:
* [pion](https://github.com/pion/webrtc)

Heavily based on:
* [rtp-to-webrtc](https://github.com/pion/webrtc/tree/master/examples/rtp-to-webrtc)
* [rtmp-to-webrtc](https://github.com/Sean-Der/rtmp-to-webrtc)


