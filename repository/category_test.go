package repository

import (
	"context"
	"testing"

	"github.com/suared/core-apiuser/model"

	_ "github.com/suared/core/infra"
	"github.com/suared/core/security"
)

//CRUD is only real test here with 1 model move as there are no pointers that are re-used yet...
func TestSimpleCategoryDynamoLifeCycle(t *testing.T) {

	ctx := context.TODO()
	ctx = security.SetupTestAuthFromContext(ctx, 1)

	repository, err := NewCategoryRepository()
	if err != nil {
		t.Errorf("Repo initialization failed with: %v", err)
	}

	root := NewCategoryUserModel("Testing")
	personalCat := model.NewCategory("Personal")
	root.AddChild(*personalCat)
	workCat := model.NewCategory("Work")
	root.AddChild(*workCat)
	busCat := model.NewCategory("Entrepeneurial")
	root.AddChild(*busCat)

	//Create/ Read
	repository.Insert(ctx, *root)
	queryModel := CategoryUserModel{}
	queryModel.ID = root.ID
	dbroot, err := repository.SelectOne(ctx, queryModel)
	if err != nil {
		t.Errorf("During select, received error: %v", err)
	}
	if len(dbroot.GetAllChildren()) != 3 {
		t.Errorf("Expected 3 children, received: %v.  Full object is: %v", len(dbroot.GetAllChildren()), dbroot)
	}
	dbFoundChild, dbFoundParent := dbroot.FindChildByID(personalCat.ID)
	if dbFoundChild.ID != personalCat.ID {
		t.Errorf("Unexpected child, received: %v, expected: %v", dbFoundChild, personalCat)
	}
	if dbFoundParent != nil {
		t.Errorf("Exected parent to be a root, instead received: %v", dbFoundParent)
	}
	//save as a set ID for future comparison
	dbroot.ID = "A"
	err = repository.Insert(ctx, dbroot)
	if err != nil {
		t.Errorf("Unexpected error saving first test scenario: %v", err)
	}

	//Update / Read

	personalCat = root.GetChildByName("Personal")
	personalCat.AddChild(*getDisconnectedCategorySet())
	err = repository.Update(ctx, *root)
	if err != nil {
		t.Errorf("error during update: %v", err)
	}
	dbroot, err = repository.SelectOne(ctx, queryModel)
	if err != nil {
		t.Errorf("error during after update: %v", err)
	}
	dbPlay, _ := dbroot.FindChildByName("Play")
	if len(dbPlay.Children) != 3 {
		t.Errorf("Expected 3 children under play, received: %v, of item: %v", len(dbPlay.Children), dbPlay)
	}
	dbPersonal := dbroot.GetChildByName("Personal")
	if dbPersonal.Children[0].Title != "Play" {
		t.Errorf("Expected Play to be child of Personal, Personal has: %v", dbPersonal)
	}

	//save current progress of testing
	dbroot.ID = "B"
	err = repository.Insert(ctx, dbroot)
	if err != nil {
		t.Errorf("Unexpected error saving first test scenario: %v", err)
	}

	//test move
	dbroot, err = repository.SelectOne(ctx, queryModel)
	if err != nil {
		t.Errorf("error during after update: %v", err)
	}

	dbVideo, _ := dbroot.FindChildByName("Videos")
	dbWork, _ := dbroot.FindChildByName("Work")
	dbroot.Move(dbVideo.ID, dbWork)

	err = repository.Update(ctx, dbroot)
	if err != nil {
		t.Errorf("after move update error: %v", err)
	}

	dbroot, err = repository.SelectOne(ctx, queryModel)
	if err != nil {
		t.Errorf("error in retrieve: %v", err)
	}

	dbVideo = dbroot.GetChildByName("Videos")
	dbWork = dbroot.GetChildByName("Work")

	dbVideo, dbVideoParent := dbroot.FindChildByName("Videos")
	if len(dbVideo.Children) != 2 {
		t.Errorf("Expected 2 children under Videos, received: %v, of item: %v", len(dbVideo.Children), dbVideo)
	}

	if dbWork.Children[0].Title != "Videos" {
		t.Errorf("Expected Videos to be child of Work, Work has: %v", dbWork)
	}

	if dbVideoParent.Title != "Work" {
		t.Errorf("Failed move to Work area, parent is: %v, expected Work", dbVideoParent)
	}
	//save current progress of testing
	dbroot.ID = "C"
	err = repository.Insert(ctx, dbroot)
	if err != nil {
		t.Errorf("Unexpected error saving first test scenario: %v", err)
	}

	//Delete / Read
	err = repository.Delete(ctx, queryModel)
	if err != nil {
		t.Errorf("Delete failed, received: %v", err)
	}

	//Dynamo has latency on these types of things but give this a shot given testing is local/ 1 instance
	dbroot, err = repository.SelectOne(ctx, queryModel)
	if err != nil {
		t.Errorf("error in delete retrieve: %v", err)
	}
	if dbroot.ID != "" {
		t.Errorf("Expected Empty Struct after delete!, instead received: %v", dbroot)
	}

}

//Leveraging this start from model test
func getDisconnectedCategorySet() *model.Category {
	/*
		Play
			- Games
			- Videos
				-Beetlejuice
				-Something About Mary
			- Music
	*/
	playCat := model.NewCategory("Play")
	playCat.AddChild(*model.NewCategory("Games"))
	videoCat := model.NewCategory("Videos")
	videoCat = playCat.AddChild(*videoCat)
	videoCat.AddChild(*model.NewCategory("Beetlejuice"))
	videoCat.AddChild(*model.NewCategory("Something About Mary"))
	playCat.AddChild(*model.NewCategory("Music"))
	return playCat
}
