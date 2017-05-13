package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
)

func main() {
	// Parse flags
	flag.Parse()
	log.Println("Start")

	// Set up logger
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))
	templateData, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatalln("Open template err:", err)
	}
	// Start an http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(templateData)
	})
	go func() {
		if err := http.ListenAndServe("127.0.0.1:4000", nil); err != nil {
			log.Fatalln("http.ListenAndServe err:", err)
		}
	}()

	resoursesPath := os.Getenv("GOPATH") + "/src/github.com/asticode/go-astilectron/examples"
	a, err := astilectron.New(astilectron.Options{
		AppName:            "Electron demo",
		AppIconDefaultPath: resoursesPath + "/gopher.png",
		AppIconDarwinPath:  resoursesPath + "/gopher.icns",
		BaseDirectoryPath:  resoursesPath,
	})
	if err != nil {
		astilog.Fatal("creating new astilectron failed err:", err)
	}
	defer a.Close()
	a.HandleSignals()
	// Start
	log.Println("Start electron")
	if err = a.Start(); err != nil {
		astilog.Fatal("starting failed, err:", err)
	}

	w, err := a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
	// Center: astilectron.PtrBool(true),
	// Height: astilectron.PtrInt(600),
	// Width:  astilectron.PtrInt(600),
	})
	if err != nil {
		astilog.Fatal("new window failed, err", err)
	}
	if err = w.Create(); err != nil {
		astilog.Fatal("creating window failed, err:", err)
	}
	log.Println("Window create success")

	// Add listener
	w.On(astilectron.EventNameWindowEventMessage, func(e astilectron.Event) (deleteListener bool) {
		var m string
		e.Message.Unmarshal(&m)
		astilog.Infof("Received message %s", m)
		return
	})
	if err = w.Maximize(); err != nil {
		astilog.Fatal("maximizing window failed", err)
	}
	time.Sleep(1 * time.Second)
	w.Send("Ping!")
	// Blocking pattern
	a.Wait()
}
