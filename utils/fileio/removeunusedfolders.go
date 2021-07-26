package fileio

// RemoveUnUsedFolders ...
func RemoveUnUsedFolders(from DirectoryReader, upto DirectoryReader) {
	curr := from

	found := false
	for found == false {
		if (len(curr.Files()) == 0 && len(curr.SubDirectories()) == 0) == false {
			return
		}

		up, ok := curr.Up()
		if !ok {
			return
		}
		up.RemoveDirectory(curr.Name())
		curr = up

		found = curr == upto
	}
}
