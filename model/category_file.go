package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//GetCategoryRootFromBytes - returns a Category Root from byte representation to retrieve from storage
func GetCategoryRootFromBytes(data []byte) (*CategoryRoot, error) {
	var catRoot *CategoryRoot

	catRoot = NewCategoryRoot("FileRetrieve")
	//no interfaces and no reuse here so potentially can be done by default encoders, let's try that first
	err := json.Unmarshal(data, catRoot)
	if err != nil {
		log.Printf("Unmarshall Category Root failed with %v", err)
		return nil, err
	}
	return catRoot, nil
}

//ConvertCategoryRootToBytes - return byte[] representation to enable storage
func ConvertCategoryRootToBytes(root *CategoryRoot) ([]byte, error) {
	data, err := json.Marshal(root)
	return data, err
}

//getCategoryRootFromFile - Restore Category Root From File
func getCategoryRootFromFile(fileloc string) *CategoryRoot {

	var theRoot *CategoryRoot

	f, err := os.Open(os.Getenv("CATEGORY_MODEL_TESTFILE_DIR") + fileloc)

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Received error on open file e: %v", err)
		panic(err)
	}
	defer f.Close()

	theRoot, err = GetCategoryRootFromBytes(data)
	if err != nil {
		log.Printf("Received error retrieving Bytes from file e: %v", err)
		panic(err)
	}

	return theRoot
}

func saveCategoryRootToFile(filename string, root *CategoryRoot) {
	fileDir := os.Getenv("CATEGORY_MODEL_TESTFILE_DIR")
	f, err := os.Create(fileDir + filename)
	if err != nil {
		log.Printf("Received error on create file e: %v", err)
		panic(err)
	}
	defer f.Close()

	data, err := ConvertCategoryRootToBytes(root)
	if err != nil {
		log.Printf("Received error on converting category to file e: %v", err)
		panic(err)
	}
	f.Write(data)
	f.Sync()
}
