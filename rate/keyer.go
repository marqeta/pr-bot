package rate

import "github.com/marqeta/pr-bot/id"

type Keyer func(id id.PR) string

func OrgKey(id id.PR) string {
	return "Org/" + id.Owner
}

func RepoKey(id id.PR) string {
	return "Repo/" + id.RepoFullName
}

func AuthorKey(id id.PR) string {
	return "Author/" + id.Author
}
