package api

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/suared/core-apiuser/model"
	"github.com/suared/core-apiuser/repository"

	"github.com/suared/core/security"
	"github.com/suared/core/uuid"

	coretest "github.com/suared/core/test"
)

// From here is the API to test

func TestAppLifeCycleE2E(t *testing.T) {
	//integration test intentional here
	appURIBase := os.Getenv("PROCESS_LISTEN_URI")
	appURI := appURIBase + os.Getenv("PROCESS_RELATIVE_PATH")
	categoriesURI := appURI + "/categories"
	lifeAppCategoriesURI := categoriesURI + "/lifeapp"

	//Tests:  Get; Patch (update, move, add, delete)
	log.Printf("Running API Lifecyce Test E2E with uri: %v", appURI)

	//Get
	response, err := coretest.SimpleGet(lifeAppCategoriesURI)
	if err != nil {
		t.Errorf("Get failed with: %v", err)
	}
	catModel := repository.CategoryUserModel{}

	err = json.Unmarshal([]byte(response), &catModel)
	if err != nil {
		t.Errorf("Error in Unmarshal: %v", err)
	}

	if catModel.ID != myLifeCategoryUserModelID {
		t.Errorf("Default model not created")
	}
	//Add
	actions := CategoryActions{
		Operation: "ADD",
		ParentID:  "",
		ID:        uuid.NewUUID(),
		Title:     "Personal",
	}

	byteArr, err := json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	response, err = coretest.SimpleGet(lifeAppCategoriesURI)
	if err != nil {
		t.Errorf("Get failed with: %v", err)
	}

	catModel = repository.CategoryUserModel{}

	err = json.Unmarshal([]byte(response), &catModel)
	if err != nil {
		t.Errorf("Error in Unmarshal: %v", err)
	}

	personalCat, personalParent := catModel.FindChildByName("Personal")

	if personalCat.Title != "Personal" {
		t.Error("Expected Personal to exist")
	}

	if personalParent != nil {
		t.Error("Expected parent to be root")
	}

	actions = CategoryActions{
		Operation: "ADD",
		ParentID:  "",
		ID:        uuid.NewUUID(),
		Title:     "Work",
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	gamesCat := model.NewCategory("Games")
	actions = CategoryActions{
		Operation: "ADD",
		ParentID:  personalCat.ID,
		ID:        gamesCat.ID,
		Title:     gamesCat.Title,
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	actions = CategoryActions{
		Operation: "ADD",
		ParentID:  gamesCat.ID,
		ID:        uuid.NewUUID(),
		Title:     "Monopoly",
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	actions = CategoryActions{
		Operation: "ADD",
		ParentID:  gamesCat.ID,
		ID:        uuid.NewUUID(),
		Title:     "Parcheesi",
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	response, err = coretest.SimpleGet(lifeAppCategoriesURI)
	if err != nil {
		t.Errorf("Get failed with: %v", err)
	}

	catModel = repository.CategoryUserModel{}
	err = json.Unmarshal([]byte(response), &catModel)
	if err != nil {
		t.Errorf("Error in Unmarshal: %v", err)
	}

	//Move
	workCat, _ := catModel.FindChildByName("Work")

	actions = CategoryActions{
		Operation: "MOVE",
		ParentID:  workCat.ID,
		ID:        gamesCat.ID,
		Title:     "",
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	response, err = coretest.SimpleGet(lifeAppCategoriesURI)
	if err != nil {
		t.Errorf("Get failed with: %v", err)
	}

	catModel = repository.CategoryUserModel{}
	err = json.Unmarshal([]byte(response), &catModel)
	if err != nil {
		t.Errorf("Error in Unmarshal: %v", err)
	}

	gamesCat, gamesParent := catModel.FindChildByName("Games")
	if gamesParent.Title != "Work" {
		t.Error("Expected parent to be work after move")
	}
	if gamesCat.Title != "Games" {
		t.Error("Title not expected to be changed in move operation")
	}

	//Update
	parcheesiCat, _ := catModel.FindChildByName("Parcheesi")
	actions = CategoryActions{
		Operation: "UPDATE",
		ParentID:  "",
		ID:        parcheesiCat.ID,
		Title:     "Operation",
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	response, err = coretest.SimpleGet(lifeAppCategoriesURI)
	if err != nil {
		t.Errorf("Get failed with: %v", err)
	}

	catModel = repository.CategoryUserModel{}
	err = json.Unmarshal([]byte(response), &catModel)
	if err != nil {
		t.Errorf("Error in Unmarshal: %v", err)
	}

	operationCat, operationParent := catModel.FindChildByName("Operation")
	if operationParent.Title != "Games" {
		t.Error("Expected parent to be Games")
	}
	if operationCat.ID != parcheesiCat.ID {
		t.Error("model id should not change for update")
	}
	if len(operationParent.Children) != 2 {
		t.Error("Expected 2 children after update")
	}

	//Delete
	actions = CategoryActions{
		Operation: "DELETE",
		ParentID:  "",
		ID:        parcheesiCat.ID,
		Title:     "",
	}

	byteArr, err = json.Marshal(actions)
	if err != nil {
		t.Errorf("Could not marshall action object: %v", err)
	}
	coretest.SimplePatch(lifeAppCategoriesURI, byteArr)

	response, err = coretest.SimpleGet(lifeAppCategoriesURI)
	if err != nil {
		t.Errorf("Get failed with: %v", err)
	}

	catModel = repository.CategoryUserModel{}
	err = json.Unmarshal([]byte(response), &catModel)
	if err != nil {
		t.Errorf("Error in Unmarshal: %v", err)
	}

	operationCat, operationParent = catModel.FindChildByName("Operation")
	if operationCat.ID != "" {
		t.Error("Expected would not exist after DELETE")
	}

	gamesCat, _ = catModel.FindChildByName("Games")
	if len(gamesCat.Children) != 1 {
		t.Error("Expected only 1 child after delete")
	}

	//Cleanup
	ctx := context.TODO()
	ctx = security.SetupTestAuthFromContext(ctx, 1)

	err = categoryService.DeleteCategoryModel(ctx, myLifeCategoryUserModelID)
	if err != nil {
		t.Errorf("Delete User Model failed, err: %v", err)
	}

}
