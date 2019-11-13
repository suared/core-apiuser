package service

import (
	"context"
	"testing"

	"github.com/suared/core-apiuser/model"

	_ "github.com/suared/core/infra"
	"github.com/suared/core/security"
)

// only repository test and 1 move, no shared pointers yet
func TestSimpleCategoryServiceLifeCycle(t *testing.T) {

	ctx := context.TODO()
	ctx = security.SetupTestAuthFromContext(ctx, 1)

	svc := NewCategoryService()

	//Get model will create the default if does not exist for the lifeApp id
	catModel, err := svc.GetCategoryModel(ctx, MyLifeCategoryUserModelID)

	if err != nil {
		t.Errorf("Get Model failed for Lifeapp id, err: %v", err)
	}

	if catModel.Name != "Life Categories" {
		t.Errorf("Default Lifeapp name not generated")
	}

	if len(catModel.Children) != 2 {
		t.Errorf("Expected 2 children in default lifeapp")
	}

	//tests to do...  replace update, move, add, delete
	childCat := model.NewCategory("Personal")
	childCat = childCat.AddChild(*getDisconnectedCategorySet())
	catModel.AddChild(*childCat)
	childCat = model.NewCategory("Work")
	catModel.AddChild(*childCat)
	childCat = model.NewCategory("Entrepeneurial")
	catModel.AddChild(*childCat)

	err = svc.ReplaceCategoryModel(ctx, catModel)
	if err != nil {
		t.Errorf("Replace failed with err: %v", err)
	}

	//Ensure using the updated version from the service
	catModel, err = svc.GetCategoryModel(ctx, MyLifeCategoryUserModelID)

	if err != nil {
		t.Errorf("Get Model failed for Lifeapp id, err: %v", err)
	}

	videos, _ := catModel.FindChildByName("Videos")
	work, _ := catModel.FindChildByName("Work")

	err = svc.MoveCategory(ctx, catModel.ID, work.ID, videos.ID)

	//Refresh and see if we get the latest
	catModel, err = svc.GetCategoryModel(ctx, MyLifeCategoryUserModelID)

	if err != nil {
		t.Errorf("Get Model failed for Lifeapp id, err: %v", err)
	}

	videos, videoParent := catModel.FindChildByName("Videos")
	work, _ = catModel.FindChildByName("Work")

	if len(videos.Children) != 2 {
		t.Errorf("Expected 2 children under Videos, received: %v, of item: %v", len(videos.Children), videos)
	}

	if work.Children[0].Title != "Videos" {
		t.Errorf("Expected Videos to be child of Work, Work has: %v", work)
	}

	if videoParent.Title != "Work" {
		t.Errorf("Failed move to Work area, parent is: %v, expected Work", videoParent)
	}

	//left to test:  update, add, delete
	err = svc.AddCategory(ctx, MyLifeCategoryUserModelID, work.ID, *model.NewCategory("School"))

	if err != nil {
		t.Errorf("Add failed with: %v", err)
	}

	catModel, err = svc.GetCategoryModel(ctx, MyLifeCategoryUserModelID)

	if err != nil {
		t.Errorf("Get Model failed for Lifeapp id, err: %v", err)
	}

	school, schoolParent := catModel.FindChildByName("School")
	if schoolParent.Title != "Work" {
		t.Error("Expected Parent to be Work")
	}
	if school.ID == "" {
		t.Error("Expected school ID")
	}

	//update test
	school.Title = "Computer Science"
	oldSchoolID := school.ID
	err = svc.UpdateCategory(ctx, MyLifeCategoryUserModelID, *school)

	catModel, err = svc.GetCategoryModel(ctx, MyLifeCategoryUserModelID)

	if err != nil {
		t.Errorf("Get Model failed for Lifeapp id, err: %v", err)
	}

	school, schoolParent = catModel.FindChildByName("Computer Science")
	if schoolParent.Title != "Work" {
		t.Error("Expected Parent to be Work")
	}
	if school.ID != oldSchoolID {
		t.Error("Expected school ID not to change during update")
	}

	//delete test
	err = svc.DeleteCategory(ctx, MyLifeCategoryUserModelID, school.ID)
	catModel, err = svc.GetCategoryModel(ctx, MyLifeCategoryUserModelID)

	if err != nil {
		t.Errorf("Get Model failed for Lifeapp id, err: %v", err)
	}

	school, _ = catModel.FindChildByName("Computer Science")
	if school.ID != "" {
		t.Error("Expected to be empty after delete")
	}

	//Cleanup...
	err = svc.DeleteCategoryModel(ctx, MyLifeCategoryUserModelID)
	if err != nil {
		t.Errorf("Delete User Model failed, err: %v", err)
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
