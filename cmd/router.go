package main

import (
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mangirdaz/image-inspector/pkg/config"
	"github.com/mangirdaz/image-inspector/pkg/storage"
)

func NewRouter() {

	//create new router
	router := mux.NewRouter().StrictSlash(false)
	storage := storage.InitKVStorage()

	//api backend init
	router.Path("/healthz").Name("Health endpoint").HandlerFunc(http.HandlerFunc(mybackendHandler(Health, storage)))
	//api backend init
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	apiV1.Methods("GET").Path("").Name("Index").Handler(http.HandlerFunc(mybackendHandler(Index, storage)))

	//notes methods
	notes := apiV1.PathPrefix("/image").Subrouter()

	notes.Methods("POST").Name("Init Scan").Handler(http.HandlerFunc(mybackendHandler(InitScan, storage)))

	//external resources
	static := apiV1.PathPrefix("/status").Subrouter()
	static.Methods("GET").Name("get images  status").Handler(http.HandlerFunc(mybackendHandler(Index, storage)))
	static.Methods("GET").Path("/{image}").Name("get particular image status").Handler(http.HandlerFunc(mybackendHandler(Index, storage)))

	//middleware intercept
	midd := http.NewServeMux()
	midd.Handle("/", router)
	midd.Handle("/api/v1", negroni.New(
		negroni.HandlerFunc(CorsHeadersMiddleware),
		negroni.Wrap(apiV1),
	))
	n := negroni.Classic()
	n.UseHandler(midd)
	url := fmt.Sprintf("%s:%s", config.Get("EnvAPIIP"), config.Get("EnvAPIPort"))

	log.WithFields(log.Fields{
		"url": url,
	}).Debug("api: starting api server")

	log.Fatal(http.ListenAndServe(url, n))

}
