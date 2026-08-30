# Project Overview

This is a Go-based metrics and alerting service. It consists of two main components: an `agent` that collects metrics and a `server` that receives and stores them. The agent periodically collects system metrics and sends them to the server. The server exposes an HTTP API to receive and query metrics. The project uses `chi` for routing, `zerolog` for logging, and an in-memory storage (`MemStorage`) for the metrics.

# Building and Running

The `Makefile` provides the necessary commands for building and running the project.

**Building the project:**

*   `make all`: Build both the agent and the server.
*   `make server`: Build the server. The binary will be located at `cmd/server/server`.
*   `make agent`: Build the agent. The binary will be located at `cmd/agent/agent`.

**Running the project:**

1.  Start the server:
    ```bash
    ./cmd/server/server
    ```
2.  Start the agent:
    ```bash
    ./cmd/agent/agent
    ```

**Running tests:**

*   `make tests`: Run the integration tests located in the `tests/` directory.

# Development Conventions

*   **Project Structure:** The project follows the standard Go project layout.
*   **Dependencies:** The project uses `go mod` for dependency management. Key libraries include:
    *   `github.com/go-chi/chi/v5` for HTTP routing.
    *   `github.com/rs/zerolog` for structured logging.
    *   `github.com/stretchr/testify` for assertions in tests.
*   **Testing:** The project has a suite of integration tests in the `tests/` directory, which can be run using the `make tests` command.
*   **Configuration:** The agent and server can be configured using command-line flags.
