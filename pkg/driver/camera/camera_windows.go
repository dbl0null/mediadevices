package camera

/*
#cgo CXXFLAGS: -std=gnu++11
#cgo LDFLAGS: -lstrmiids -lole32 -loleaut32 -lquartz
#include <dshow.h>
#include "camera_windows.hpp"
*/
import "C"

import (
	"fmt"
	"image"
	"io"
	"sync"
	"unsafe"

	"github.com/pion/mediadevices/pkg/driver"
	"github.com/pion/mediadevices/pkg/frame"
	"github.com/pion/mediadevices/pkg/io/video"
	"github.com/pion/mediadevices/pkg/prop"
)

var (
	callbacks   = make(map[uintptr]*Camera)
	callbacksMu sync.RWMutex
)

type Camera struct {
	UID   string
	Name  string
	cam   *C.camera
	ch    chan []byte
	buf   []byte
	bufGo []byte
}

func init() {
	C.CoInitializeEx(nil, C.COINIT_MULTITHREADED)

	var list C.cameraList
	var errStr *C.char
	if C.listCamera(&list, &errStr) != 0 {
		fmt.Printf("Failed to list camera: %s\n", C.GoString(errStr))
		return
	}
	fmt.Printf("listCamera: %d\n", list.num)
	for i := 0; i < int(list.num); i++ {
		uid := C.GoString(C.getUid(&list, C.int(i)))
		name := C.GoString(C.getName(&list, C.int(i)))
		cam := &Camera{UID: uid, Name: name}
		driver.GetManager().Register(cam, driver.Info{
			Label:      uid,
			Name:       name,
			DeviceType: driver.Camera,
		})
	}

	C.freeCameraList(&list, &errStr)
}

func (cam *Camera) Open() error {
	cam.ch = make(chan []byte)
	cam.cam = &C.camera{name: C.CString(cam.Name), uid: C.CString(cam.UID)}

	var errStr *C.char
	if C.listResolution(cam.cam, &errStr) != 0 {
		return fmt.Errorf("failed to open device: %s", C.GoString(errStr))
	}

	return nil
}

//export imageCallback
func imageCallback(cam uintptr) {
	callbacksMu.RLock()
	cb, ok := callbacks[uintptr(unsafe.Pointer(cam))]
	callbacksMu.RUnlock()
	if !ok {
		return
	}

	copy(cb.bufGo, cb.buf)
	cb.ch <- cb.bufGo
}

func (cam *Camera) Close() error {
	callbacksMu.Lock()
	key := uintptr(unsafe.Pointer(cam.cam))
	if _, ok := callbacks[key]; ok {
		delete(callbacks, key)
	}
	callbacksMu.Unlock()
	close(cam.ch)

	if cam.cam != nil {
		C.free(unsafe.Pointer(cam.cam.name))
		C.freeCamera(cam.cam)
		cam.cam = nil
	}
	return nil
}

func (cam *Camera) VideoRecord(p prop.Media) (video.Reader, error) {
	nPix := p.Width * p.Height
	cam.buf = make([]byte, nPix*2) // for YUY2
	cam.bufGo = make([]byte, nPix*2)
	cam.cam.width = C.int(p.Width)
	cam.cam.height = C.int(p.Height)
	cam.cam.buf = C.size_t(uintptr(unsafe.Pointer(&cam.buf[0])))

	var errStr *C.char
	if C.openCamera(cam.cam, &errStr) != 0 {
		return nil, fmt.Errorf("failed to open device: %s", C.GoString(errStr))
	}

	callbacksMu.Lock()
	callbacks[uintptr(unsafe.Pointer(cam.cam))] = cam
	callbacksMu.Unlock()

	img := &image.YCbCr{}

	r := video.ReaderFunc(func() (image.Image, func(), error) {
		b, ok := <-cam.ch
		if !ok {
			return nil, func() {}, io.EOF
		}
		img.Y = b[:nPix]
		img.Cb = b[nPix : nPix+nPix/2]
		img.Cr = b[nPix+nPix/2 : nPix*2]
		img.YStride = p.Width
		img.CStride = p.Width / 2
		img.SubsampleRatio = image.YCbCrSubsampleRatio422
		img.Rect = image.Rect(0, 0, p.Width, p.Height)
		return img, func() {}, nil
	})
	return r, nil
}

func (cam *Camera) Properties() []prop.Media {
	properties := []prop.Media{}
	for i := 0; i < int(cam.cam.numProps); i++ {
		p := C.getProp(cam.cam, C.int(i))
		// TODO: support other FOURCC
		if p.fcc == fourccYUY2 {
			properties = append(properties, prop.Media{
				Video: prop.Video{
					Width:       int(p.width),
					Height:      int(p.height),
					FrameFormat: frame.FormatYUY2,
				},
			})
		}
	}
	return properties
}

const (
	fourccYUY2 = 0x32595559
)
