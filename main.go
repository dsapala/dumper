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

func main() {
	flag.Parse()

	sm := http.NewServeMux()

	sm.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		Handler: sm,
	}

	a := server.Addr
	if a == "" {
		a = *addr + ":" + strconv.Itoa(*port)
	}

	listen, err := net.Listen("tcp", a)
	if err != nil {
		panic(err)
	}

	address := listen.Addr().String()
	log.Printf("listening at %s...\n", address)
	log.Fatal(server.Serve(listen.(*net.TCPListener)))
}
