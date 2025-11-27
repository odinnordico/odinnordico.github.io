# Development Guide

This guide is for developers who want to contribute to odinnordico.github.io or understand its internal architecture.

## Project Architecture

odinnordico.github.io follows a standard Go project layout:

*   `cmd/`: Contains the main entry points for the CLI commands (using Cobra).
*   `internal/`: Private application code.
    *   `generator/`: Core logic for PDF and Website generation.
    *   `loader/`: Logic for loading and parsing YAML data files.
    *   `models/`: Go structs representing the resume data.
    *   `utils/`: Helper functions.
*   `templates/`: Default templates embedded in the binary.

## Key Libraries

*   **[Cobra](https://github.com/spf13/cobra)**: Used for CLI command structure and flag parsing.
*   **[Viper](https://github.com/spf13/viper)**: Handles configuration loading.
*   **[Maroto](https://github.com/johnfercher/maroto)**: The engine behind PDF generation.
*   **[Scriggo](https://github.com/open2b/scriggo)**: Template engine for HTML generation.
*   **[fsnotify](https://github.com/fsnotify/fsnotify)**: Used for the live reload feature.

## Building and Testing

### Prerequisites

*   Go 1.21+
*   Make (optional, if a Makefile exists)

### Running Tests

Run all tests in the project:

```bash
go test ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

### Building

Build the binary:

```bash
go build -o odinnordico.github.io .
```

## Contributing

1.  **Fork** the repository on GitHub.
2.  **Clone** your fork locally.
3.  Create a **feature branch** (`git checkout -b feature/my-new-feature`).
4.  **Commit** your changes.
5.  **Push** to the branch.
6.  Open a **Pull Request**.

Please ensure your code passes `go fmt` and `go vet` before submitting.
