package project

import (
	"ele/utils/fileio"
	"errors"
)

// PromoteOptions ...
type PromoteOptions struct {
	Target    string
	CreateTag bool
}

// Promote ...
func (project *Project) Promote(options *PromoteOptions) error {

	if len(options.Target) == 0 {
		return errors.New("No target provided")
	}
	if project.IsValidTag(options.Target) == false {
		if options.CreateTag {
			project.AddTag(options.Target)
		} else {
			return errors.New("Invalid target provided, consider providing CreateTag argument")
		}
	}
	source := project.Search(&SearchOptions{OnlyStaged: true, VersionStart: project.ActiveStart, VersionEnd: project.ActiveEnd})

	if len(*source) == 0 {
		source = project.Search(&SearchOptions{VersionStart: project.ActiveStart, VersionEnd: project.ActiveEnd})
	}

	if len(*source) == 0 {
		return errors.New("Nothing to promote")
	}

	for _, mig := range *source {

		var dir fileio.DirectoryReader = project.dir
		if mig.VersionDir != project.dir {
			upOneDir, ok := mig.VersionDir.Up()
			if ok {
				dir = upOneDir
			}
		}

		targetDir := dir.CreateDirectory(options.Target)
		fromPath := fileio.AppendPathChar(mig.VersionDir.Name()) + fileio.AppendPathChar(mig.PathFromVersion) + mig.FileName
		toPath := fileio.AppendPathChar(targetDir.Name()) + fileio.AppendPathChar(mig.PathFromVersion) + mig.FileName

		dir.MoveFile(fromPath, toPath)
		fileio.RemoveUnUsedFolders(mig.Dir, mig.VersionDir)
	}

	project.stagingArea.Clear()

	return nil
}
