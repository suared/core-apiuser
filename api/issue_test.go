package api

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	coretest "github.com/suared/core/test"
)

// From here is the API to test

func TestAppLifeCycleE2E(t *testing.T) {
	//Setting as real integration style client vs. unit test style in mux example
	//Intentionally skipping event catching and testing here, will use request bin for now
	appURIBase := os.Getenv("PROCESS_LISTEN_URI")
	appURI := appURIBase + os.Getenv("PROCESS_RELATIVE_PATH") + "/app"
	issueURI := appURI + "/issue"
	//monitoringURI := appURI + "/monitoring"

	log.Printf("Running API Lifecyce Test E2E with uri: %v", appURI)

	//Case: XXX
	issue := Issue{
		UserText: "My Login is not working for SAP",
	}
	byteArr, err := json.Marshal(issue)
	resourceURI, err := coretest.SimplePost(issueURI, byteArr)
	if err != nil {
		t.Errorf("Expected Successful add, received error: %v, response: %v", err, resourceURI)
	}

}
