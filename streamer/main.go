package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	_ "github.com/pion/mediadevices/pkg/driver/camera"
	_ "github.com/pion/mediadevices/pkg/driver/microphone"

	"github.com/denisbrodbeck/machineid"
)

var Uid string

func onExit() {
	now := time.Now()
	ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}

func main() {
	log.SetFlags(0)
	var err error

	Uid, err = machineid.ProtectedID(os.Args[0])
	if err != nil {
		Uid = "temp-" + uuid.New().String()
	}

	systray.Run(onReady, onExit)

	fmt.Printf("\n\nApplication ID: %v\n\n", Uid)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	addr := "phobia.bsjsc.ru"
	// TODO: ws host as parameter, adjustable RTP payload types, bitrate and resolution

	u := url.URL{Scheme: "wss", Host: addr, Path: "streamers/" + Uid}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		log.Printf("connecting to %s", u.String())

		ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Printf("dial: %s", err.Error())
			continue
		}
		// defer ws.Close()

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				log.Println("read loop")
				t, data, err := ws.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					continue
				}
				log.Printf("recv [%d]: %s", t, data)

				if data[0] == '!' {
					args := strings.SplitN(string(data[1:]), "/", 3)
					if len(args) > 1 {
						cmd := args[0]
						addr := args[1]
						switch cmd {
						case "start":
							deviceId := args[2]
							fmt.Printf("[%s] %s to udp://%s\n", cmd, deviceId, addr)
							go Start(addr, deviceId)
						case "stop":
							fmt.Printf("[%s] udp://%s\n", cmd, addr)
							Stop(addr)
						}
					}
				} else {
					log.Printf("[%s]\n", data)
				}
			}
		}()

	loop:
		for {
			log.Println("loop")

			select {
			case <-done:
				log.Println("done")
				break
			case <-ticker.C:
				log.Println("tick")
				if err := ws.WriteMessage(websocket.TextMessage, []byte(State())); err != nil {
					log.Println("write:", err)
					break loop
				}
			case <-interrupt:
				log.Println("interrupt")
				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				if err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
					log.Println("write close:", err)
					break loop
				}
				select {
				case <-done:
				case <-time.After(time.Second * 30):
				}
				break loop
			}
		}
		fmt.Println("restarting in 10s...")
		time.Sleep(time.Second * 10)
	}
}
