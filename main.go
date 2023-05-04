package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/riltech/streamer"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var PATH string = "/feeds/video"

func InitConfig() *Specification {
	var s Specification
	err := envconfig.Process("RTSP_STREAM", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	setting := EndpointYML{}
	defer func() {
		s.EndpointYML = setting
	}()
	dat, err := ioutil.ReadFile("rtsp-stream.yml")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &s
	}
	err = yaml.Unmarshal(dat, &setting)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &s
	}
	return &s
}

func main() {

	os.RemoveAll(PATH)
	time.Sleep(time.Second * 1)
	ensureDir(PATH)
	time.Sleep(time.Second * 1)

	spec := InitConfig()
	var streaml []streamer.Stream

	for _, item := range spec.EndpointYML.Listen {
		stream, _ := streamer.NewStream(
			item.Uri,
			PATH,
			true,
			true,
			streamer.ProcessLoggingOpts{
				Enabled:    true,
				Compress:   true,
				Directory:  PATH,
				MaxAge:     10,
				MaxBackups: 100,
				MaxSize:    1000,
			},
			25*time.Second,
		)
		streaml = append(streaml, *stream)
	}

	server := GetNewApiServer(":8081", streaml)




	done := server.ExitPreHook()
	go func() {
		logrus.Println("rtsp-stream transcoder started on %d | MainProcess")

		server.Run()
	}()
	<-done
	if err := server.Server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("HTTP server Shutdown: %v", err)
	}
	os.Exit(0)

}
