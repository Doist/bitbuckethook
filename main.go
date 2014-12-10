package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	port           = flag.Int("p", 4007, "Hook listener port number")
	qSize          = flag.Int("q", 10, "Request backlog length")
	token          = flag.String("t", "", "Secret token")
	configFilename = flag.String("c", "bitbuckethook.json", "Hook listener config")
)

func newPayloadHandler(config hookConfig, qSize int) *payloadHandler {
	if qSize < 0 {
		qSize = 0
	}
	return &payloadHandler{
		qSize:    qSize,
		config:   config,
		incoming: make(chan *Payload, 10),
		reqs:     make(map[string]chan *Payload),
	}
}

func (ph *payloadHandler) Loop() {
	for p := range ph.incoming {
		if args, ok := ph.config[p.Repository.Name]; !ok || len(args) == 0 {
			continue
		}
		if _, ok := ph.reqs[p.Repository.Name]; !ok {
			c := make(chan *Payload, ph.qSize)
			ph.reqs[p.Repository.Name] = c
			go payloadProcessor(c, ph.config[p.Repository.Name])
		}
		select {
		case ph.reqs[p.Repository.Name] <- p:
		default: // spillover
		}
	}
}

type payloadHandler struct {
	qSize    int
	config   hookConfig
	incoming chan *Payload
	reqs     map[string]chan *Payload
}

func (this *payloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if *token != "" {
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		getToken, ok := values["token"]
		if !ok || getToken[0] != *token {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := r.FormValue("payload")
	payload, err := parsePayload(body)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	this.incoming <- payload
}

func payloadProcessor(ch chan *Payload, args []string) {
	cmdString := strings.Join(args, " ")
	for p := range ch {
		log.Printf("%s: %s", p.Repository.Name, cmdString)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	config, err := parseConfig(*configFilename)
	if err != nil {
		log.Fatal(err)
	}
	if len(config) == 0 {
		log.Fatal("empty config")
	}

	handler := newPayloadHandler(config, *qSize)
	go handler.Loop()

	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, handler))
}
