ORAS x Harbor - Push & Pull Arbitrary Files as OCI Artifacts

This repository demonstrates how to push and pull arbitrary files (such as .txt, .yaml, .json) to a Harbor OCI registry using ORAS Go SDK, without wrapping them as Docker container images.

It follows the OCI Artifacts specification — enabling clean, type-safe usage of registries beyond container use cases.

---

What is an OCI Artifact?

OCI Artifacts allow storing any file type (not just container images) in an OCI-compliant registry like Harbor. Instead of forcing data into fake image layers, each artifact is declared with:

- A config blob with custom mediaType
- One or more layer blobs (your file)
- An OCI Manifest describing it all

This lets registries like Harbor manage files like Helm charts, SBOMs, policies, templates, configs, and more — the right way.

---

Project Structure

.
├── main.go Entry point with CLI flags for push/pull
├── oras/
│ ├── push.go PushFileToOCI(): Push file as OCI artifact
│ └── pull.go PullFromOCI(): Download artifact to local path
├── template/
│ └── hello.txt Sample file to push
├── download/ Target folder for pulled artifacts
└── go.mod Go module definition

---

How to Use

1. Start Harbor on localhost:8080

Make sure your Harbor instance is running at http://localhost:8080
Create a project named "demo".

If using plain HTTP, edit Docker daemon config:
{
"insecure-registries": ["localhost:8080"]
}
Then restart Docker.

---

2. Run the CLI

Push File

go run main.go --mode=push

This pushes ./template/hello.txt to:

localhost:8080/demo/template:latest

Pull File

go run main.go --mode=pull

This pulls the artifact from Harbor and saves it to ./download/

---

Authentication

Your Harbor credentials are specified in HarborConfig:

HarborConfig{
URL: "localhost:8080",
Repo: "demo/template",
Username: "admin",
Password: "Harbor12345",
}

Modify these as needed to match your Harbor setup.

---

Technical Highlights

- Uses "application/vnd.test.file" as the custom mediaType
- No fake Docker image needed — fully OCI-compliant
- Uses oras.Copy() from ORAS Go SDK for push & pull
- Enables secure & structured artifact handling in Harbor

---

References

- ORAS Go SDK: https://github.com/oras-project/oras-go
- OCI Artifacts Spec: https://github.com/opencontainers/artifacts
- Harbor: https://goharbor.io

---

Author

Developed by miha for educational and prototyping purposes.
Feel free to fork and adapt!
