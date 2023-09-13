package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-rpc/nrpc"

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
	// this is only for testing, it's not the way to do it
	// you have to create a request that is going to stream responses...
	resp := &provisioningpb.SayHelloResponse{Message: "Hello " + req.Name + "!"}
	for i := 0; i < 19; i++ {
		s.log.Info("SayHello", "count", s.count.Add(1), "request", req.Name)
		r := nrpc.GetRequest(ctx)
		if err := r.SendReply(resp.ProtoReflect().Interface(), nil); err != nil {
			return nil, err
		}
		time.Sleep(200 * time.Millisecond)
	}
	s.log.Info("SayHello", "count", s.count.Add(1), "request", req.Name)
	return resp, nil
}
