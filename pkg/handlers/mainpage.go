package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	structs "github.com/serhio83/shell-bot/pkg/structs"
	utils "github.com/serhio83/shell-bot/pkg/utils"
)

// returns a simple HTTP handler function which writes a response.
func mainpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check if Content-Type: application/json used and fail if not
		ah := r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if ah != "application/json" {
			log.Println(utils.StringDecorator("[error] you should use Content-Type: application/json"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// fail if zero content length
		if r.ContentLength == 0 {
			log.Println(utils.StringDecorator("[error] Invalid request payload"))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// decode payload
		var msg structs.Message
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&msg); err != nil {
			log.Println(utils.StringDecorator(fmt.Sprintf("[error] can`t decode json payload: %v", err)))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// run remote ssh command if Alerts exists
		if len(msg.Alerts) > 0 {

			var stderr bytes.Buffer
			cmd := exec.Command("ssh", "-o StrictHostKeyChecking=no", "gw.tp.fbs", "nginx -t && nginx -s reload")
			cmd.Stderr = &stderr
			errr := cmd.Run()
			if errr != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				log.Println(utils.StringDecorator("[error] exec.Command failed: " + utils.StringSplitter(stderr)))
				return
			}

			// parse command execution result, split new line and write to log
			log.Println(
				utils.StringDecorator(
					fmt.Sprintf("[ok] instance: %s, alertname: %s, UA: %s",
						msg.Alerts[0].Labels.Instance,
						msg.Alerts[0].Labels.Alertname,
						r.Header.Get("User-Agent"))))

			log.Println(utils.StringDecorator("[ok] " + utils.StringSplitter(stderr)))

			// give responce to client
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)

		} else {

			// fail when no Alerts in payload
			log.Println(utils.StringDecorator("[error] bad json payload"))
			http.Error(w, "bad json payload", http.StatusBadRequest)
			return
		}

	}
}
