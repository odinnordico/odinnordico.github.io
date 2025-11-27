# Installation

This guide covers the different ways to install odinnordico.github.io on your system.

## Prerequisites

Before installing odinnordico.github.io, ensure you have the following installed:

*   **Go**: Version 1.21 or higher is required. You can download it from [go.dev](https://go.dev/dl/).
*   **Git**: Required for cloning the repository. Download from [git-scm.com](https://git-scm.com/downloads).

## Installation Methods

### Option 1: Install via `go install` (Recommended)

If you have Go installed and configured in your path, this is the easiest way to get the latest version.

```bash
go install github.com/odinnordico/odinnordico.github.io@latest
```

Ensure your `$(go env GOPATH)/bin` is in your system's `PATH`.

### Option 2: Build from Source

If you want to modify the code or contribute, building from source is the best option.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/odinnordico/odinnordico.github.io.git
    cd odinnordico.github.io
    ```

2.  **Build the binary:**

    ```bash
    go build -o odinnordico.github.io .
    ```

3.  **Move to a directory in your PATH (Optional):**

    ```bash
    sudo mv odinnordico.github.io /usr/local/bin/
    ```

### Option 3: Download Pre-built Binaries

*(Note: If you have set up a release pipeline, you would link to the Releases page here. For now, we assume users build from source or use `go install`)*.

Check the [Releases](https://github.com/odinnordico/odinnordico.github.io/releases) page for pre-built binaries for your operating system (Windows, macOS, Linux).

## Verifying Installation

To verify that odinnordico.github.io is installed correctly, run:

```bash
odinnordico.github.io --help
```

You should see the help message with a list of available commands.

```text
A powerful resume generator for developers.

Usage:
  odinnordico.github.io [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pdf         Generate PDF resume
  serve       Start development server
  website     Generate static website

Flags:
      --config string       config file (default is $HOME/.odinnordico.github.io.yaml)
      --data-dir string     Directory containing data files (default "data")
  -h, --help                help for odinnordico.github.io
      --output-dir string   Directory for generated output (default "public")
```
