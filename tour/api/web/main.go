package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jaugustosaba/ebury-exam/tour"
	"github.com/jaugustosaba/ebury-exam/tour/api"
)

var (
	port = flag.Int("port", 8080, "server port")
	// using a mutex to avoid concurrent access on our in-memory tour implementation
	mutex = sync.Mutex{}
)

func loadCSV(t *tour.Tour, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return t.LoadFromCSV(file)
}

func createServiceHandler[Request any, Response any](httpMethod string, service func(Request) Response) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s", req.Method, req.RequestURI)

		if req.Method != httpMethod {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if req.Header.Get("Content-Type") != "application/json" {
			resp.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		var jsonReq Request

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&jsonReq)
		if err != nil {
			log.Printf("cannot decode request body: %s", err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		mutex.Lock()
		jsonResp := (func() Response {
			defer mutex.Unlock()
			return service(jsonReq) // mutex protected service invocation
		})()

		resp.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(resp)
		err = encoder.Encode(jsonResp)
		if err != nil {
			log.Printf("cannot encode response body: %s", err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	flag.Parse()

	// server's tour
	t := tour.NewTour()

	// loads CSV
	for _, path := range flag.Args() {
		err := loadCSV(t, path)
		if err != nil {
			log.Printf("cannot read CSV file %s: %s", path, err.Error())
		}
	}

	// start REST server
	service := api.NewService(t)
	http.HandleFunc("/route/add", createServiceHandler(http.MethodPost, service.AddRoute))
	http.HandleFunc("/route/shortest", createServiceHandler(http.MethodPost, service.ShortestRoute))
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Printf("cannot open server: %s", err.Error())
		os.Exit(1)
	}
}
