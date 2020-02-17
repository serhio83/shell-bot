package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/serhio83/shell-bot/pkg/structs"
	"github.com/serhio83/shell-bot/pkg/utils"
)

// returns a simple HTTP handler function which writes a response.
func mainpage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rHost := strings.Split(r.Host, ":")[0]
		rAddr := strings.Split(r.RemoteAddr, ":")[0]

		// check if Content-Type: application/json used and fail if not
		ah := r.Header.Get("Content-Type")
		if ah != "application/json" {
			log.Println(utils.StringDecorator(
				fmt.Sprintf("%s - %s [400] you should use Content-Type: application/json",
					rAddr,
					rHost)))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// fail if zero content length
		if r.ContentLength == 0 {
			log.Println(utils.StringDecorator(
				fmt.Sprintf("%s - %s [400] Invalid request payload",
					rAddr,
					rHost)))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// decode payload
		var msg structs.Message
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&msg); err != nil {
			log.Println(utils.StringDecorator(
				fmt.Sprintf("%s - %s [400] can`t decode json payload: %v",
					rAddr,
					rHost,
					err)))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// run remote ssh command if Alerts exists
		if len(msg.Alerts) > 0 {
			instanceHost := strings.Split(msg.Alerts[0].Labels.Instance, "://")[1]
			if len(instanceHost) > 0 {
				var stderr bytes.Buffer
				cmd := exec.Command("ssh", "-o StrictHostKeyChecking=no", instanceHost, "nginx -t && nginx -s reload")
				cmd.Stderr = &stderr
				errr := cmd.Run()
				if errr != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					log.Println(utils.StringDecorator(
						fmt.Sprintf("%s - %s [400] exec.Command failed: %v",
							rAddr,
							rHost,
							utils.StringSplitter(stderr))))
					return
				}

				// parse command execution result, split new line and write to log
				log.Println(
					utils.StringDecorator(
						fmt.Sprintf("%s - %s [200] instance: %s, alertname: %s, UA: %s",
							rAddr,
							rHost,
							msg.Alerts[0].Labels.Instance,
							msg.Alerts[0].Labels.Alertname,
							r.Header.Get("User-Agent"))))

				log.Println(utils.StringDecorator("[ssh.exec] " + utils.StringSplitter(stderr)))

				// give responce to client
				w.WriteHeader(200)
			}
		} else {

			// fail when no Alerts in payload
			log.Println(utils.StringDecorator(
				fmt.Sprintf("%s - %s [400] bad json payload",
					rAddr,
					rHost)))
			http.Error(w, "bad json payload", http.StatusBadRequest)
			return
		}

	}
}
