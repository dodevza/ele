package loaders

import "testing"

func Test_PathsWithDifferenceCasing_SortedIgnoringCasing(t *testing.T) {
	s := newPathSorter()
	sorted := s.Sort("Abd", "abc")

	if sorted[0] == "Abd" {
		t.Errorf("Case was ignored")
	}
}

func Test_PathsWithAfterKeyword_SortedAfterAtTheEnd(t *testing.T) {
	s := newPathSorter()
	sorted := s.Sort("after.file", "create.file")

	if sorted[0] == "after.file" {
		t.Errorf("After keyword was ignored")
	}
}

func Test_PathsWithBeforeKeyword_SortedBeforeAtTheStart(t *testing.T) {
	s := newPathSorter()
	sorted := s.Sort("create.file", "before.file")

	if sorted[0] != "before.file" {
		t.Errorf("Before keyword was ignored")
	}
}

func Test_PathAndFileAfterKeyword_SortedAtTheEndAndPathAfterLast(t *testing.T) {
	s := newPathSorter()
	sorted := s.Sort("after.file", "create.file", "create.after.file")

	if sorted[1] != "create.after.file" || sorted[2] != "after.file" {
		t.Errorf("After keyword not sorted correctly")
	}
}

func Test_PathAndFileBeforeKeyword_SortedAtTheStartAndPathBeforeFirst(t *testing.T) {
	s := newPathSorter()
	sorted := s.Sort("create.file", "create.before.file", "before.file")

	if sorted[1] != "create.before.file" || sorted[0] != "before.file" {
		t.Errorf("After keyword not sorted correctly")
	}
}

func Test_MutliplePathAfter_AllPathAftersAtTheEnd(t *testing.T) {
	s := newPathSorter()
	sorted := s.Sort("after.init.file", "after.cleanup.file", "create.file")

	if sorted[1] != "after.cleanup.file" || sorted[2] != "after.init.file" {
		t.Errorf("After path keywords keyword not sorted correctly, %s", sorted)
	}
}
