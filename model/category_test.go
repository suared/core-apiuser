package model

import (
	"testing"

	_ "github.com/suared/core/infra"
)

/*
Tests:
   -  Standard Category setup - CRUD
   -  "Move" Children from another parent (ensure levels reset)
   -  "Outdent" Children
   -  "Indent" Children
*/

//TODO: add findchildbyname tests
func TestBaseCategorySetup(t *testing.T) {
	root := NewCategoryRoot("Testing")
	personalCat := NewCategory("Personal")
	root.AddChild(*personalCat)
	workCat := NewCategory("Work")
	root.AddChild(*workCat)
	busCat := NewCategory("Entrepeneurial")
	root.AddChild(*busCat)

	//READ
	child := root.GetChildByName("Work")
	if child.Level != 1 && child.Title != "Work" {
		t.Errorf("Expected Level 1 Child Work, Got: %v", child)
	}

	child = root.GetChildByName("Entrepeneurial")
	if child.Level != 1 && child.Title != "Entrepeneurial" {
		t.Errorf("Expected Level 1 Child Entrepeneurial, Got: %v", child)
	}

	//UPDATE
	child = root.GetChildByName("Entrepeneurial")
	child.Title = "New Business"
	child = root.GetChildByName("New Business")
	if child.Level != 1 && child.Title != "New Business" {
		t.Errorf("Expected Level 1 Child Entrepeneurial, Got: %v", child)
	}

	//DELETE
	tempCategory := NewCategory("DELETEME")
	child.AddChild(*tempCategory)
	tempCategory = child.GetChildByName("DELETEME")
	if tempCategory.Level != 2 && tempCategory.Title != "DELETEME" {
		t.Errorf("Expected Level 2 Child DELETEME, Got: %v", tempCategory)
	}
	child.RemoveChildByName("DELETEME")
	tempCategory = child.GetChildByName("DELETEME")
	if tempCategory.ID != "" {
		t.Errorf("Expected Empty Struct after DELETEME, Got: %v", tempCategory)
	}

	//shallow copy is ok here (no children yet)
	tmpRecover := *child
	root.RemoveChildByName("New Business")
	child = root.GetChildByName("New Business")
	if child.ID != "" {
		t.Errorf("2 - Expected Empty Struct after Second Delete, Got: %v", child)
	}
	root.AddChild(tmpRecover)
	child = root.GetChildByName("New Business")
	if child.Level != 1 && child.Title != "New Business" {
		t.Errorf("2- Recovery of New Business Failed, Got: %v", child)
	}

	//L2 check & L3 Check
	child.AddChild(*getDisconnectedCategorySet())
	play := child.GetChildByName("Play")
	if play.Level != 2 && play.Title != "Play" {
		t.Errorf("Expected Play, Got: %v", play)
	}

	videos := play.GetChildByName("Videos")
	if videos.Level != 3 && videos.Title != "Videos" {
		t.Errorf("Expected Videos, Got: %v", videos)
	}

	beetlejuice := videos.GetChildByName("Beetlejuice")
	if beetlejuice.Level != 4 && beetlejuice.Title != "Beetlejuice" {
		t.Errorf("Expected Beetlejuice, Got: %v", beetlejuice)
	}

	outdentTestIDSave := beetlejuice.ID
	// Note: this is just the same as above at a deeper nesting level
	beetlejuice.AddChild(*getDisconnectedCategorySet())
	play = beetlejuice.GetChildByName("Play")
	if play.Level != 5 && play.Title != "Play" {
		t.Errorf("Expected Play, Got: %v", play)
	}

	videos = play.GetChildByName("Videos")
	if videos.Level != 6 && videos.Title != "Videos" {
		t.Errorf("Expected Videos, Got: %v", videos)
	}

	beetlejuice = videos.GetChildByName("Beetlejuice")
	if beetlejuice.Level != 7 && beetlejuice.Title != "Beetlejuice" {
		t.Errorf("Expected Beetlejuice, Got: %v", beetlejuice)
	}

	//Listing check
	personal := root.GetChildByName("Personal")
	personal.AddChild(*getDisconnectedCategorySet())
	list := root.GetAllChildren()
	if list[0].Title != "Personal" {
		t.Errorf("Listing 0 expected Personal, Got: %v", list[0])
	}

	if list[1].Title != "Play" {
		t.Errorf("Listing 1 expected Play, Got: %v", list[1])
	}

	var currentLevel int
	for i := range list {
		if list[i].Level > currentLevel {
			//went to a child - can only go down 1 level at a time
			if list[i].Level != currentLevel+1 {
				t.Errorf("Expected only 1 level move, failed with prev: %v, current: %v", list[i-1], list[i])
			}
			//can go up as many levels as makes sense, same for = so no tests otherwise
			currentLevel = list[i].Level
		}
	}

	/* At this point looks like:

	Personal
		Play
		- Games
		- Videos
			-Beetlejuice
			-Something About Mary
		- Music
	Work
	New Business
		Play
		- Games
		- Videos
			-Beetlejuice
				Play  (** For outdent/ indent test next, going to move this back, then indent under Something about Mary, then direct move back so it ends here **)
				- Games
				- Videos
					-Beetlejuice
					-Something About Mary
				- Music
			-Something About Mary
		- Music
	*/

	if list[13].Title != "Play" {
		t.Errorf("Listing 13 expected Play, Got: %v, fulllist: %v", list[13], list)
	}

	//Save this version to enable rerun of the tests later to validate file recovery
	saveCategoryRootToFile("CtxRoot-Initial.json", root)

	//Move Tests..
	//child is Beetlejuice, Parent is Videos in above picture, see above outdent/ indent test
	targetChild, targetParent := root.FindChildByID(outdentTestIDSave)
	targetPlay := targetChild.GetChildByName("Play")
	targetNewParent := targetParent.GetChildByName("Something About Mary")
	root.Move(targetPlay.ID, targetNewParent)
	if targetNewParent.GetChildByName("Play").ID != targetPlay.ID {
		t.Error("Move failed, expected child under Something About Mary")
	}

	saveCategoryRootToFile("CtxRoot-SAMMoveCheck.json", root)

	list = root.GetAllChildren()
	//Make sure all children travelled
	if list[14].Title != "Play" {
		t.Errorf("Listing 14 expected Play, Got: %v, fulllist: %v", list[14], list)
	}

	if list[15].Title != "Games" {
		t.Errorf("Listing 15 expected Movies, Got: %v, fulllist: %v", list[15], list)
	}

	//Reset to previous state and then do the same with convenience methods
	root.Move(targetPlay.ID, targetChild)

	newRoot := getCategoryRootFromFile("CtxRoot-Initial.json")
	list = newRoot.GetAllChildren()
	if list[13].Title != "Play" {
		t.Errorf("FileRestore: Listing 13 expected Play, Got: %v, fulllist: %v", list[13], list)
	}

	if !newRoot.Equals(root) {
		t.Errorf("restore equals post-baseline check failed!  root: %v", newRoot)
		saveCategoryRootToFile("CtxRoot-postbaselineroot.json", root)
		saveCategoryRootToFile("CtxRoot-postbaselinenew.json", newRoot)

	}

	list = root.GetAllChildren()
	if list[13].Title != "Play" {
		t.Errorf("Listing 13 expected Play, Got: %v, fulllist: %v", list[13], list)
	}

	//root.Outdent/ Indent- repeating previous test with convenience methods
	//child is Beetlejuice, Parent is Videos in above picture, see above outdent/ indent test
	targetChild, targetParent = root.FindChildByID(outdentTestIDSave)
	targetChild = targetChild.GetChildByName("Play")
	root.Outdent(targetChild.ID)
	root.Outdent(targetChild.ID)
	root.Indent(targetChild.ID, targetParent.GetChildByName("Something About Mary").ID)
	list = root.GetAllChildren()
	//Make sure all children travelled
	if list[14].Title != "Play" {
		t.Errorf("Listing 14 expected Play, Got: %v, fulllist: %v", list[14], list)
	}

	if list[15].Title != "Games" {
		t.Errorf("Listing 15 expected Movies, Got: %v, fulllist: %v", list[15], list)
	}

	newRoot = getCategoryRootFromFile("CtxRoot-SAMMoveCheck.json")
	if !newRoot.Equals(root) {
		t.Errorf("restore equals SAM Checkfailed!  root: %v", newRoot)
		saveCategoryRootToFile("CtxRoot-postbaselineroot.json", root)
		saveCategoryRootToFile("CtxRoot-postbaselinenew.json", newRoot)

	}

}

