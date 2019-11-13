package main

import (
	"os"

	"github.com/gorilla/mux"
	coreapi "github.com/suared/core/api"
	_ "github.com/suared/core/infra"

	"github.com/suared/core-apiuser/api"

	"github.com/akrylysov/algnhsa"
)

func main() {
	if os.Getenv("LAMBDA_ENV") == "true" {
		startLambdaAPI()
	} else {
		startWebAPI()
	}
}

type apiRoutes struct {
	refRouter *mux.Router
	autoStart bool
}

func (routes *apiRoutes) SetupRoutes(router *mux.Router) {
	api.SetupAppRoutes(router)
	api.SetupAppObjectiveRoutes(router)
	routes.refRouter = router

}

func (routes *apiRoutes) StartServer() bool {
	return routes.autoStart
}

func startWebAPI() {
	config := &apiRoutes{}
	config.autoStart = true
	coreapi.StartHTTPListener(config)
}

func startLambdaAPI() {
	config := &apiRoutes{}
	config.autoStart = false
	coreapi.StartHTTPListener(config)
	//In the lambda environment, this will strip out the earlier proxy portions
	opts := &algnhsa.Options{UseProxyPath: true}
	//Since autoStart is false, it is up to the caller to start the engine to be used
	//log.Printf("at this point, config has: %v, refRouter has: %v", config, config.refRouter)
	algnhsa.ListenAndServe(config.refRouter, opts)
	/*
		Note from this library docs if I add binary content in the future (likel):
		To make the API Gateway treat certain content types as binary, you need to add the desired types to your API's "Binary Media Types" (Settings section) and also pass them to the algnhsa.ListenAndServe function:

		algnhsa.ListenAndServe(http.DefaultServeMux, []string{"image/jpeg", "image/png"}])
		You can find the algnhsa package (Lambda Go net/http server adapter) on GitHub - https://github.com/akrylysov/algnhsa.

		"github.com/akrylysov/algnhsa"
	*/
}
