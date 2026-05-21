package domain

type Template struct {
	Name    string
	Columns []string
}

var Templates = map[string]Template{
	"basic-kanban": {
		Name:    "Basic Kanban",
		Columns: []string{"To Do", "In Progress", "Done"},
	},
	"bug-tracker": {
		Name:    "Bug Tracker",
		Columns: []string{"Backlog", "Investigating", "Fixing", "Verified"},
	},
	"content-pipeline": {
		Name:    "Content Pipeline",
		Columns: []string{"Ideas", "Drafting", "Review", "Published"},
	},
}
