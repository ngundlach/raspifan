package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ngundlach/raspifan/sensor"
)

type RestService struct {
	sd  *sensor.SensorData
	Srv *http.Server
}

func NewRestService(sd *sensor.SensorData, addr string) RestService {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", fanHandler(sd))
	mux.HandleFunc("GET /stream", fanStreamHandler(sd))
	mux.HandleFunc("GET /healthcheck", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	srv := &http.Server{Addr: addr, Handler: mux}
	return RestService{sd, srv}
}

func fanHandler(sd *sensor.SensorData) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		d := sd.Get()
		response := sensorModel{}
		response.Fan = string(d.Fan)
		response.Temp = d.Temp
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func fanStreamHandler(sd *sensor.SensorData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		sendT := time.NewTicker(time.Second * 2)
		defer sendT.Stop()

		clientDisconnect := req.Context().Done()

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		}
		data := sd.Get()

		response := sensorModel{Fan: string(data.Fan), Temp: data.Temp}

		sendMessage := func(resp *sensorModel) error {
			jsonData, err := json.Marshal(resp)
			if err != nil {
				log.Printf("Error marshalling JSON %v", err)
				return err
			}
			fmt.Fprintf(w, "event: temp\ndata: %s\n\n", jsonData)
			flusher.Flush()
			return nil
		}

		if err := sendMessage(&response); err != nil {
			log.Printf("Error sending initial message %v", err)
		}
		for {
			select {
			case <-clientDisconnect:
				fmt.Println("Client has disconnected")
				return
			case <-sendT.C:
				response.Fan = string(sd.Get().Fan)
				response.Temp = sd.Get().Temp
				if err := sendMessage(&response); err != nil {
					log.Printf("Error sending message %v", err)
					continue
				}
			}
		}
	}
}
