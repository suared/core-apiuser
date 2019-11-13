package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/suared/core/repository/dynamodb"
	"github.com/suared/core/security"
	"github.com/suared/core/ziptools"
)

//CategoryDAO - Caller would ipmlement the relevant DAO which would include the model to be saved along with any other key values (e.g. hash/sort for dynamo)
type CategoryDAO struct {
	CategoryHashKey string
	CategorySortKey string
	UserID          string

	//using zip for storage to keep dynamo costs low, hence removing the UserModel from unmarshal to replace with zip equivalent
	//because the life of a dao is only for a db interaction, handling the conversion in Refresh is fine
	CategoryUserModel     `json:"-"`
	CategoryUserModelData []byte
}

//HashKey - This is the value that would be set as the dynamo hashkey
func (dao *CategoryDAO) HashKey() string {
	return dao.CategoryHashKey
}

//SortKey - This is the value that would be set as the dynamo sortKey
func (dao *CategoryDAO) SortKey() string {
	return dao.CategorySortKey
}

//User - the user that made this call
func (dao *CategoryDAO) User() string {
	return dao.UserID
}

//New - creates a new instance of this specific type to support return values of the right type
func (dao *CategoryDAO) New() dynamodb.DAO {
	return new(CategoryDAO)
}

//Refresh - updates the Hashkey and SortKey.  Used by the library before calls
func (dao *CategoryDAO) Refresh() {
	dao.CategoryHashKey = "category_" + dao.UserID
	dao.CategorySortKey = dao.ID
}

//Populate - Called for any post processing after the DAO is generated from the library (e.g. calculated fields, unzip, etc)
func (dao *CategoryDAO) Populate() {
	var buf bytes.Buffer

	//Populate the model from the zip - First step is to unzip the bits
	err := ziptools.GetGunzipData(&buf, dao.CategoryUserModelData)
	if err != nil {
		panic(fmt.Errorf("Unable to unzip Category dao for hash: %v", dao.CategoryHashKey))
	}
	//Then marshall the bits into the model object for use (Note: intentionally not calling refresh here yet as DAOs are one time use and should not be necessary)
	err = json.Unmarshal(buf.Bytes(), &dao.CategoryUserModel)
	if err != nil {
		panic(fmt.Errorf("Unable to unmarshal unzip Category dao for hash: %v", dao.CategoryHashKey))
	}

}

//NewCategoryDAO - Initializes this object with the user ID from context
func NewCategoryDAO(ctx context.Context) *CategoryDAO {
	dao := new(CategoryDAO)
	dao.UserID = security.GetAuth(ctx).GetUser()
	return dao
}
