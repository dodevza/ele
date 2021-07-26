package project

import (
	"ele/constants"
	"ele/utils/doc"
	"fmt"
	"strings"
)

// Tag ...
type Tag struct {
	Name     string
	IsActive bool
	IsTag    bool
}

// Tags ...
type Tags []Tag

// Print ...
func (tags Tags) Print() int {
	rows := doc.NewTable().
		Column("Tag", 56).Column("Is Custom", 12).Column("Is Active", 12).
		StartRows()
	for _, tag := range tags {
		rows.Row(tag.Name, doc.YesNo(tag.IsTag), doc.YesNo(tag.IsActive))
	}

	count := len(tags)

	if count > 0 {
		rows.Divider()
		doc.Line(doc.Infof("%d tags", count))
	}
	return count
}

// GetTags ...
func (project *Project) GetTags() Tags {
	tree := project.QueryFromRoot().Tags()

	tags := make(Tags, 0)
	versions := tree.ToArray()

	for _, v := range versions {
		if len(v.Name) == 0 || v.Name == constants.REPEATABLE {
			continue
		}
		isTag := tree.IsTag(v.Name)
		tag := Tag{Name: v.Name, IsTag: isTag}
		tags = append(tags, tag)
	}

	return tags
}

// AddTag ...
func (project *Project) AddTag(tag string) {

	key := strings.ToUpper(tag)
	if project.tagMap[key] {
		return
	}
	project.config.Tags = append(project.config.Tags, tag)
	project.tagMap[key] = true
	project.tags = append(project.tags, tag)

	project.SaveConfig()
}

// RemoveTag ...
func (project *Project) RemoveTag(tag string) {

	key := strings.ToUpper(tag)
	if project.tagMap[key] == false {
		return
	}
	tags := make([]string, 0)
	for _, t := range project.config.Tags {
		if t != tag {
			tags = append(tags, t)
		}
	}
	project.config.Tags = tags
	delete(project.tagMap, key)
	project.tags = tags

	project.SaveConfig()
}

// InspectTag ...
func (project *Project) InspectTag(tag string) int {
	tree := project.QueryFromRoot().Tags()
	migrations := tree.Forward(tag, tag)
	migrations.PrintExcecutionPlan()
	return len(migrations)
}

// ActivateTags ...
func (project *Project) ActivateTags(start string, end string) error {
	workspace := project.root.CreateDirectory(".ele")
	err := workspace.MoveFile("active", "active_")
	appender := workspace.FileAppender("active")
	appender.Text(fmt.Sprintf("%s:%s", start, end))
	appender.Close()
	if err == nil {
		workspace.RemoveFile("active_")
	}

	project.ActiveStart = start
	project.ActiveEnd = end
	return nil
}

// DeactivateTags ...
func (project *Project) DeactivateTags() error {
	workspace := project.root.CreateDirectory(".ele")
	workspace.RemoveFile("active")
	project.ActiveStart = ""
	project.ActiveEnd = ""
	return nil
}

func (project *Project) loadActiveTags() {

	workingDir := project.root.CreateDirectory(".ele")
	scn, err := workingDir.FileScanner("active")
	if err != nil {
		return
	}
	if scn.Scan() {
		activeVersions := scn.Text()
		start, end := StringToTagRange(activeVersions)
		project.ActiveStart = start
		project.ActiveEnd = end
	}
}

// StringToTagRange ...
func StringToTagRange(tagValue string) (string, string) {
	var start strings.Builder
	var end strings.Builder
	foundSep := false
	for _, s := range tagValue {
		if s == ':' {
			foundSep = true
			continue
		}
		if !foundSep {
			start.WriteString(string(s))
		} else {
			end.WriteString(string(s))
		}
	}

	if !foundSep {
		return start.String(), start.String()
	}

	return start.String(), end.String()
}
