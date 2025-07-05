package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/baominh/go-oras/oras"
)

func main() {
	mode := flag.String("mode", "push", "Push mode (push/pull)")
	flag.Parse()

	config := oras.HarborConfig{
		URL:      "localhost:8080",
		Repo:     "demo/template",
		Username: "admin",
		Password: "Harbor12345",
	}

	// File path to push
	filePath := "./template/hello.zip"
	// Local folder path to store pulled artifact
	outputPath := "./download"

	switch *mode {
	case "push":
		fmt.Println("Pushing file to Harbor...")
		err := oras.PushFileToOCI(config, filePath, "latest")
		if err != nil {
			log.Fatalf("❌ Failed to push file: %v", err)
		}
		fmt.Println("✅ File pushed successfully.")

	case "pull":
		fmt.Println("Pulling file from Harbor...")
		err := oras.PullFromOCI(config, "latest", outputPath)
		if err != nil {
			log.Fatalf("❌ Failed to pull file: %v", err)
		}
		fmt.Println("✅ File pulled successfully.")

	default:
		log.Fatalf("Invalid mode: %s. Use 'push' or 'pull'.", *mode)
	}
}

