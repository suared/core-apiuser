package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/suared/core-apiuser/model"
	"github.com/suared/core-apiuser/repository"
)

var categoryRepo *repository.CategoryRepository

//MyLifeCategoryUserModelID - One time generated UUID will be used for the life app as there is only 1 per user
const MyLifeCategoryUserModelID = "1SBsF9WrcSmBwWvzWVojegYR6z2"

func init() {
	catRepo, err := repository.NewCategoryRepository()
	if err != nil {
		panic("Unable to setup Category Repository while initializing the category service")
	}
	categoryRepo = catRepo
}

//CategoryService - The service interface for working with categories.
type CategoryService struct {
}

//GetCategoryModel - Returns the requested Category Model.  For lifeapp, creates the default model if it does not yet exist for this user
func (t *CategoryService) GetCategoryModel(ctx context.Context, categoryModelID string) (*repository.CategoryUserModel, error) {
	catModel := repository.CategoryUserModel{}
	catModel.ID = categoryModelID
	catModel, err := categoryRepo.SelectOne(ctx, catModel)
	if err != nil {
		return nil, fmt.Errorf("Service Get Model Failed with: %v", err)
	}
	//First time user, initialize the base model
	if catModel.ID == "" && categoryModelID == MyLifeCategoryUserModelID {
		catModel.ID = MyLifeCategoryUserModelID
		catModel.Name = "Life Categories"
		//Future: make this more flexible for user types
		personalCat := model.NewCategory("Life")
		catModel.AddChild(*personalCat)
		workCat := model.NewCategory("Work")
		catModel.AddChild(*workCat)

		err = categoryRepo.Insert(ctx, catModel)
		if err != nil {
			return nil, fmt.Errorf("Could not initialize lifeapp user model with err: %v", err)
		}

	}
	return &catModel, nil
}

//ReplaceCategoryModel - Expert use only, replaces the full model
func (t *CategoryService) ReplaceCategoryModel(ctx context.Context, newUserModel *repository.CategoryUserModel) error {
	if newUserModel.ID == "" {
		return fmt.Errorf("Model id required for Replace")
	}
	err := categoryRepo.Update(ctx, *newUserModel)
	if err != nil {
		return fmt.Errorf("Category Model replace failed with: %v", err)
	}
	return nil
}

//DeleteCategoryModel - Expert use only, removes the full model
func (t *CategoryService) DeleteCategoryModel(ctx context.Context, userModelID string) error {
	if userModelID == "" {
		return fmt.Errorf("Model id required for Delete")
	}
	delTemplate := repository.CategoryUserModel{}
	delTemplate.ID = userModelID
	err := categoryRepo.Delete(ctx, delTemplate)
	if err != nil {
		return fmt.Errorf("Category Model delete failed with: %v", err)
	}
	return nil
}

//UpdateCategory updates the title of an existing category
func (t *CategoryService) UpdateCategory(ctx context.Context, categoryModelID string, updatedCategory model.Category) error {
	catModel := repository.CategoryUserModel{}
	catModel.ID = categoryModelID
	model, err := categoryRepo.SelectOne(ctx, catModel)
	if err != nil {
		return fmt.Errorf("Service Update Category  Failed with: %v", err)
	}
	if model.ID == "" {
		return errors.New("Service Update Category Model not found")
	}
	catItem, _ := model.FindChildByID(updatedCategory.ID)
	catItem.Title = updatedCategory.Title
	err = categoryRepo.Update(ctx, model)
	if err != nil {
		return fmt.Errorf("Category Model update failed with: %v", err)
	}
	return nil
}

//MoveCategory moves a categoryfrom one location to a new location
func (t *CategoryService) MoveCategory(ctx context.Context, categoryModelID string, newParentID string, categoryIDToMove string) error {
	catModel := repository.CategoryUserModel{}
	catModel.ID = categoryModelID
	model, err := categoryRepo.SelectOne(ctx, catModel)
	if err != nil {
		return fmt.Errorf("Service Update Category  Failed with: %v", err)
	}
	if model.ID == "" {
		return errors.New("Service Update Category Model not found")
	}
	if categoryIDToMove == "" {
		return errors.New("No category to move was selected")
	}
	catItem, _ := model.FindChildByID(newParentID)
	model.Move(categoryIDToMove, catItem)
	err = categoryRepo.Update(ctx, model)
	if err != nil {
		return fmt.Errorf("Category Model move failed with: %v", err)
	}

	return nil
}

//AddCategory - adds a category under the provided parent
func (t *CategoryService) AddCategory(ctx context.Context, categoryModelID string, newParentID string, newCategory model.Category) error {
	catModel := repository.CategoryUserModel{}
	catModel.ID = categoryModelID
	model, err := categoryRepo.SelectOne(ctx, catModel)
	if err != nil {
		return fmt.Errorf("Service Add Category  Failed with: %v", err)
	}
	if model.ID == "" {
		return errors.New("Service Add Category Model not found")
	}
	if newCategory.ID == "" {
		return errors.New("No new category id to add was selected")
	}
	catItem, _ := model.FindChildByID(newParentID)
	//If no parent was found, this is a root menu add
	if catItem.ID == "" {
		model.AddChild(newCategory)
	} else {
		catItem.AddChild(newCategory)
	}

	err = categoryRepo.Update(ctx, model)
	if err != nil {
		return fmt.Errorf("Category Model add failed with: %v", err)
	}

	return nil
}

//DeleteCategory - removes a category from the tree
func (t *CategoryService) DeleteCategory(ctx context.Context, categoryModelID string, categoryIDToDelete string) error {
	catModel := repository.CategoryUserModel{}
	catModel.ID = categoryModelID
	model, err := categoryRepo.SelectOne(ctx, catModel)
	if err != nil {
		return fmt.Errorf("Service Delete Category  Failed with: %v", err)
	}
	if model.ID == "" {
		return errors.New("Service Delete Category Model not found")
	}
	if categoryIDToDelete == "" {
		return errors.New("No new category id to delete was selected")
	}
	_, catItem := model.FindChildByID(categoryIDToDelete)
	catItem.RemoveChildByID(categoryIDToDelete)
	err = categoryRepo.Update(ctx, model)
	if err != nil {
		return fmt.Errorf("Category Model delete failed with: %v", err)
	}

	return nil
}

//NewCategoryService - returns a service interface for the category user model domain
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}
