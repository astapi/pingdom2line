package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port is server port to be listened.
	Port        string `envconfig:"PORT" default:"3000"`
	AccessToken string `envconfig:"LINE_NOTIFY_TOKEN"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

var env envConfig

func _main(args []string) int {
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	http.HandleFunc("/ping", ping)
	http.Handle("/pingdomhook", LineNotifyHandler{
		AccessToken: env.AccessToken,
	})

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("pong")
	w.WriteHeader(http.StatusOK)
}