func TestCategorySliceRemover(t *testing.T) {
	testSlice := removeCategorySliceIndex(getCategoryArray(), 0)
	if testSlice[0].Title != "test2" {
		t.Error("Front of slice delete error")
	}

	testSlice = removeCategorySliceIndex(getCategoryArray(), 4)
	if testSlice[len(testSlice)-1].Title != "test4" {
		t.Error("End of slice delete error")
	}

	testSlice = removeCategorySliceIndex(getCategoryArray(), 3)
	if testSlice[len(testSlice)-1].Title != "test5" {
		t.Error("Middle of slice delete error (1)")
	}
	if testSlice[3].Title != "test5" {
		t.Errorf("Middle of slice delete error (2), slice: %v", testSlice)
	}
	if testSlice[2].Title != "test3" {
		t.Error("Middle of slice delete error (3)")
	}
}

func getCategoryArray() []*Category {
	var catSlice []*Category
	cat1 := NewCategory("test1")
	cat2 := NewCategory("test2")
	cat3 := NewCategory("test3")
	cat4 := NewCategory("test4")
	cat5 := NewCategory("test5")
	return append(catSlice, cat1, cat2, cat3, cat4, cat5)
}
func getDisconnectedCategorySet() *Category {
	/*
		Play
			- Games
			- Videos
				-Beetlejuice
				-Something About Mary
			- Music
	*/
	playCat := NewCategory("Play")
	playCat.AddChild(*NewCategory("Games"))
	videoCat := NewCategory("Videos")
	videoCat = playCat.AddChild(*videoCat)
	videoCat.AddChild(*NewCategory("Beetlejuice"))
	videoCat.AddChild(*NewCategory("Something About Mary"))
	playCat.AddChild(*NewCategory("Music"))
	return playCat
}
