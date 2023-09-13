// This code was autogenerated from provisioningpb/api.proto, do not edit.
package provisioningpb

import (
	"context"
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/nats-rpc/nrpc"
)

// ProvisioningServiceNRPCServer is the interface that providers of the service
// ProvisioningService should implement.
type ProvisioningServiceNRPCServer interface {
	SayHello(ctx context.Context, req *SayHelloRequest) (*SayHelloResponse, error)
}

var (
	// The request completion time, measured at client-side.
	clientRCTForProvisioningService = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "nrpc_client_request_completion_time_seconds",
			Help:       "The request completion time for calls, measured client-side.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			ConstLabels: map[string]string{
				"service": "ProvisioningService",
			},
		},
		[]string{"method"})

	// The handler execution time, measured at server-side.
	serverHETForProvisioningService = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "nrpc_server_handler_execution_time_seconds",
			Help:       "The handler execution time for calls, measured server-side.",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			ConstLabels: map[string]string{
				"service": "ProvisioningService",
			},
		},
		[]string{"method"})

	// The counts of calls made by the client, classified by result type.
	clientCallsForProvisioningService = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nrpc_client_calls_count",
			Help: "The count of calls made by the client.",
			ConstLabels: map[string]string{
				"service": "ProvisioningService",
			},
		},
		[]string{"method", "encoding", "result_type"})

	// The counts of requests handled by the server, classified by result type.
	serverRequestsForProvisioningService = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nrpc_server_requests_count",
			Help: "The count of requests handled by the server.",
			ConstLabels: map[string]string{
				"service": "ProvisioningService",
			},
		},
		[]string{"method", "encoding", "result_type"})
)

// ProvisioningServiceHandler provides a NATS subscription handler that can serve a
// subscription using a given ProvisioningServiceNRPCServer implementation.
type ProvisioningServiceHandler struct {
	ctx     context.Context
	workers *nrpc.WorkerPool
	nc      nrpc.NatsConn
	server  ProvisioningServiceNRPCServer

	encodings []string
}

func NewProvisioningServiceHandler(ctx context.Context, nc nrpc.NatsConn, s ProvisioningServiceNRPCServer) *ProvisioningServiceHandler {
	return &ProvisioningServiceHandler{
		ctx:    ctx,
		nc:     nc,
		server: s,

		encodings: []string{"protobuf"},
	}
}

func NewProvisioningServiceConcurrentHandler(workers *nrpc.WorkerPool, nc nrpc.NatsConn, s ProvisioningServiceNRPCServer) *ProvisioningServiceHandler {
	return &ProvisioningServiceHandler{
		workers: workers,
		nc:      nc,
		server:  s,
	}
}

// SetEncodings sets the output encodings when using a '*Publish' function
func (h *ProvisioningServiceHandler) SetEncodings(encodings []string) {
	h.encodings = encodings
}

func (h *ProvisioningServiceHandler) Subject() string {
	return "ProvisioningService.>"
}

func (h *ProvisioningServiceHandler) Handler(msg *nats.Msg) {
	var ctx context.Context
	if h.workers != nil {
		ctx = h.workers.Context
	} else {
		ctx = h.ctx
	}
	request := nrpc.NewRequest(ctx, h.nc, msg.Subject, msg.Reply)
	// extract method name & encoding from subject
	_, _, name, tail, err := nrpc.ParseSubject(
		"", 0, "ProvisioningService", 0, msg.Subject)
	if err != nil {
		log.Printf("ProvisioningServiceHanlder: ProvisioningService subject parsing failed: %v", err)
		return
	}

	request.MethodName = name
	request.SubjectTail = tail

	// call handler and form response
	var immediateError *nrpc.Error
	switch name {
	case "SayHello":
		_, request.Encoding, err = nrpc.ParseSubjectTail(0, request.SubjectTail)
		if err != nil {
			log.Printf("SayHelloHanlder: SayHello subject parsing failed: %v", err)
			break
		}
		var req SayHelloRequest
		if err := nrpc.Unmarshal(request.Encoding, msg.Data, &req); err != nil {
			log.Printf("SayHelloHandler: SayHello request unmarshal failed: %v", err)
			immediateError = &nrpc.Error{
				Type: nrpc.Error_CLIENT,
				Message: "bad request received: " + err.Error(),
			}
			serverRequestsForProvisioningService.WithLabelValues(
				"SayHello", request.Encoding, "unmarshal_fail").Inc()
		} else {
			request.Handler = func(ctx context.Context)(proto.Message, error){
				innerResp, err := h.server.SayHello(ctx, &req)
				if err != nil {
					return nil, err
				}
				return innerResp, err
			}
		}
	default:
		log.Printf("ProvisioningServiceHandler: unknown name %q", name)
		immediateError = &nrpc.Error{
			Type: nrpc.Error_CLIENT,
			Message: "unknown name: " + name,
		}
		serverRequestsForProvisioningService.WithLabelValues(
			"ProvisioningService", request.Encoding, "name_fail").Inc()
	}
	request.AfterReply = func(request *nrpc.Request, success, replySuccess bool) {
		if !replySuccess {
			serverRequestsForProvisioningService.WithLabelValues(
				request.MethodName, request.Encoding, "sendreply_fail").Inc()
		}
		if success {
			serverRequestsForProvisioningService.WithLabelValues(
				request.MethodName, request.Encoding, "success").Inc()
		} else {
			serverRequestsForProvisioningService.WithLabelValues(
				request.MethodName, request.Encoding, "handler_fail").Inc()
		}
		// report metric to Prometheus
		serverHETForProvisioningService.WithLabelValues(request.MethodName).Observe(
			request.Elapsed().Seconds())
	}
	if immediateError == nil {
		if h.workers != nil {
			// Try queuing the request
			if err := h.workers.QueueRequest(request); err != nil {
				log.Printf("nrpc: Error queuing the request: %s", err)
			}
		} else {
			// Run the handler synchronously
			request.RunAndReply()
		}
	}

	if immediateError != nil {
		if err := request.SendReply(nil, immediateError); err != nil {
			log.Printf("ProvisioningServiceHandler: ProvisioningService handler failed to publish the response: %s", err)
			serverRequestsForProvisioningService.WithLabelValues(
				request.MethodName, request.Encoding, "handler_fail").Inc()
		}
		serverHETForProvisioningService.WithLabelValues(request.MethodName).Observe(
			request.Elapsed().Seconds())
	} else {
	}
}

