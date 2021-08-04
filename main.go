package main

import (
	"github.com/negasus/haproxy-spoe-go/action"
	"github.com/negasus/haproxy-spoe-go/agent"
	"github.com/negasus/haproxy-spoe-go/request"
	"github.com/negasus/haproxy-spoe-go/message"
  "github.com/negasus/haproxy-spoe-go/varint"
	"log"
	"math/rand"
	"net"
	"os"
)

func extractHdrs(prefix string, buf []byte) {
  for {
    keyLen, i := varint.Uvarint(buf)
    buf = buf[i:]
    key := string(buf[:keyLen])
    buf = buf[keyLen:]
    valLen, i := varint.Uvarint(buf)
    buf = buf[i:]
    val := string(buf[:valLen])
    buf = buf[valLen:]
    if keyLen == 0 && valLen == 0 {
      return
    }
    log.Printf("cap-hdr-%s %s: %s\n", prefix, key, val)
  }
}

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

  hdrsValue, ok := msg.KV.Get("hdrs")
  if !ok {
		log.Printf("var 'hdrs' not found in message")
		return
  }
  hdrs, ok := hdrsValue.([]byte)
	if !ok {
		log.Printf("var 'hdrs' has wrong type. expect IP addr")
		return
	}
  extractHdrs("request", hdrs)

	ipScore := rand.Intn(100)

	log.Printf("IP: %s, send score '%d'", ip.String(), ipScore)

	req.Actions.SetVar(action.ScopeSession, "ip_score", ipScore)
	req.Actions.SetVar(action.ScopeSession, "trace_context", "abc-123")
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
  log.Printf("body: %s\n", string(body))

  hdrsValue, ok := msg.KV.Get("hdrs")
  if !ok {
		log.Printf("var 'hdrs' not found in message")
		return
  }
  hdrs, ok := hdrsValue.([]byte)
	if !ok {
		log.Printf("var 'hdrs' has wrong type. expect IP addr")
		return
	}
  extractHdrs("response", hdrs)


  tracectxValue, ok := msg.KV.Get("trace_context")
  if !ok {
		log.Printf("var 'trace_context' not found in message")
		return
  }
  tracectx, ok := tracectxValue.(string)
	if !ok {
		log.Printf("var 'trace_context' has wrong type. expect IP addr")
		return
	}
  log.Printf("trace-context: %s", tracectx)
}
