package migrations

import "testing"

func Test_Version_SequencialyAdded_CorrectOrder(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.0.1")
	tree.addVersion("0.1.0")

	result := tree.ToString()

	if result != "0.0.1,0.1.0" {
		t.Errorf("Incorrect Order")
	}
}

func Test_Version_UnsequencialyAdded_CorrectOrder(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.1.0")
	tree.addVersion("0.0.1")

	result := tree.ToString()

	if result != "0.0.1,0.1.0" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_Version_UnsequencialyAdded3Items_CorrectOrder(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.1.0")
	tree.addVersion("0.0.1")
	tree.addVersion("2.0.1")

	result := tree.ToString()

	if result != "0.0.1,0.1.0,2.0.1" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_Version_AddNonNumericalTags_NonNumericalTagsAtTheEnd(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.1.0")
	tree.addVersion("NEXT")
	tree.addVersion("2.0.1")

	result := tree.ToString()

	if result != "0.1.0,2.0.1,NEXT" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_Version_AddSemanticTags_SematicTagsAfterTheirOwnNumericTags(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.1.0")
	tree.addVersion("0.1.0-NEXT")
	tree.addVersion("2.0.1-NEXT")
	tree.addVersion("2.0.1")

	result := tree.ToString()

	if result != "0.1.0,0.1.0-NEXT,2.0.1,2.0.1-NEXT" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_VersionsWithTags_TagsOrderedAtTheEndAndInSquenceItWasSpecified(t *testing.T) {
	tree := EmptyVersionTreeWithTags([]string{"UAT", "QA"})

	tree.addVersion("V0.1.0")
	tree.addVersion("V2.0.1")

	result := tree.ToString()

	if result != "V0.1.0,V2.0.1,UAT,QA" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_VersionsWithDuplicateTags_TagsOrderedAtTheEndAndInSquenceItWasSpecified(t *testing.T) {
	tree := EmptyVersionTreeWithTags([]string{"UAT", "QA"})

	tree.addVersion("V0.1.0")
	tree.addVersion("QA")
	tree.addVersion("V2.0.1")
	tree.addVersion("UAT")

	result := tree.ToString()

	if result != "V0.1.0,V2.0.1,UAT,QA" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_MultipleTags_TagsOrderedAtTheEndAndInSquenceItWasSpecified(t *testing.T) {
	tree := EmptyVersionTreeWithTags([]string{"UAT", "QA", "PROD"})

	result := tree.ToString()

	if result != "UAT,QA,PROD" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_Version_UnderscoreVersions_UnderScoreVersionsShouldBeAtTheTop(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.1.0")
	tree.addVersion("_init")
	tree.addVersion("2.0.1")

	result := tree.ToString()

	if result != "_INIT,0.1.0,2.0.1" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_Version_UnderscoreDepthVersions_MoreUnderScoresShouldBeAfterLess(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("0.1.0")
	tree.addVersion("__1")
	tree.addVersion("_1")

	result := tree.ToString()

	if result != "_1,__1,0.1.0" {
		t.Errorf("Incorrect Order: %s", result)
	}
}

func Test_RepeatableWithTag_TagAfterRepeatable(t *testing.T) {
	tree := EmptyVersionTree()

	tree.addVersion("integration")
	tree.addVersion("~")

	result := tree.ToString()

	if result != "INTEGRATION,~" {
		t.Errorf("Incorrect Order: %s", result)
	}
}