type ProvisioningServiceNRPCClient struct {
	nc      nrpc.NatsConn
	Subject string
	Encoding string
	Timeout time.Duration
}

func NewProvisioningServiceNRPCClient(nc nrpc.NatsConn) *ProvisioningServiceNRPCClient {
	return &ProvisioningServiceNRPCClient{
		nc:      nc,
		Subject: "ProvisioningService",
		Encoding: "protobuf",
		Timeout: 5 * time.Second,
	}
}

func (c *ProvisioningServiceNRPCClient) SayHello(req *SayHelloRequest) (*SayHelloResponse, error) {
	start := time.Now()

	subject := c.Subject + "." + "SayHello"

	// call
	var resp = SayHelloResponse{}
	if err := nrpc.Call(req, &resp, c.nc, subject, c.Encoding, c.Timeout); err != nil {
		clientCallsForProvisioningService.WithLabelValues(
			"SayHello", c.Encoding, "call_fail").Inc()
		return nil, err
	}

	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForProvisioningService.WithLabelValues("SayHello").Observe(elapsed)
	clientCallsForProvisioningService.WithLabelValues(
		"SayHello", c.Encoding, "success").Inc()

	return &resp, nil
}

func (c *ProvisioningServiceNRPCClient) SayHelloPoll(req *SayHelloRequest,maxreplies int, cb func (*SayHelloResponse) error,
) (error) {
	start := time.Now()

	subject := c.Subject + "." + "SayHello"

	var resp SayHelloResponse

	err := nrpc.Poll(req, &resp, c.nc, subject, c.Encoding, c.Timeout, maxreplies,
		func() error {
			return cb(&resp)
		},
	)
	if err != nil {
		clientCallsForProvisioningService.WithLabelValues(
			"SayHello", c.Encoding, "poll_fail").Inc()
		return err
	}

	// report total time taken to Prometheus
	elapsed := time.Since(start).Seconds()
	clientRCTForProvisioningService.WithLabelValues("SayHello").Observe(elapsed)
	clientCallsForProvisioningService.WithLabelValues(
		"SayHello", c.Encoding, "poll_success").Inc()

	return nil
}

type NRPCClient struct {
	nc      nrpc.NatsConn
	defaultEncoding string
	defaultTimeout time.Duration
	ProvisioningService *ProvisioningServiceNRPCClient
}

func NewNRPCClient(nc nrpc.NatsConn) *NRPCClient {
	c := NRPCClient{
		nc: nc,
		defaultEncoding: "protobuf",
		defaultTimeout: 5*time.Second,
	}
	c.ProvisioningService = NewProvisioningServiceNRPCClient(nc)
	return &c
}

func (c *NRPCClient) SetEncoding(encoding string) {
	c.defaultEncoding = encoding
	if c.ProvisioningService != nil {
		c.ProvisioningService.Encoding = encoding
	}
}

func (c *NRPCClient) SetTimeout(t time.Duration) {
	c.defaultTimeout = t
	if c.ProvisioningService != nil {
		c.ProvisioningService.Timeout = t
	}
}

func init() {
	// register metrics for service ProvisioningService
	prometheus.MustRegister(clientRCTForProvisioningService)
	prometheus.MustRegister(serverHETForProvisioningService)
	prometheus.MustRegister(clientCallsForProvisioningService)
	prometheus.MustRegister(serverRequestsForProvisioningService)
}