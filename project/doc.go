package project

import (
	"ele/utils/doc"
)

// PrintStatus ...
func (project *Project) PrintStatus() {
	count := len(project.QueryFromRoot().Limit(project.ActiveStart, project.ActiveEnd))
	if count == 0 {
		doc.Line(doc.Info("No migrations"))
	} else {
		doc.Line(doc.Infof("%d migrations", count))
	}
	doc.NewLine()
	project.PrintBadges()
}

// PrintBadges ...
func (project *Project) PrintBadges() {
	name := project.Environments.ActiveName()
	if len(name) > 0 {
		doc.Line(doc.Env(project.Environments.ActiveName()), doc.Tag(project.ActiveStart, project.ActiveEnd))
	} else {
		doc.Line(doc.Tag(project.ActiveStart, project.ActiveEnd))
	}

	doc.NewLine()
}
