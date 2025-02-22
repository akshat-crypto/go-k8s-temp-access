# kube-temp-access

A CLI tool to create temporary Kubernetes access for dashboard login or auditing purposes. This tool generates a service account, assigns configurable permissions (e.g., `view`, `deployments`, `pods`), provides a token, and auto-deletes resources after a specified duration.

## Features
- Supports multiple namespaces (e.g., `default`, `kube-system`).
- Configurable permissions via the `--resources` flag (e.g., `view`, `deployments`, `pods`, `statefulsets`).
- Adjustable token expiration and cleanup time with the `--expiration` flag (default: 15 minutes).
- Automatic resource cleanup using Kubernetes Jobs.

## Prerequisites
- Go 1.18 or later installed.
- A valid kubeconfig file with cluster-admin permissions (required for creating `ServiceAccounts`, `ClusterRoles`, `ClusterRoleBindings`, and `Jobs`).

## Project Structure

golang-k8s-temp-access/
├── cmd/
│   └── root.go           # Cobra CLI root and create command
├── internal/
│   ├── k8s/
│   │   ├── client.go     # Kubernetes client setup
│   │   ├── resources.go  # Resource creation and deletion logic
│   │   └── token.go      # Token generation
│   └── utils/
│       └── uuid.go       # UUID utility
├── main.go               # Entry point
├── go.mod                # Go module file
└── go.sum                # Dependency checksums


## Build Instructions
1. **Set Up the Repository:**
   ```bash
   mkdir -p ~/Desktop/work/golang-post-project-1/golang-k8sconfig-creator
   cd ~/Desktop/work/golang-post-project-1/golang-k8sconfig-creator
   go mod init golang-k8s-temp-access
    ```
2. **Install Dependencies:**
    ```bash
    go get github.com/spf13/cobra@v1.8.0
    go get k8s.io/client-go@v0.28.0
    go get github.com/google/uuid
    go mod tidy
    ```
3. **Build the Binary (CGO Enabled):**
    ```bash
    CGO_ENABLED=1 go build -o kube-temp-access
    ```

## Usage

- Run the create command with your kubeconfig and desired options:

```bash
    ./kube-temp-access create --kubeconfig ~/.kube/k8s-local --namespace default,kube-system --resources deployments,pods --expiration 30m
```
- Shorthand: 

```bash
./kube-temp-access create -k ~/.kube/k8s-local -n default,kube-system -r deployments,pods -e 30m
```

- Flags:

```bash
-k, --kubeconfig string: Path to kubeconfig file (default: ~/.kube/config).
-n, --namespace strings: Comma-separated list of namespaces (default: default).
-r, --resources strings: Comma-separated list of resources to grant access to (e.g., view, deployments, pods, statefulsets; default: view).
-e, --expiration duration: Token expiration and cleanup duration (e.g., 15m, 1h; default: 15m).
```

## Example Output

Here’s an example run:

