package oci

import "context"

type Puller interface {
	Pull(ctx context.Context, id ArtifactID, path string) error
}

type ArtifactID struct {
	Registry string
	Repo     string
	Tag      string
}

type Reader interface {
	ListDirs(ctx context.Context, filePath string) ([]string, error)
	FilterModules(ctx context.Context, dirs []string) []string
}
