package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flagMode := flag.String("mode", "client", "app working mode client or server")
	flag.Parse()

	conf, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatalf("can't read the config: %s", err)
	}

	if *flagMode == "server" {

		go func() {
			if err := NewServerStart(conf.GRPC); err != nil {
				log.Fatalf("couldn't launch the server: %s", err)
			}
		}()
		if err := NewHTTPServerStart(conf.GRPC, conf.HTTP); err != nil {
			log.Fatalf("couldn't launch the HTTP proxy: %s", err)
		}

	} else if *flagMode == "client" {

		client, err := NewClient(conf.GRPC)
		if err != nil {
			log.Fatalf("couldn't connect to the server: %s", err)
		}
		fmt.Println(client.ServerInfo())

	} else {
		fmt.Printf("unknown mode %s, server or client allowed", *flagMode)
	}

}
