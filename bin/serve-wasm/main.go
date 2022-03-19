package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func compile(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	cmd := exec.Command("make", "build-wasm")
	if output, err := cmd.CombinedOutput(); err != nil {
		handleErr(w, fmt.Errorf("%s: %s", err, string(output)))
		return
	}
	wasm, err := ioutil.ReadFile("webhookdb.wasm")
	if err != nil {
		handleErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/wasm")
	w.WriteHeader(200)
	_, _ = w.Write(wasm)
}

func handleErr(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(fmt.Sprintf("Error: %s", err)))
}

func main() {
	http.HandleFunc("/compile", compile)
	if err := http.ListenAndServe(":18008", nil); err != nil {
		log.Fatal(err)
	}
}
