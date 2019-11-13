package model

import (
	"strconv"

	"github.com/suared/core/uuid"

	//Setup environment
	_ "github.com/suared/core/infra"
)

//Category - The model object
type Category struct {
	ID       string      `json:"id"`
	Level    int         `json:"level"`
	Title    string      `json:"title"`
	Children []*Category `json:"categories"`
}

//Implement stringer interface to facilitate debugging
func (cat *Category) String() string {
	return cat.ID + "," + strconv.Itoa(cat.Level) + "," + cat.Title + ", children:" + strconv.Itoa(len(cat.Children)) + "\n"
}

//Equals for debugging
func (cat *Category) Equals(compare *Category) bool {
	if cat.ID != compare.ID {
		return false
	}
	if cat.Level != compare.Level {
		return false
	}
	if cat.Title != compare.Title {
		return false
	}
	if cat.ID != compare.ID {
		return false
	}

	baseLength := len(cat.Children)
	compareLength := len(compare.Children)
	if baseLength != compareLength {
		return false
	}

	//Intentionally making order critical for now, may revisit in the future
	for i := range cat.Children {
		if !(cat.Children[i].Equals(compare.Children[i])) {
			return false
		}
	}

	return true
}

//AddChild - Adds a new child and returns the pointer to the newly created Child
func (cat *Category) AddChild(child Category) *Category {
	//Ensure the level is correct
	child.Level = cat.Level + 1
	//Disconnect the children so the levels can be updated as needed
	children := child.Children
	child.Children = nil
	//Add back the children to reset levels - this works as it is recursive in nature
	for i := range children {
		child.AddChild(*children[i])
	}
	cat.Children = append(cat.Children, &child)
	return &child
}

//GetChildByName - Returns first child with matching name, immediate child search only
func (cat *Category) GetChildByName(name string) *Category {
	for i := range cat.Children {
		if cat.Children[i].Title == name {
			return cat.Children[i]
		}
	}
	return &Category{}
}

//FindChildByName - Returns first child with matching name followed by the parent Category if not the root, recursive search
func (cat *Category) FindChildByName(name string) (*Category, *Category) {
	for i := range cat.Children {
		if cat.Children[i].Title == name {
			return cat.Children[i], cat
		}
		foundChild, foundParent := cat.Children[i].FindChildByName(name)
		if foundChild.ID != "" {
			return foundChild, foundParent
		}
	}
	return &Category{}, &Category{}
}

//FindChildByID - Returns first child with matching id followed by the parent Category if not the root, recursive search
func (cat *Category) FindChildByID(id string) (*Category, *Category) {
	for i := range cat.Children {
		if cat.Children[i].ID == id {
			return cat.Children[i], cat
		}
		foundChild, foundParent := cat.Children[i].FindChildByID(id)
		if foundChild.ID != "" {
			return foundChild, foundParent
		}
	}
	return &Category{}, &Category{}
}

//RemoveChildByName - Removes first child with matching name, immediate child search only
func (cat *Category) RemoveChildByName(name string) {
	cat.Children = removeCategoryItemByName(cat.Children, name)
}

//RemoveChildByID - Removes first child with matching ID, will search recursively through children
func (cat *Category) RemoveChildByID(id string) {
	//First Find the Category
	_, parent := cat.FindChildByID(id)
	//Remove it
	parent.Children = removeCategoryItemByID(parent.Children, id)
}

func removeCategoryItemByName(list []*Category, name string) []*Category {
	var foundDelete bool
	var indexDelete int
	for i := range list {
		if list[i].Title == name {
			foundDelete = true
			indexDelete = i
			break
		}
	}

	if foundDelete {
		return removeCategorySliceIndex(list, indexDelete)
	}
	return list

}

func removeCategoryItemByID(list []*Category, id string) []*Category {
	var foundDelete bool
	var indexDelete int
	for i := range list {
		if list[i].ID == id {
			foundDelete = true
			indexDelete = i
			break
		}
	}

	if foundDelete {
		return removeCategorySliceIndex(list, indexDelete)
	}
	return list
}

//GetAllChildren - returns a sorted array of the category tree
func (cat *Category) GetAllChildren() []*Category {
	var categoryArray []*Category
	var child *Category
	for idx := range cat.Children {
		child = cat.Children[idx]
		categoryArray = append(categoryArray, cat.Children[idx])
		categoryArray = append(categoryArray, child.GetAllChildren()...)
	}
	return categoryArray
}

//NewCategory - Constructs a new Category
func NewCategory(title string) *Category {
	return &Category{ID: uuid.NewUUID(), Level: 1, Title: title}
}

//CategoryRoot - The base of the category tree
type CategoryRoot struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Children []*Category `json:"categories"`
}

