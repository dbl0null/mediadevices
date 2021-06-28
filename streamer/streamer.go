package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	"github.com/getlantern/systray"
	"github.com/pion/webrtc/v3"
	"github.com/skratchdot/open-golang/open"

	md "github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/io/video"

	// h264 "github.com/pion/mediadevices/pkg/codec/x264"
	h264 "github.com/pion/mediadevices/pkg/codec/openh264"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/mediadevices/pkg/driver"
	"github.com/pion/mediadevices/pkg/prop"
)

const (
	mtu = 1000
)

func must(err error) bool {
	if err != nil {
		//fmt.Printf("ERROR: %s\n", err.Error())
		return true
	}
	return false
}

type Stream struct {
	Addr     string
	Active   bool
	DeviceId string
	Started  time.Time
	Sent     uint64
	Stop     chan interface{}
}

func (s Stream) String() string {
	return fmt.Sprintf("%s active: %t, device id: %s, start time: %s, bytes sent: %s", s.Addr, s.Active, s.DeviceId, strconv.FormatInt(s.Started.Unix(), 10), strconv.FormatUint(s.Sent, 10))
}

func (s Stream) Serialize() string {
	return fmt.Sprintf("s/%s/%t/%s/%s/%s", s.Addr, s.Active, s.DeviceId, strconv.FormatInt(s.Started.Unix(), 10), strconv.FormatUint(s.Sent, 10))
}

var Streams = make(map[string]*Stream)
var Devices = make(map[string]*md.MediaDeviceInfo)
var Menu = make(map[string]*systray.MenuItem)
var mVideo *systray.MenuItem
var mAudio *systray.MenuItem
var mStreams *systray.MenuItem

func init() {
}

func State() string {
	fmt.Printf("STATE:\n\n")

	state := ""
	for _, stream := range Streams {
		if stream != nil {
			state += stream.Serialize()
			state += "\n"
		}
		fmt.Printf("%s\n", state)
	}
	for _, device := range Devices {
		if device != nil {
			state += device.Label
			state += "\n"
		}
	}

	return state
}

func Enumerate() map[string]*md.MediaDeviceInfo {
	fmt.Printf("ENUMERATING:\n")

	drivers := driver.GetManager().Query(func(drv driver.Driver) bool {
		return true
	})
	devices := make(map[string]*md.MediaDeviceInfo, len(drivers))

	for _, drv := range drivers {
		var kind md.MediaDeviceType
		deviceID := drv.ID()
		drvInfo := drv.Info()

		switch {
		case driver.FilterVideoRecorder()(drv):
			mVideoSub := mVideo.AddSubMenuItem(drvInfo.Name, drvInfo.Label)
			Menu[deviceID] = mVideoSub
			kind = md.VideoInput
			if drvInfo.Name == "HD Pro Webcam C920" {
				go Start("192.168.1.20:15000", deviceID)
			}
		case driver.FilterAudioRecorder()(drv):
			Menu[deviceID] = mAudio.AddSubMenuItem(drvInfo.Name, drvInfo.Label)
			kind = md.AudioInput
		default:
			continue
		}

		deviceInfo := md.MediaDeviceInfo{DeviceID: deviceID, Kind: kind, Label: drvInfo.Label, Name: drvInfo.Name, DeviceType: drvInfo.DeviceType}
		Devices[deviceID] = &deviceInfo

		fmt.Printf("\t%s\n", deviceInfo.String())

		// if driver.FilterVideoRecorder()(drv) && sending == "" {
		// 	sending = deviceID
		// }
	}

	// TODO: check if device not available anymore
	Devices = devices

	// go Start("192.168.1.20:15000", sending)

	return devices
}