```bash

./kube-temp-access create --kubeconfig ~/.kube/k8s-local --namespace default,kube-system

Created service account temp-sa-ddbfd88c-bc2e-4422-bee9-aa9577c6fdc7 in namespace default
Created cluster role binding temp-binding-ddbfd88c-bc2e-4422-bee9-aa9577c6fdc7
Token for dashboard login in namespace default: eyJhbGciOiJSUzI1NiIsImtpZCI6ImN5b3NENzdfZ05oZlFnS00xanBWTlI1YkIxY3l4ZVFrcGFUMnlTRmV2RVkifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzQwMjM0OTM0LCJpYXQiOjE3NDAyMzEzMzQsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiOTM5ZTJjNGItZDg1OC00Zjk2LWEzZGEtOTg5NmMwODdkNjFkIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6InRlbXAtc2EtZGRiZmQ4OGMtYmMyZS00NDIyLWJlZTktYWE5NTc3YzZmZGM3IiwidWlkIjoiOGJlMTViMDYtM2U0Ni00MzkwLTliMWUtODZhMmEwYTdlODNlIn19LCJuYmYiOjE3NDAyMzEzMzQsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnRlbXAtc2EtZGRiZmQ4OGMtYmMyZS00NDIyLWJlZTktYWE5NTc3YzZmZGM3In0.EZbhXQF9XhEY7JvpprCLxHxGvWirqZ-5GkoKxnvu4SMDEMTZMCoF_A_sdaLliqmrf1lfKQ0G0Ub5PyZJW4xGeezN0s4CAT6WGGT_Oi23jYVYdJ04Agk1kHVepUmQZaMMLdpTGuLYkmXZx0DBfTIFpxifJ0BqrpXhuD-IMk4ECZjQWOPr2-YTq02xUzPSTsaS-6g7Y98D1PkWsA7Xt69eqLMFa4VNHiF5_uNsA6Yj6cNTa6y1_d_9j0SCkyIOpKMaiQwdFp2hXMMHqFgCK6zHQuo2Vx5Ji2qGstWZcBre8O_LLQJZP0BTK_TLPqDJlR_jkFg--Ilvfg0hfhAHE4fX7Q
Created deletion job temp-deletion-job-ddbfd88c-bc2e-4422-bee9-aa9577c6fdc7 in namespace default, resources will be deleted in 15 minutes
Created service account temp-sa-adfd0336-ccb0-4a6b-ba46-aa1f797f5536 in namespace kube-system
Created cluster role binding temp-binding-adfd0336-ccb0-4a6b-ba46-aa1f797f5536
Token for dashboard login in namespace kube-system: eyJhbGciOiJSUzI1NiIsImtpZCI6ImN5b3NENzdfZ05oZlFnS00xanBWTlI1YkIxY3l4ZVFrcGFUMnlTRmV2RVkifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzQwMjM0OTM1LCJpYXQiOjE3NDAyMzEzMzUsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiZmFlNzQzY2YtOTIzMC00ZjE3LTk4MDMtYWU3MGY5OTQ4ZDQ0Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsInNlcnZpY2VhY2NvdW50Ijp7Im5hbWUiOiJ0ZW1wLXNhLWFkZmQwMzM2LWNjYjAtNGE2Yi1iYTQ2LWFhMWY3OTdmNTUzNiIsInVpZCI6IjhjMmViZDQxLWJkNWMtNDEyNS1hOTllLWQ1MzFhOGZlM2MzNCJ9fSwibmJmIjoxNzQwMjMxMzM1LCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06dGVtcC1zYS1hZGZkMDMzNi1jY2IwLTRhNmItYmE0Ni1hYTFmNzk3ZjU1MzYifQ.Fb7FxrOl-WQC20E3Ez5_sPvQ06KhWqDE0HbeI9ZO5mbRBbDB3Y6mGvaErkqHQhHVwHt7QAZzmJeYucFqg54NCcl5wcwYUIq4joJNB1xipUV-_FqVYQc3_3gq3qkn-6LMgfaJMMIdMoANG_lnXJ4I_OelfIbIqr8H-54_VxySDZ17VGPY3Pf48UCQ77c6sdeqd86cPoUKyoeHBwZYMRGcwxjuW3Yb0Zp_gQPQkLBTkqfp4b4LaFyLLMZLc6UuZuKz0J_3y0cEWegVymHhn7jvrKyPY8lxKtHTwiRw1-nCrbVTej-PWi1WcpFlZVvI9esaX4ZdRm4m1jF16U55o1x-0g
Created deletion job temp-deletion-job-adfd0336-ccb0-4a6b-ba46-aa1f797f5536 in namespace kube-system, resources will be deleted in 15 minutes

```


### Notes

- Ensure your kubeconfig has sufficient permissions.
- The deletion job uses bitnami/kubectl:latest; adjust if your cluster requires a specific version.
- For production, consider adding error handling and logging.

