package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lucasepe/codename"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nuid"
	"github.com/nats-rpc/nrpc"

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
	inboxPrefix := fmt.Sprintf("%s%s", nats.InboxPrefix, name)
	inboxPrefixLen := len(inboxPrefix)
	nuidSize := 22
	inboxLen := inboxPrefixLen + nuidSize

	opts := []nats.Option{
		nats.Name(name),
		nats.CustomInboxPrefix(inboxPrefix),
	}
	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Error("nats connection failed", "error", err)
	}
	defer nc.Close()

	client := provisioningpb.NewProvisioningServiceNRPCClient(nc)

	// nrpc.GetReplyInbox is used by StreamCall to get a inbox subject. It can
	// be changed by a client lib that needs custom inbox subjects.
	nrpc.GetReplyInbox = func(_ nrpc.NatsConn) string {
		b := make([]byte, inboxLen)
		pres := b[:inboxPrefixLen]
		copy(pres, inboxPrefix)
		ns := b[inboxPrefixLen:]
		copy(ns, nuid.Next())
		return string(b[:])
	}

	if err := client.SayHelloPoll(&provisioningpb.SayHelloRequest{Name: name}, 20,
		func(r *provisioningpb.SayHelloResponse) error {
			m := r.GetMessage()
			fmt.Println(">>>>>>>>>>>>>>", m)
			if r.Message != "Hello "+name+"!" {
				return fmt.Errorf("invalid response: %s", m)
			}
			log.Info("Poll SayHello", "request", name, "response", m)
			return nil
		},
	); err != nil {
		log.Error("Poll SayHello failed, exiting...", "error", err)
		os.Exit(1)
	}
}
