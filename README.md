# Podinfo Token Validator

## Overview

**Performs the following steps:**

1. Call the POST `/token` endpoint of the Podinfo app to generate a token.
2. Validate the generated token using the GET `/token/validate` endpoint.
3. Store the validation result in the Podinfo app's cache using the POST/PUT `/cache/validation_result` endpoint.
4. Connect to Redis to verify that the cache key exists and contains the correct validation result.

## Development and Debugging

### Prerequisites

Managing dependencies is done using NixOS+Devbox. While Devbox is recommended for an efficient setup, it is not mandatory. If you choose to proceed without Devbox, ensure that you manually fit the dependencies listed in `devbox.json`.

* Docker/Podman
* The project uses a `Taskfile` to streamline development workflow. `task` is a modern alternative to `Make`. Install `task` to execute predefined commands, such as building the app, tests, and deploying to kubernetes.
* This project relies on the **Podinfo app** being pre-installed in the Kubernetes cluster. Ensure that the Podinfo service is running and accessible before deploying or running this application.
* Access to the `bringg-local` (Kind) Kubernetes cluster. Devbox ensures the correct `KUBECONFIG` is exported for this purpose.

### Configuration

The app accepts both: config file and environment variables.

#### Configuration file

Update the `config.yaml` file to define Podinfo and Redis connection details:

```yaml
log:
  level: debug
  format: json

podinfo:
  base_url: "http://localhost:29898"
  token_endpoint: "/token"
  token_validate: "/token/validate"
  cache_endpoint: "/cache/validation_result"

redis:
  host: "localhost"
  port: 26379
  validation_result_key: "validation_result"
```

#### Environment Variables

Using environment variables:

* `LOG_LEVEL`
* `LOG_FORMAT`
* `PODINFO_BASE_URL`
* `REDIS_HOST`
* `REDIS_PORT`

### Development and Debugging

#### Using Tasks

This project uses a Taskfile to automate common workflows. Here’s how you can use it:

* **Key Commands:**

   * `task test`: Runs the application tests.
   * `task run`: Starts the application locally.
   * `task local-release`: Starts test,docker build and push, update manifests. Then you will have to commit changes, argocd will pick up the new version.
   * `task -a` - list all available tasks and their descriptions.

* **Version Management:** The `Taskfile` includes a `VERSION` variable. If you plan to release a new version, update this variable to the desired version number.
* **Container Deployment:**

   * Use `task docker-push` to push the Docker image to a registry.
   * Alternatively, use `task kind-load` to load the image directly into `bringg-local` (kind) cluster.

* **Podinfo Integration:**

   * Use `task podinfo-port-fwd` to set up port forwarding to podinfo backend and redis.

* Use `task podinfo-redis-cli` to connect podinfo-redis using redis-cli.

