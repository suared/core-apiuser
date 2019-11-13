package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	coreapi "github.com/suared/core/api"
	coreerrors "github.com/suared/core/errors"

	"github.com/suared/core-apiuser/model"
	"github.com/suared/core-apiuser/service"
)

var relPathCategory string

const myLifeCategoryUserModelID = service.MyLifeCategoryUserModelID

var categoryService *service.CategoryService

func init() {
	// local shared variables
	//Will be /api/domain... or /api/version/domain (must start with a /)..
	relPathCategory = os.Getenv("PROCESS_RELATIVE_PATH") + "/categories"
	log.Printf("App service started with relPathApp: %v", relPathCategory)
	categoryService = service.NewCategoryService()

}

// SetupAppRoutes - Setup Routes for the App API
func SetupAppRoutes(router *mux.Router) {
	//Syntax:
	//router.HandleFunc("pathuri"), methodname).Methods("HTTPVERB")
	//{variableName}
	//Example:
	//router.HandleFunc("/api/tasks"), processSummary).Methods("GET") (or "POST", etc...)

	/* This API has (Note: added named category id to enable to reuse for other purposes):
	* Get lifeapp life org category model - GET Lifeapp/Categories/myLife;  Returns CategoryUserModel
	* Add a category  -  PATCH Lifeapp/Categories/myLife   <Category Object w/  Action>; Returns Success/Failure
	* Delete a category - PATCH Lifeapp/Categories/myLife	<Category Object w/  Action>; Returns Success/Failure
	* Move a category - PATCH Lifeapp/Categories/myLife 	<Category Object w/  Action>; Returns Success/Failure
	 */

	urlToHandle := relPathCategory + "/lifeapp" //  -->  lifeApp/categories/lifeapp
	router.HandleFunc(urlToHandle, getLifeCategoryModel).Methods("GET")
	router.HandleFunc(relPathCategory+"/lifeappList", getLifeCategoryList).Methods("GET")
	router.HandleFunc(urlToHandle, patchCategoryLifeModel).Methods("PATCH")
}

//API Request object(s)

//CategoryActions - Defines the patch object expected when interacting with life app category actions
//Operation is required, one of:  ADD, MOVE, DELETE, UPDATE
//ParentID - Required for Add and Move.
//ID - Required for All Actions
//Title - Required for ADD and UPDATE
type CategoryActions struct {
	Operation string `json:"operation"` //Required for All Actions - Add, Move, Delete, Update
	ParentID  string `json:"parentID"`  //Add = parent ID, Move = New Parent ID
	ID        string `json:"id"`        //Required for All Actions
	Title     string `json:"title"`     //Required for Add, Update
}

//repository.CategoryUserModel is the other API object that will be used

//GET /Lifeapp/Categories/lifeapp
func getLifeCategoryModel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := categoryService.GetCategoryModel(ctx, myLifeCategoryUserModelID)
	//Assume all are system errors to start, will start converting to split out user errors later
	//getProcessAPIError internal call will be used to convert to user/ client errors in one central location as common erors are found
	if err != nil {
		apiErr := getCategoryError(r, "get", err)
		coreapi.WriteGetAPIResponse(ctx, w, r, categories, apiErr)
	} else {
		coreapi.WriteGetAPIResponse(ctx, w, r, categories, nil)
	}
}

//GET /Lifeapp/categories/lifeapplist
func getLifeCategoryList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := categoryService.GetCategoryModel(ctx, myLifeCategoryUserModelID)
	//lifeapp ui manages an array of categories for simplicity
	catList := categories.GetAllChildren()
	//Assume all are system errors to start, will start converting to split out user errors later
	//getProcessAPIError internal call will be used to convert to user/ client errors in one central location as common erors are found
	if err != nil {
		apiErr := getCategoryError(r, "get", err)
		coreapi.WriteGetAPIResponse(ctx, w, r, catList, apiErr)
	} else {
		coreapi.WriteGetAPIResponse(ctx, w, r, catList, nil)
	}
}

//PATCH /Lifeapp/Categories/lifeapp
func patchCategoryLifeModel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//First - Get the payload action object
	categoryAction := &CategoryActions{}
	byteMessage, _ := ioutil.ReadAll(r.Body)

	//log.Printf("TaskMessagePatch Body contains: %v", string(byteMessage))

	err := json.Unmarshal(byteMessage, categoryAction)

	if err != nil || categoryAction.ID == "" {
		apiErr := coreerrors.NewClientError("Body of message sent does not meet the category action structure (1):" + err.Error())
		coreapi.WritePatchAPIResponse(ctx, w, r, apiErr)
		return
	}

	//Finally, based on the operation, make the associated change to the process model
	if categoryAction.Operation == "UPDATE" {
		err = categoryService.UpdateCategory(ctx, myLifeCategoryUserModelID, model.Category{ID: categoryAction.ID, Title: categoryAction.Title})
		if err != nil {
			apiErr := coreerrors.NewClientError(fmt.Sprintf("Update failed during Category Patch Request: %v", err))
			coreapi.WritePatchAPIResponse(ctx, w, r, apiErr)
			return
		}
	} else if categoryAction.Operation == "MOVE" {
		err = categoryService.MoveCategory(ctx, myLifeCategoryUserModelID, categoryAction.ParentID, categoryAction.ID)
		if err != nil {
			apiErr := coreerrors.NewClientError(fmt.Sprintf("Update failed during Category Patch Delete Request: %v", err))
			coreapi.WritePatchAPIResponse(ctx, w, r, apiErr)
			return
		}
	} else if categoryAction.Operation == "ADD" {
		err = categoryService.AddCategory(ctx, myLifeCategoryUserModelID, categoryAction.ParentID, model.Category{ID: categoryAction.ID, Title: categoryAction.Title})
		if err != nil {
			apiErr := coreerrors.NewClientError(fmt.Sprintf("Update failed during Category Patch Delete Request: %v", err))
			coreapi.WritePatchAPIResponse(ctx, w, r, apiErr)
			return
		}
	} else if categoryAction.Operation == "DELETE" {
		err = categoryService.DeleteCategory(ctx, myLifeCategoryUserModelID, categoryAction.ID)
		if err != nil {
			apiErr := coreerrors.NewClientError(fmt.Sprintf("Update failed during Category Patch Delete Request: %v", err))
			coreapi.WritePatchAPIResponse(ctx, w, r, apiErr)
			return
		}
	} else {
		apiErr := coreerrors.NewClientError("No Matching Operation defined in the Category Patch Request")
		coreapi.WritePatchAPIResponse(ctx, w, r, apiErr)
		return
	}

}

func getCategoryError(r *http.Request, when string, err error) error {
	//To start will always assume a system error and build up library of user errors over time
	//If known up front, the appropriate message will be created by api code, if common will add translator here
	_, ok := err.(coreerrors.Error)
	if ok {
		return err
	}
	apiError := coreerrors.NewError(err)
	return apiError
}
