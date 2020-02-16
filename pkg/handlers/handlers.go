package handlers

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	utils "github.com/serhio83/shell-bot/pkg/utils"
)

//Router ...
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
	r.HandleFunc("/", mainpage())
	r.HandleFunc("/home", home(buildTime, commit, release)).Methods("GET")
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/readyz", readyz(isReady))
	return r
}
