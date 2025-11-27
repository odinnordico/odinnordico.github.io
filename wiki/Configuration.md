# Configuration

odinnordico.github.io can be configured using command-line flags or a configuration file.

## Command Line Flags

The following global flags are available for all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to config file | `$HOME/.odinnordico.github.io.yaml` |
| `--data-dir` | Directory containing data files | `data` |
| `--output-dir` | Directory for generated output | `public` |

### Command-Specific Flags

#### `pdf`
| Flag | Description | Default |
|------|-------------|---------|
| `--theme` | Theme name to use | `default` |

#### `website`
| Flag | Description | Default |
|------|-------------|---------|
| `--theme` | Theme name to use | `default` |

#### `serve`
| Flag | Description | Default |
|------|-------------|---------|
| `--host` | Host to serve on | `localhost` |
| `--port` | Port to serve on | `8080` |
| `--watch` | Enable live reload | `false` |
| `--theme` | Theme name to use | `default` |

## Configuration File

You can also configure odinnordico.github.io using a YAML configuration file. By default, odinnordico.github.io looks for `.odinnordico.github.io.yaml` in your home directory.

Example `.odinnordico.github.io.yaml`:

```yaml
data-dir: "/path/to/my/resume/data"
output-dir: "/path/to/my/resume/public"
theme: "modern"
```

### Precedence

Configuration is applied in the following order (highest precedence first):

1.  Command-line flags
2.  Configuration file
3.  Default values
