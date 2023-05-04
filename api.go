package main

import (
	"encoding/json"
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

	for i, _ := range s.Stream {
		s.Stream[i].Start().Wait()
	
		logrus.Infof("folder name stream  %s | ", s.Stream[i].ID)
		logrus.Infof("folder name stream  %s | ", s.Stream[i].Logger)
			
		WriteJson(w, http.StatusOK, s.Stream[i].ID)
	}

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
