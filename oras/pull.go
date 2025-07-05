package oras

import (
	"context"
	"fmt"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)


// PullFromOCI pulls an artifact with the given tag from the OCI registry into the target local path
func PullFromOCI(config HarborConfig, tag, localPath string) error {
	ctx := context.Background()

	// 1. Create a local file store
	fs, err := file.New(localPath)
	if err != nil {
		return fmt.Errorf("failed to create file store: %w", err)
	}
	defer fs.Close()

	// 2. Create the remote repository
	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", config.URL, config.Repo))
	if err != nil {
		return fmt.Errorf("failed to create remote repository: %w", err)
	}
	repo.PlainHTTP = true // Disable HTTPS for local testing

	// 3. Add authentication
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(config.URL, auth.Credential{
			Username: config.Username,
			Password: config.Password,
		}),
	}

	// 4. Pull artifact
	desc, err := oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
	if err != nil {
		return fmt.Errorf("failed to pull artifact: %w", err)
	}

	fmt.Println("✅ Pulled manifest descriptor:", desc)
	return nil
}
