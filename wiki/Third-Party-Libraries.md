# Third Party Libraries

odinnordico.github.io relies on several powerful open-source libraries. Understanding these libraries can help you when customizing or extending odinnordico.github.io.

## gofpdf (PDF Generation)

[gofpdf](https://github.com/grafana/gofpdf) is a library for generating PDF documents with high level support for text, drawing and images.

*   **Usage**: It is used in `internal/generator/pdf.go` to render the resume PDF.
*   **Customization**: The `resume.yaml.tmpl` file maps to a grid system implemented using gofpdf primitives.
*   **Documentation**: [gofpdf Docs](https://pkg.go.dev/github.com/grafana/gofpdf)

## Scriggo (Template Engine)

[Scriggo](https://github.com/open2b/scriggo) is a template engine for Go that looks like `html/template` but is more powerful and allows for more complex logic.

*   **Usage**: It is used in `internal/generator/website.go` to render the HTML website.
*   **Customization**: You can use standard Go template syntax in `index.html.tmpl`.
*   **Documentation**: [Scriggo Docs](https://scriggo.com/)

## Cobra (CLI Framework)

[Cobra](https://github.com/spf13/cobra) is a library for creating powerful modern CLI applications.

*   **Usage**: It defines the command structure (`odinnordico.github.io pdf`, `odinnordico.github.io website`, etc.) in the `cmd/` directory.
*   **Customization**: If you want to add a new command, you would use Cobra.
*   **Documentation**: [Cobra Docs](https://cobra.dev/)

## Viper (Configuration)

[Viper](https://github.com/spf13/viper) is a complete configuration solution for Go applications.

*   **Usage**: It handles reading configuration from flags, config files, and environment variables.
*   **Customization**: You can add new configuration options in `cmd/root.go`.
*   **Documentation**: [Viper Docs](https://github.com/spf13/viper)

## fsnotify (File Watching)

[fsnotify](https://github.com/fsnotify/fsnotify) provides cross-platform file system notifications.

*   **Usage**: It is used by the `serve` command to detect changes in data or template files and trigger a regeneration.
*   **Documentation**: [fsnotify Docs](https://github.com/fsnotify/fsnotify)
