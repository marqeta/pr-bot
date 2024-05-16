package oci

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type bundleFileReader struct {
}

func (r *bundleFileReader) ListDirs(_ context.Context, localfilepath string) ([]string, error) {
	file, err := os.Open(localfilepath)
	if err != nil {
		log.Info().Msgf("Error opening file %s", localfilepath)
		return nil, err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		log.Info().Msgf("Error creating gzip reader for file %s", localfilepath)
		return nil, err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	modulesMap := make(map[string]bool)
	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		if header == nil {
			continue
		}

		switch header.Typeflag {

		case tar.TypeReg:
			// /abc/def/ghi.rego => /abc/def
			dir := filepath.Dir(header.Name)
			modulesMap[dir] = true
		}
	}

	modules := make([]string, 0)
	for module := range modulesMap {
		modules = append(modules, module)
	}

	return modules, nil
}

func (r *bundleFileReader) FilterModules(_ context.Context, dirs []string) []string {
	modules := make([]string, 0)
	for _, dir := range dirs {
		if strings.Count(dir, string(os.PathSeparator)) != 2 {
			// skip files in the root directory and files in depth > 2
			continue
		}

		nodes := strings.Split(dir, string(os.PathSeparator))
		hidden := false
		for _, node := range nodes {
			if strings.HasPrefix(node, "_") {
				hidden = true
				break
			}
		}
		if hidden {
			// dirs that start with _ are used for hosting library functions and utility functions
			continue
		}
		modules = append(modules, dir)
	}

	return modules

}
func NewReader() Reader {
	return &bundleFileReader{}
}
