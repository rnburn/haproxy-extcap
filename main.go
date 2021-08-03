package main

import (
	"github.com/negasus/haproxy-spoe-go/action"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/negasus/haproxy-spoe-go/request"
	"github.com/negasus/haproxy-spoe-go/message"
	"log"
	"math/rand"
	"net"
	"os"
)

func main() {

	log.Print("listen 9000")

	listener, err := net.Listen("tcp4", ":9000")
	if err != nil {
		log.Printf("error create listener, %v", err)
		os.Exit(1)
	}
	defer listener.Close()

	a := agent.New(handler)

	if err := a.Serve(listener); err != nil {
		log.Printf("error agent serve: %+v\n", err)
	}
}

func handler(req *request.Request) {

	log.Printf("handle request EngineID: '%s', StreamID: '%d', FrameID: '%d' with %d messages\n", req.EngineID, req.StreamID, req.FrameID, req.Messages.Len())

  for _, msg := range *req.Messages {
    switch msg.Name {
    case "check-client-ip":
      handleRequestMessage(req, msg)
    case "extcap-response":
      handleResponseMessage(req, msg)
    default:
      log.Printf("unkown message %s", msg.Name)
    }
  }
}

func handleRequestMessage(req *request.Request, msg *message.Message) {
	ipValue, ok := msg.KV.Get("ip")
	if !ok {
		log.Printf("var 'ip' not found in message")
		return
	}

	ip, ok := ipValue.(net.IP)
	if !ok {
		log.Printf("var 'ip' has wrong type. expect IP addr")
		return
	}

  bodyValue, ok := msg.KV.Get("body")
  if !ok {
		log.Printf("var 'body' not found in message")
		return
  }
  body, ok := bodyValue.([]byte)
	if !ok {
		log.Printf("var 'body' has wrong type. expect IP addr")
		return
	}
  log.Printf("request body length %d\n", len(body))

	ipScore := rand.Intn(100)

	log.Printf("IP: %s, send score '%d'", ip.String(), ipScore)

	req.Actions.SetVar(action.ScopeSession, "ip_score", ipScore)
}

func handleResponseMessage(req *request.Request, msg *message.Message) {
  bodyValue, ok := msg.KV.Get("body")
  if !ok {
		log.Printf("var 'body' not found in message")
		return
  }
  body, ok := bodyValue.([]byte)
	if !ok {
		log.Printf("var 'body' has wrong type. expect IP addr")
		return
	}
  log.Printf("response body length %d\n", len(body))
}
