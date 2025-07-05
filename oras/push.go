package oras

import (
	"context"
	"fmt"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

type HarborConfig struct {
	URL      string
	Repo     string
	Username string
	Password string
}

func PushFileToOCI(config HarborConfig, filePath, tag string) error {
	// 0. Create a file store 
	fs, err := file.New("./")
	if err != nil {
		panic(err)
	}
	defer fs.Close()
	ctx := context.Background()

	// 1. Add files to the file store
	mediaType := "application/vnd.test.file" 
	fileNames := []string{filePath} // List of files to be pushed
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))

	for _, name := range fileNames {
		fileDescriptor, err := fs.Add(ctx, name, mediaType, "")
		if err != nil {
			panic(err)
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
		fmt.Printf("file descriptor for %s: %v\n", name, fileDescriptor)
	}

	// 2. Pack the files and tag the packed manifest
  artifactType := "application/vnd.test.artifact"
	opts := oras.PackManifestOptions{
		Layers: fileDescriptors,
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1, artifactType, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("manifest descriptor:", manifestDescriptor)

	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		panic(err)
	}

	// 3. Connect to a remote repository
	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", config.URL, config.Repo))
	if err != nil {
		panic(err)
	}
	repo.PlainHTTP = true
	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(config.URL, auth.Credential{
			Username: config.Username,
			Password: config.Password,
		}),
	}

	// 4. Push the manifest to the remote repository 
	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(fmt.Errorf("failed to push file: %w", err))
	}

	return nil
}

