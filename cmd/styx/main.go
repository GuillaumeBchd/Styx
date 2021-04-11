package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/guillaumebchd/styx/pkg/conf"
	"github.com/guillaumebchd/styx/pkg/ddos"
	"github.com/guillaumebchd/styx/pkg/rvp"
	"github.com/pelletier/go-toml"
)

func main() {

	// We read our configuration file
	configuration, err := toml.LoadFile("./configuration.toml")
	if err != nil {
		log.Fatal(err)
	}
	conf := conf.Get(configuration)

	// We create our router
	r := mux.NewRouter()

	// We create our reverse proxy from our configuration object object
	reverseProxy := rvp.GenerateProxy(conf.ReverseProxy)

	// We capture all the paths and we redirect it to the reverse proxy
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)
	})

	ddosProtect := ddos.New("GlobalLimit", 1, 2)

	r.Use(ddosProtect.Proctection)

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("0.0.0.0:%d", conf.Server.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting server " + conf.Server.ServerName + " on : " + srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
