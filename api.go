package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/riltech/streamer"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func GetNewApiServer(port string, stream []streamer.Stream) *ApiServer {
	return &ApiServer{
		PortNumber: port,
		Stream:     stream,
	}
}
func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Contanet-type", "application/json")

	json.NewEncoder(w).Encode(v)

}
func (s *ApiServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/Start", s.StartRtsp).Methods("GET")
	router.HandleFunc("/Stop", s.StopRtsp).Methods("GET")

	handler := cors.AllowAll().Handler(router)

	s.Server = &http.Server{
		Addr:    s.PortNumber,
		Handler: handler,
	}
	s.Server.ListenAndServe()
}

func (s *ApiServer) StartRtsp(w http.ResponseWriter, r *http.Request) {
	var folderInfoList []FolderInfo

	for i, _ := range s.Stream {
		// processCommands := []string{
		// 	"ffmpeg",
		// 	"-fflags",
		// 	"nobuffer",
		// 	"-rtsp_transport",
		// 	"tcp",
		// 	"-i",
		// 	s.Stream[i].CMD.Args[7],
		// 	"-c:v",
		// 	"libx264",
		// 	"-movflags",
		// 	"frag_keyframe+empty_moov",

		// 	"-an",
		// 	"-hls_flags",
		// 	"delete_segments+append_list",
		// 	"-f",
		// 	"segment",
		// 	"-segment_list_flags",
		// 	"live",

		// 	"-segment_time",
		// 	"4",

		// 	"-segment_list_size",
		// 	"3",

		// 	"-segment_format",
		// 	"mpegts",
		// 	"-segment_list",
		// 	s.Stream[i].CMD.Args[31],
		// 	"-segment_list_type",
		// 	"m3u8",
		// 	"-segment_list_entry_prefix",
		// 	s.Stream[i].CMD.Args[30],
		// }

		// s.Stream[i].CMD.Args = processCommands

		s.Stream[i].Start().Wait()

		// logrus.Infof("folder name stream  %s | ", s.Stream[i].Logger)
		fmt.Println(s.Stream[i].CMD.Args)

		folderInfo := FolderInfo{
			FolderName: s.Stream[i].ID,
			Url:        s.Stream[i].OriginalURI,
			Count:      "0",
			MacAdd:     find(s.Stream[i].OriginalURI),
		}
		folderInfoList = append(folderInfoList, folderInfo)

	}

	logrus.Infof("folder name stream  %s | ", folderInfoList)
	WriteJson(w, http.StatusOK, folderInfoList)
}
func (s *ApiServer) StopRtsp(w http.ResponseWriter, r *http.Request) {

	// s.Stream.Stop()
	WriteJson(w, http.StatusOK, "ok")

}
func (c *ApiServer) ExitPreHook() chan bool {
	done := make(chan bool)
	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		for uri, strm := range c.Stream {
			logrus.Debugf("Closing processing of %s", uri)
			if err := strm.Stop(); err != nil {
				logrus.Error(err)
				return
			}
			logrus.Debugf("Succesfully closed processing for %s", uri)
		}
		done <- true
	}()
	return done
}
