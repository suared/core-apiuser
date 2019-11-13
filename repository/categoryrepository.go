package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/suared/core/infra"
	"github.com/suared/core/repository"
	"github.com/suared/core/repository/dynamodb"
	"github.com/suared/core/ziptools"
)

//CategoryRepository - Database interface to the Category table
type CategoryRepository struct {
	config  repository.Config
	session repository.Session
}

//Config - Returns the current configuration
func (repo *CategoryRepository) Config() repository.Config {
	return repo.config
}

//DAO - Returns a DAO associated with this repository from a model object
func (repo *CategoryRepository) DAO(ctx context.Context, userModel CategoryUserModel, zipme bool, active bool, audit bool) (dynamodb.DAO, error) {
	dao := NewCategoryDAO(ctx)
	dao.CategoryUserModel = userModel

	if zipme == true {
		dao.CategoryUserModelData = ziptools.GetGzipDataFromStruct(userModel)
	}
	return dao, nil
}

//Insert - Sample of a basic insert method with validation
func (repo *CategoryRepository) Insert(ctx context.Context, userModel CategoryUserModel) error {
	// Populate the Data object First //  active?, audit?
	dao, err := repo.DAO(ctx, userModel, true, false, false)

	if err != nil {
		log.Printf("Unable to Insert DAO, error Getting DAO: %v", err)
		return err
	}

	// Repository layer is responsible for validating auth rules
	err = dynamodb.ValidAction(ctx, "insert", dao)
	if err != nil {
		return err
	}

	return dynamodb.InsertOrUpdate(ctx, repo, dao)
}

//Update - Sample of updating a DB entry
func (repo *CategoryRepository) Update(ctx context.Context, userModel CategoryUserModel) error {
	dao, err := repo.DAO(ctx, userModel, true, false, false)
	if err != nil {
		log.Printf("Unable to Update, error getting DAO, err: %v", err)
		return err
	}

	err = dynamodb.ValidAction(ctx, "update", dao)
	if err != nil {
		return err
	}

	return dynamodb.InsertOrUpdate(ctx, repo, dao)

}

//Delete - Sample of deleting a DB entry
func (repo *CategoryRepository) Delete(ctx context.Context, template CategoryUserModel) error {
	dao, err := repo.DAO(ctx, template, false, false, false)
	if err != nil {
		log.Printf("Unable getting Dao in Delete, err: %v", err)
		return err
	}

	return dynamodb.Delete(ctx, repo, dao)
}

//Select - Sample of a get all by hashkey
func (repo *CategoryRepository) Select(ctx context.Context, template CategoryUserModel) ([]CategoryUserModel, error) {
	dao, err := repo.DAO(ctx, template, false, false, false)
	if err != nil {
		log.Printf("Unable to Select, error getting DAO, err: %v", err)
		return nil, err
	}

	result, err := dynamodb.Select(ctx, repo, dao)

	var outputList []CategoryUserModel
	//since the search is for user, validation only needs to occur on one item..
	var validated bool
	for i := range result {
		resultDAO := result[i]
		//Check once only...
		if !validated {
			//ignore if no result/ empty
			if resultDAO.HashKey() != "" {
				err = dynamodb.ValidAction(ctx, "selectAll", resultDAO)

				if err != nil {
					return []CategoryUserModel{}, err
				}
			}
		}
		//Convert DAO to Request here then add to list
		categoryDao, ok := resultDAO.(*CategoryDAO)
		if !ok {
			return []CategoryUserModel{}, errors.New("Unable to convert back to categoryDao, DB results unexpected")
		}
		resultItem := categoryDao.CategoryUserModel
		outputList = append(outputList, resultItem)
		//log.Printf("output list: %v", outputList)
	}

	return outputList, err

}

//SelectOne - Returns one model object, can be empty if no results
func (repo *CategoryRepository) SelectOne(ctx context.Context, template CategoryUserModel) (CategoryUserModel, error) {
	dao, err := repo.DAO(ctx, template, false, false, false)
	if err != nil {
		log.Printf("Unable to SelectOne, error getting DAO, err: %v", err)
		return CategoryUserModel{}, err
	}

	result, err := dynamodb.SelectOne(ctx, repo, dao)

	//Convert DAO to Request here then add to list
	categoryDao, ok := result.(*CategoryDAO)
	if !ok {
		return CategoryUserModel{}, errors.New("Unable to convert back to categoryDao, DB results unexpected")
	}
	resultItem := categoryDao.CategoryUserModel

	if resultItem.ID != "" {
		err = dynamodb.ValidAction(ctx, "selectOne", result)

		if err != nil {
			return CategoryUserModel{}, err
		}
	}

	return resultItem, err

}

//SetSession - enables the library to store/ reuse the session for efficiency vs. creating new on each call
func (repo *CategoryRepository) SetSession(session repository.Session) {
	repo.session = session
}

//Session - Returns the session associated with this repository
func (repo *CategoryRepository) Session() repository.Session {
	return repo.session
}

//NewCategoryRepository - Initializes a sample repository with config values set
func NewCategoryRepository() (*CategoryRepository, error) {
	repo := new(CategoryRepository)
	configMap := repository.NewBasicConfig("categoryDatabase")
	configMap.AddEntry("backend", os.Getenv("PROCESS_REPOSITORY"))
	configMap.AddEntry("table", os.Getenv("PROCESS_AWS_DYNAMOTABLE_CATEGORY"))
	configMap.AddEntry("region", os.Getenv("PROCESS_AWS_REGION"))
	configMap.AddEntry("endpoint", os.Getenv("PROCESS_AWS_DYNAMOENDPOINT"))
	configMap.AddEntry("rcu", os.Getenv("PROCESS_AWS_DYNAMOTABLE_RCU"))
	configMap.AddEntry("wcu", os.Getenv("PROCESS_AWS_DYNAMOTABLE_WCU"))

	configMap.AddEntry("hashKeyName", "CategoryHashKey")
	configMap.AddEntry("sortKeyName", "CategorySortKey")
	configMap.AddEntry("env", os.Getenv("PROCESS_ENV"))

	repo.config = configMap

	//Convert the config into an initialized dynamoo table
	repositoryInit, err := dynamodb.CreateTable(repo)
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize category database session, received error: %v", err)
	}

	repository, ok := repositoryInit.(*CategoryRepository)

	if !ok {
		return nil, fmt.Errorf("Repository Category cast did not succeed, have a: %v", repositoryInit)
	}

	return repository, nil
}
