package repository

import (
	"github.com/suared/core-apiuser/model"

	"github.com/suared/core/uuid"
)

//CategoryUserModel - repository model object to enable future non-direct model adds where appropriate. Intentionally saving/ enabling only this tier for customizations thus far
type CategoryUserModel struct {
	model.CategoryRoot
}

//NewCategoryUserModel - initializes the user model
func NewCategoryUserModel(name string) *CategoryUserModel {
	catModel := new(CategoryUserModel)
	catModel.ID = uuid.NewUUID()
	catModel.Name = name
	return catModel
}
