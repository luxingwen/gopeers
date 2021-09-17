package main

import (
	"flag"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	mirror   = flag.String("mirror", "", "Mirror Web Base URL")
	logfile  = flag.String("log", "-", "Set log file, default STDOUT")
	upstream = flag.String("upstream", "", "Server base URL, conflict with -mirror")
	address  = flag.String("addr", ":5000", "Listen address")
	token    = flag.String("token", "1234567890ABCDEFG", "slave and master token should be same")
)

func InitSignal() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for {
			s := <-sig
			fmt.Println("Got signal:", s)
			if state.Closed {
				fmt.Println("Cold close !!!")
				os.Exit(1)
			}
			fmt.Println("Warm close, waiting ...")
			go func() {
				state.Close()
				os.Exit(0)
			}()
		}
	}()
}

func main() {
	flag.Parse()

	if *upstream != "" {
		if err := InitSlave(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := InitMaster(); err != nil {
			log.Fatal(err)
		}
	}

	InitSignal()

	log.Printf("Listening on %s", *address)
	log.Fatal(http.ListenAndServe(*address, nil))
}
