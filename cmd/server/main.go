package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-rpc/nrpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"labnrpc/provisioning/provisioningpb"
)

func main() {
	log := slog.Default()
	var natsURL = nats.DefaultURL
	if len(os.Args) == 2 {
		natsURL = os.Args[1]
	}
	opts := []nats.Option{nats.Name("provisioning-service")}
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Error("nats connection failed", err)
		os.Exit(1)
	}
	defer nc.Close()

	ctx := context.Background()
	pool := nrpc.NewWorkerPool(ctx, uint(runtime.NumCPU()), 10, 4*time.Second)
	h := provisioningpb.NewProvisioningServiceConcurrentHandler(pool, nc, newProvisioningServer(log))
	h.SetEncodings([]string{"protobuf", "json"})

	// Start a NATS subscription using the handler. You can also use the
	// QueueSubscribe() method for a load-balanced set of servers.
	sub, err := nc.Subscribe(h.Subject(), h.Handler)
	if err != nil {
		log.Error("nats subscription failed", err)
		os.Exit(1)
	}
	defer func() {
		_ = sub.Unsubscribe()
	}()

	// Do this block only if you generated the code with the prometheus plugin.
	metricsURL := fmt.Sprintf("http://localhost:%d/metrics", 6060)
	fmt.Printf("Check metrics at %s\n", metricsURL)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":6060", nil)

	_ = exec.Command("xdg-open", metricsURL).Start()

	fmt.Println("Server is running, Ctrl+C to quit.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	close(c)
}

type provisioningServer struct {
	provisioningpb.UnimplementedProvisioningServiceServer
	log   *slog.Logger
	count atomic.Int64
}

func newProvisioningServer(log *slog.Logger) *provisioningServer {
	return &provisioningServer{
		log: log,
	}
}

func (s *provisioningServer) SayHello(ctx context.Context, req *provisioningpb.SayHelloRequest) (*provisioningpb.SayHelloResponse, error) {
	c := s.count.Add(1)
	s.log.Info("SayHello", "count", c, "request", req.Name)
	return &provisioningpb.SayHelloResponse{Message: "Hello " + req.Name + "!"}, nil
}
