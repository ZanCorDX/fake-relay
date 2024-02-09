package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fb/httpserver/common"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Sample code to receive a block submission EXACTLY as our relay does now.
// Based on https://github.com/flashbots/mev-boost-relay taking code from:
// mev-boost-relay/services/api/service.go handleSubmitNewBlock for the parsing
// mev-boost-relay/common/types_spec.go
func main() {
	var	pathSubmitNewBlock = "/relay/v1/builder/blocks"
    // Create a new Gorilla Mux router.
    r := mux.NewRouter()
    // Define your API routes and handlers.
	r.HandleFunc(pathSubmitNewBlock, handleSubmitNewBlock).Methods(http.MethodPost);
	

    // Create an HTTP server with Gorilla Mux router.
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }

    // Start the server in a separate Goroutine.
    go func() {
        log.Println("Starting the server on :8080")
        if err := srv.ListenAndServe(); err != nil {
            log.Fatal(err)
        }
    }()

    // Implement graceful shutdown.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down the server...")

    // Set a timeout for shutdown (for example, 5 seconds).
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown error: %v", err)
    }
    log.Println("Server gracefully stopped")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome to the GoLang HTTP Server!")
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    fmt.Fprintf(w, "Hello, %s!", name)
}

func handleSubmitNewBlock(w http.ResponseWriter, req *http.Request) {
	var r io.Reader = req.Body
	var err error
	isGzip := req.Header.Get("Content-Encoding") == "gzip"
	if isGzip {
		r, err = gzip.NewReader(req.Body)
		if err != nil {
			fmt.Printf("api.RespondError(w, http.StatusBadRequest:", err.Error());
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	limitReader := io.LimitReader(r, 10*1024*1024) // 10 MB
	requestPayloadBytes, err := io.ReadAll(limitReader)
	if err != nil {
		fmt.Printf("api.RespondError(w, http.StatusBadRequest:", err.Error());
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload := new(common.VersionedSubmitBlockRequest)
	// Check for SSZ encoding
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/octet-stream" {
		if err = payload.UnmarshalSSZ(requestPayloadBytes); err != nil {
			fmt.Printf("could not decode payload - SSZ");

			// SSZ decoding failed. try JSON as fallback (some builders used octet-stream for json before)
			if err2 := json.Unmarshal(requestPayloadBytes, payload); err2 != nil {
				fmt.Printf("api.RespondError(w, http.StatusBadRequest;", err.Error());
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			fmt.Printf("received ssz-encoded payload")
		}
	} else {
		if err := json.Unmarshal(requestPayloadBytes, payload); err != nil {
			fmt.Printf("api.RespondError(w, http.StatusBadRequest: ",err.Error());
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	fmt.Printf("Packet deserialized!!\n");
	
	var reference_file_name = "./deneb_test_packet.bin"
	//Set generate_reference to true the first time to generate reference packets that you can compare later setting generate_reference to false
	generate_reference := false
	if generate_reference  {
		marshaled,_ := json.Marshal(payload)
		os.WriteFile(reference_file_name,marshaled,0644);
	} else {
		marshaled,_ := json.Marshal(payload)
		reference_marshaled,_ := os.ReadFile(reference_file_name)
	
		if bytes.Equal(marshaled,reference_marshaled) {
			fmt.Printf("OK!!");	
		} else {
			fmt.Printf("FAILED :()");
		}
		fmt.Printf("\n=============\n");
	}

	w.WriteHeader(http.StatusOK)
}