//Equals - compares two roots for equality
func (root *CategoryRoot) Equals(compare *CategoryRoot) bool {
	if root.ID != compare.ID {
		return false
	}
	if root.Name != compare.Name {
		return false
	}

	rootLength := len(root.Children)
	compareLength := len(compare.Children)
	if rootLength != compareLength {
		return false
	}

	rootChildRep := root.GetAllChildren()
	compareChildRep := compare.GetAllChildren()

	//Intentionally making order critical for now, may revisit in the future
	for i := range rootChildRep {
		if !(rootChildRep[i].Equals(compareChildRep[i])) {
			return false
		}
	}

	return true
}

//AddChild - Adds a new level 1 child to the root
func (root *CategoryRoot) AddChild(category Category) *Category {
	category.Level = 1
	root.Children = append(root.Children, &category)
	return &category
}

//RemoveChildByName - Removes first child with matching name, immediate child search only
func (root *CategoryRoot) RemoveChildByName(name string) {
	root.Children = removeCategoryItemByName(root.Children, name)
}

//RemoveChildByID - Removes first child with matching ID, will search recursively through children
func (root *CategoryRoot) RemoveChildByID(id string) {
	//First Find the Category
	_, parent := root.FindChildByID(id)
	//Remove it from root if parent is empty, otherwise from parent
	if parent.ID == "" {
		root.Children = removeCategoryItemByID(root.Children, id)
	} else {
		parent.Children = removeCategoryItemByID(parent.Children, id)
	}
}

//GetChildByName - Returns first child with matching name, immediate child search only
func (root *CategoryRoot) GetChildByName(name string) *Category {
	for i := range root.Children {
		if root.Children[i].Title == name {
			return root.Children[i]
		}
	}
	return &Category{}
}

//FindChildByName - Returns first child with matching name,recursive.  Note: if category root, the parent returns nil to signify the root was reached, otherwise all children would return an empty parent to signify not the root
func (root *CategoryRoot) FindChildByName(name string) (*Category, *Category) {
	for i := range root.Children {
		if root.Children[i].Title == name {
			return root.Children[i], nil
		}
		foundChild, foundParent := root.Children[i].FindChildByName(name)
		if foundChild.ID != "" {
			return foundChild, foundParent
		}
	}
	return &Category{}, nil
}

//FindChildByID - Returns first child with matching id followed by the parent Category if not the root, recursive search
func (root *CategoryRoot) FindChildByID(id string) (*Category, *Category) {
	for i := range root.Children {
		if root.Children[i].ID == id {
			return root.Children[i], nil
		}
		foundChild, foundParent := root.Children[i].FindChildByID(id)
		if foundChild.ID != "" {
			return foundChild, foundParent
		}
	}
	return &Category{}, &Category{}
}

//GetAllChildren - returns a sorted array of the category tree
func (root *CategoryRoot) GetAllChildren() []*Category {
	var categoryArray []*Category
	var child *Category
	for idx := range root.Children {
		child = root.Children[idx]
		categoryArray = append(categoryArray, root.Children[idx])
		categoryArray = append(categoryArray, child.GetAllChildren()...)
	}
	return categoryArray
}

//Move - relocates the category identified by ID as a child to the provided NewParent Category
func (root *CategoryRoot) Move(catID string, NewParent *Category) {
	currentCat, currentParent := root.FindChildByID(catID)
	//If parent is null, this is moving out of root
	if currentParent.ID == "" {
		root.RemoveChildByID(currentCat.ID)
	} else {
		currentParent.RemoveChildByID(currentCat.ID)
	}
	//If New Parent is null this is moving into root
	if NewParent.ID == "" {
		root.AddChild(*currentCat)
	} else {
		NewParent.AddChild(*currentCat)
	}
}

//Outdent - Move out one level
func (root *CategoryRoot) Outdent(catID string) {
	_, currentParent := root.FindChildByID(catID)
	if currentParent.ID == "" {
		//ignore, can't outdent further
		return
	}
	_, currentGrandParent := root.FindChildByID(currentParent.ID)

	root.Move(catID, currentGrandParent)

}

//Indent - Move under a new Parent
func (root *CategoryRoot) Indent(catID, newParentID string) {
	futureParentCat, _ := root.FindChildByID(newParentID)
	root.Move(catID, futureParentCat)
}

//NewCategoryRoot - Constructs a new Root Category Holder
func NewCategoryRoot(name string) *CategoryRoot {
	return &CategoryRoot{ID: uuid.NewUUID(), Name: name}
}

//removeCategorySliceIndex - Returns a new slice with the removed item
//TODO: Autogen these helpers till generics exist
func removeCategorySliceIndex(list []*Category, idx int) []*Category {

	//front of line delete
	if idx == 0 {
		list = list[1:len(list)]
		return list
	}
	//end of line delete
	if idx == len(list)-1 {
		list = list[0 : len(list)-1]
		return list
	}
	//middle delete
	firstHalf := list[0:idx]
	secondHalf := list[idx+1 : len(list)]
	list = append(firstHalf, secondHalf...)
	return list

}
