package fileio

// FindRoot ...
func FindRoot(dr DirectoryReader) (DirectoryReader, bool) {

	if containsConfig(dr) {
		return dr, true
	}

	up, ok := dr.Up()

	if !ok {
		return nil, false
	}

	return FindRoot(up)

}

func containsConfig(dr DirectoryReader) bool {

	for _, dir := range dr.SubDirectories() {
		if dir == ".ele" {
			return true
		}
	}

	return false
}