func Start(addr, deviceId string) {
	Menu[deviceId].Check()

	device, ok := Devices[deviceId]
	if !ok {
		fmt.Printf("[%s->%s] no such device\n", deviceId, addr)
		return
	}

	a, err := net.ResolveUDPAddr("udp", addr)
	if must(err) {
		fmt.Printf("[%s->%s] bad addr\n", deviceId, addr)
		return
	}

	conn, err := net.DialUDP("udp", nil, a)
	if must(err) {
		fmt.Printf("[%s->%s] can't connect to addr\n", deviceId, addr)
		return
	}

	constraints := md.MediaStreamConstraints{}

	var codecName string
	var payloadType uint8

	switch device.Kind {
	case md.VideoInput:
		h264Params, err := h264.NewParams()
		if must(err) {
			fmt.Printf("[%s->%s] can't make x264 params\n", deviceId, addr)
			return
		}
		h264Params.BitRate = 3_000_000
		//h264Params.Preset = x264.PresetUltrafast
		h264Params.KeyFrameInterval = 1

		constraints.Codec = md.NewCodecSelector(md.WithVideoEncoders(&h264Params))
		constraints.Video = func(c *md.MediaTrackConstraints) {
			c.FrameFormat = prop.FrameFormatOneOf{frame.FormatI420, frame.FormatYUY2}
			c.DeviceID = prop.StringExact(deviceId)
			// c.Width = prop.IntRanged{Min: 640, Max: 1920, Ideal: 1280}
			c.Width = prop.Int(640)
			c.Height = prop.Int(480)
		}

		codecName = webrtc.MimeTypeH264
		payloadType = 125 // corresponding value from sdpOffer to Kurento
	case md.AudioInput:
		opusParams, err := opus.NewParams()
		must(err)
		if must(err) {
			fmt.Printf("[%s->%s] can't make Opus params\n", deviceId, addr)
			return
		}
		opusParams.BitRate = 256_000

		constraints.Audio = func(c *md.MediaTrackConstraints) { c.DeviceID = prop.StringExact(deviceId) }
		constraints.Codec = md.NewCodecSelector(md.WithAudioEncoders(&opusParams))

		codecName = webrtc.MimeTypeOpus
		payloadType = 96
	default:
		return
	}

	mediaStream, err := md.GetUserMedia(constraints)
	if must(err) {
		fmt.Printf("[%s->%s] can't get media stream: %s\n", deviceId, addr, err.Error())
		return
	}

	track := mediaStream.GetTracks()[0]
	defer track.Close()

	rtpReader, err := track.NewRTPReader(codecName, rand.Uint32(), mtu) //nolint:gosec
	if must(err) {
		fmt.Printf("[%s->%s] can't make rtp reader: %s\n", deviceId, addr, err.Error())
		return
	}
	defer rtpReader.Close()

	ticker := time.NewTicker(time.Millisecond * 250)
	defer ticker.Stop()

	stop := make(chan interface{})
	defer close(stop)

	var videoReader video.Reader

	if track.Kind() == webrtc.RTPCodecTypeVideo {
		videoReader = track.(*md.VideoTrack).NewReader(true)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			mimeWriter := multipart.NewWriter(w)

			contentType := fmt.Sprintf("multipart/x-mixed-replace;boundary=%s", mimeWriter.Boundary())
			w.Header().Add("Content-Type", contentType)

			partHeader := make(textproto.MIMEHeader)
			partHeader.Add("Content-Type", "image/jpeg")

			for {
				frame, release, err := videoReader.Read()
				if err == io.EOF {
					return
				}
				must(err)

				err = jpeg.Encode(&buf, frame, nil)
				// Since we're done with img, we need to release img so that that the original owner can reuse
				// this memory.
				release()
				must(err)

				partWriter, err := mimeWriter.CreatePart(partHeader)
				must(err)

				_, err = partWriter.Write(buf.Bytes())
				buf.Reset()
				must(err)
			}
		})

	}

	stream := Stream{DeviceId: deviceId, Addr: addr, Started: time.Now(), Stop: stop}
	Streams[addr] = &stream
	mStreamItem := mStreams.AddSubMenuItem(addr, device.Name)
	mStreamItemStop := mStreamItem.AddSubMenuItem("Stop", "Stop")
	mStreamItemPreview := mStreamItem.AddSubMenuItem("Preview", "Preview")
	Menu[addr] = mStreamItem

	fmt.Printf("[%s->%s] new stream: %s\n", deviceId, addr, stream.String())
	go http.ListenAndServe(":8877", nil)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			if videoReader != nil {
	// 				frm_, release, err := videoReader.Read()
	// 				must(err)
	// 				frm := frm_.(*image.YCbCr)
	// 				bounds := frm.Bounds()
	// 				cascadeParams := pigo.CascadeParams{
	// 					MinSize:     100,
	// 					MaxSize:     600,
	// 					ShiftFactor: 0.15,
	// 					ScaleFactor: 1.1,
	// 					ImageParams: pigo.ImageParams{
	// 						Pixels: frm.Y, // Y in YCbCr should be enough to detect faces
	// 						Rows:   bounds.Dy(),
	// 						Cols:   bounds.Dx(),
	// 						Dim:    bounds.Dx(),
	// 					},
	// 				}
	//
	// 				// Run the classifier over the obtained leaf nodes and return the detection results.
	// 				// The result contains quadruplets representing the row, column, scale and detection score.
	// 				dets := classifier.RunCascade(cascadeParams, 0.0)
	//
	// 				// Calculate the intersection over union (IoU) of two clusters.
	// 				dets = classifier.ClusterDetections(dets, 0)
	//
	// 				for _, det := range dets {
	// 					if det.Q >= confidenceLevel {
	// 						log.Println("Detect a face")
	// 					}
	// 				}
	//
	// 				release()
	// 			}
	// 		}
	// 	}
	// }()

	buff := make([]byte, mtu)

	for {
		select {
		case <-mStreamItemPreview.ClickedCh:
			open.Run("http://localhost:8877")
		case <-mStreamItemStop.ClickedCh:
			close(stop)
		case <-stop:
			stream.Active = false
			fmt.Printf("[%s->%s] stop stream: %s\n", deviceId, addr, stream.String())
			return
		default:
			pkts, release, err := rtpReader.Read()
			if must(err) {
				return
			}

			stream.Active = true

			for _, pkt := range pkts {
				pkt.PayloadType = payloadType
				n, err := pkt.MarshalTo(buff)
				if must(err) {
					continue
				}

				b, err := conn.Write(buff[:n])
				if must(err) {
					continue
				}
				stream.Sent += uint64(b)
				//fmt.Printf("sent: %d\n", b)
			}
			release()
		}
	}
}

func Stop(addr string) {
	stream, ok := Streams[addr]
	if ok {
		close(stream.Stop)
	}
}
