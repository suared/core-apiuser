package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	coreapi "github.com/suared/core/api"
	coreerrors "github.com/suared/core/errors"
)

var relPathApp string

func init() {
	// local shared variables
	//Will be /api/domain... or /api/version/domain (must start with a /)..
	relPathApp = os.Getenv("PROCESS_RELATIVE_PATH") + "/app"
	log.Printf("App service started with relPathApp: %v", relPathApp)

}

// SetupAppRoutes - Setup Routes for the App API
func SetupAppRoutes(router *mux.Router) {
	//Syntax:
	//router.HandleFunc("pathuri"), methodname).Methods("HTTPVERB")
	//{variableName}
	//Example:
	//router.HandleFunc("/api/tasks"), processSummary).Methods("GET") (or "POST", etc...)

	/* TODO: Example
	router.HandleFunc(relPathApp+"/issue", newIssue).Methods("POST")
	*/
}

//API Request object(s)
/*TODO: Example

type Issue struct {
	Application string
	Status      string
	Details     string
	ErrNum      string
}
*/

//POST /conversation
func newIssue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//do work here - nil is error object if exists - first nil below should be the return object, second is the error
	byteMessage, _ := ioutil.ReadAll(r.Body)

	issue := &Issue{}
	err := json.Unmarshal(byteMessage, issue)
	//perform validations
	if err != nil {
		apiErr := coreerrors.NewClientError("Conversation structure not met: " + err.Error())
		coreapi.WritePostAPIResponse(ctx, w, r, "", apiErr)
		return
	}

	//handle request (some service call out in real life)
	convoLocation, err := handleNewIssue(issue)
	if err != nil {
		apiErr := coreerrors.NewClientError("Conversation could not be handled: " + err.Error())
		coreapi.WritePostAPIResponse(ctx, w, r, "", apiErr)
		return
	}
	coreapi.WritePostAPIResponse(ctx, w, r, convoLocation, nil)
}

func getIssueError(r *http.Request, when string, err error) error {
	//To start will always assume a system error and build up library of user errors over time
	//If known up front, the appropriate message will be created by api code, if common will add translator here
	_, ok := err.(coreerrors.Error)
	if ok {
		return err
	}
	apiError := coreerrors.NewError(err)
	return apiError
}
