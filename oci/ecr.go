package oci

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

type ecrPuller struct {
	credentialRetriever CredentialRetriever
}

// Pull implements Puller.
func (e *ecrPuller) Pull(ctx context.Context, id ArtifactID, path string) error {
	fs, err := file.New(path)
	if err != nil {
		return err
	}
	defer fs.Close()
	// Create a remote repository client
	repository, err := remote.NewRepository(fmt.Sprintf("%s/%s", id.Registry, id.Repo))
	if err != nil {
		return err
	}
	creds, err := e.credentialRetriever.RetrieveCredential(ctx)
	if err != nil {
		return err
	}
	repository.Client = &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.DefaultCache,
		Credential: auth.StaticCredential(id.Registry, creds),
	}

	// Copy from the remote repository to the OCI layout store
	manifestDescriptor, err := oras.Copy(ctx, repository, id.Tag, fs, id.Tag, oras.DefaultCopyOptions)
	if err != nil {
		return err
	}
	log.Info().Interface("manifest", manifestDescriptor).Msgf("pulled image")
	return nil
}

func NewECRPuller(credentialRetriever CredentialRetriever) Puller {
	return &ecrPuller{credentialRetriever}
}
