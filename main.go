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
	"sync"
)

var (
	port           = flag.Int("p", 4007, "Hook listener port number")
	token          = flag.String("t", "", "Secret token")
	configFilename = flag.String("c", "bitbuckethook.json", "Hook listener config")
)

type payloadHandler struct {
	config hookConfig
	mutex  sync.Mutex
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

	go this.processPayload(payload)
}

func (this *payloadHandler) processPayload(payload Payload) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	repoName := payload.Repository.Name
	command, ok := this.config[repoName]
	if ok {
		commandString := strings.Join(command, " ")
		log.Println(repoName, ":", commandString)
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println(repoName, ":", "handler not found")
	}
}

func main() {
	flag.Parse()
	config, err := parseConfig(*configFilename)
	if err != nil {
		log.Fatal(err)
	}

	handler := payloadHandler{config: config}

	addr := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(addr, &handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
