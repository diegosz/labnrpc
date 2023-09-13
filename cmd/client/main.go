package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/lucasepe/codename"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"labnrpc/provisioning/provisioningpb"
)

func main() {
	log := slog.Default()
	var natsURL = nats.DefaultURL
	if len(os.Args) == 2 {
		natsURL = os.Args[1]
	}

	rng, err := codename.DefaultRNG()
	if err != nil {
		log.Error("codename rng failed", err)
		os.Exit(1)
	}
	name := codename.Generate(rng, 2)

	opts := []nats.Option{
		nats.Name(name),
		nats.CustomInboxPrefix(fmt.Sprintf("%s%s", nats.InboxPrefix, name)),
	}
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Error("nats connection failed", "error", err)
	}
	defer nc.Close()

	// Do this block only if you generated the code with the prometheus plugin.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Error("http subscription failed", err)
		os.Exit(1)
	}
	metricsURL := fmt.Sprintf("http://localhost:%d/metrics", listener.Addr().(*net.TCPAddr).Port)
	fmt.Printf("Check metrics at %s\n", metricsURL)
	http.Handle("/metrics", promhttp.Handler())
	go http.Serve(listener, nil) //nolint:errcheck

	_ = exec.Command("xdg-open", metricsURL).Start()

	client := provisioningpb.NewProvisioningServiceNRPCClient(nc)
	client.Encoding = "json"
	done := make(chan struct{})
	fails := 0
	maxFails := 10

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("\nExiting...")
				return
			case t := <-ticker.C:
				req := &provisioningpb.SayHelloRequest{
					Name: faker.Name(),
				}
				res, err := client.SayHello(req)
				if err != nil {
					log.Error("SayHello failed", "error", err)
					fails++
					if fails > maxFails {
						log.Error("max fails reached, exiting...")
						os.Exit(1)
					}
					continue
				}
				fails = 0
				log.Info("SayHello", "time", t, "request", req.Name, "response", res.Message)
			}
		}
	}()

	fmt.Printf("Client %q is running the provisioning client, Ctrl+C to quit.\n", name)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(done)
	close(c)
}
