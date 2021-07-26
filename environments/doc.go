package environments

import (
	"ele/utils/doc"
	"strings"
)

// Print ...
func (state EnvironmentState) Print(environments EnvironmentCollection) int {
	rows := doc.NewTable().
		Column("Environment", 15).Column("Database", 15).Column("Groups", 50).
		StartRows()
	for _, env := range environments {
		groups := make([]string, 0)
		for grp, values := range state.Groups {
			key := strings.ToLower(env.Name)
			_, foundInGrp := values[key]
			if foundInGrp {
				groups = append(groups, grp)
			}
		}
		dbname := ""
		if env.Database != nil {
			dbname = env.Database.DBName
		}
		groupNames := strings.Join(groups, ", ")
		rows.Row(env.Name, dbname, groupNames)
	}

	count := len(environments)

	if count > 0 {
		rows.Divider()
		doc.Line(doc.Infof("%d environments", count))
	}
	return count
}

// Print ...
func (env *Environment) Print() {
	doc.Paragraph(doc.Data(env.Name))
	if env.Database != nil {

		doc.Line(doc.Info(doc.FitLeft("Database", 40)), doc.Info(doc.FitRight(env.Database.DBName, 40)))
		doc.Line(doc.Info(doc.FitLeft("Driver", 40)), doc.Info(doc.FitRight(env.Database.Driver, 40)))
		doc.Line(doc.Info(doc.FitLeft("User", 40)), doc.Info(doc.FitRight(env.Database.User, 40)))
		doc.Line(doc.Info(doc.FitLeft("Password", 40)), doc.Info(doc.FitRight(env.Database.Password, 40)))
		doc.Line(doc.Info(doc.FitLeft("SSL Mode", 40)), doc.Info(doc.FitRight(env.Database.SSLMode, 40)))
	}
	doc.NewLine()
}

// Inspect ...
func (state *EnvironmentState) Inspect(name string) {
	environments := state.ByName(name)

	count := len(environments)

	if count == 0 {
		doc.Paragraph(doc.Infof("No environment %s found", name))
	}

	for i, env := range environments {
		env.Print()
		if i != len(environments)-1 {
			doc.Divider(80)
		}
	}
}
