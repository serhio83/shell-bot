package handlers

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/serhio83/shell-bot/pkg/utils"
)

//Router return new mux.Router
func Router(buildTime, commit, release string) *mux.Router {
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Println(utils.StringDecorator("Readyz probe is negative by default..."))
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		log.Println(utils.StringDecorator("Readyz probe is positive."))
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", checkHeaders(mainpage))
	r.HandleFunc("/home", home(buildTime, commit, release)).Methods("GET")
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/readyz", readyz(isReady))
	return r
}
