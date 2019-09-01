package api

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	coreapi "github.com/suared/core/api"
)

//Using this as more of the integration test across the task API.  Marking with "E2E" to enable only thosoe tests to run separately in future outside box testing
//Should catch most errors at minimal cost since all local while allowing reuse for E2E
//Note that I am killing the server in this same run so if there are multiple tests, all other tests should happen first

//This is just a copy of the main calling the same initial setup for process
type apiRoutes struct {
	startServer bool
}

func (routes *apiRoutes) SetupRoutes(router *mux.Router) {
	SetupProcessRoutes(router)
	SetupTaskRoutes(router)
}

func (routes *apiRoutes) StartServer() bool {
	//tester will always run in integration mode for now
	return routes.startServer
}

func init() {
	//start listener
	log.Println("Init called on task api integration test")

}

func TestMain(m *testing.M) {
	log.Printf("loading test with e2e set as: %v", os.Getenv("LAMBDA_ENV"))
	//Going to use Lambda env synonumously w. e2e only for now since both stop the loading of the server
	if os.Getenv("LAMBDA_ENV") == "true" {
		//config := &apiRoutes{startServer: false}
	} else {
		config := &apiRoutes{startServer: true}
		go coreapi.StartHTTPListener(config)
		// TODO: listen to startup vs. assuming like below, this will work for now
		time.Sleep(1 * time.Second)

	}
	//
	os.Exit(m.Run())
}
