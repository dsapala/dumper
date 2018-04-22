package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	REQUESTS_DIR = "./requests"
	TIME_FORMAT  = "20060102150405.999999-0700"
	FILE_FORMAT  = "%s-request.txt"
	FILE_MODE    = 0644
)

var addr = flag.String("addr", "127.0.0.1", "the address to bind to")
var port = flag.Int("port", 8080, "the port to bind to")
var readTimeout = flag.Int("read-timeout", 60, "the amount of seconds to wait before timing out on a read from the client")
var writeTimeout = flag.Int("write-timeout", 60, "the amount of seconds to wait before timing out on a write to the client")
var keepalive = flag.Bool("keepalive", true, "whether or not keepalives are enabled")
var keepaliveTimeout = flag.Int("keepalive-timeout", 10, "the amount of seconds to wait before closing the client connection if it becomes idle")
var maxBodyBytes = flag.Int64("max-body-bytes", 10<<20, "the maximum bytes allowed in the body of a request")

func main() {
	flag.Parse()

	sm := http.NewServeMux()

	sm.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// limit body size
		r.Body = http.MaxBytesReader(w, r.Body, *maxBodyBytes)

		t := time.Now()
		ts := t.Format(TIME_FORMAT)
		yearS := t.Format("2006")
		monthS := t.Format("01")
		dayS := t.Format("02")

		dirS := path.Join(REQUESTS_DIR, yearS, monthS, dayS)

		if _, err := os.Stat(dirS); os.IsNotExist(err) {
			os.MkdirAll(dirS, os.ModePerm)
		}

		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf(FILE_FORMAT, path.Join(dirS, ts)), requestDump, FILE_MODE)
		if err != nil {
			log.Println(err)
		}

		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Handler:      sm,
		ReadTimeout:  time.Duration(*readTimeout) * time.Second,
		WriteTimeout: time.Duration(*writeTimeout) * time.Second,
	}

	if *keepalive {
		server.SetKeepAlivesEnabled(true)
		server.IdleTimeout = time.Duration(*keepaliveTimeout) * time.Second
	}

	listen, err := net.Listen("tcp", *addr+":"+strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}

	log.Printf("listening at %s...\n", listen.Addr().String())
	log.Fatal(server.Serve(listen.(*net.TCPListener)))
}
