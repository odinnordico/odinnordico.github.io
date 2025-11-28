# Templates and Customization

odinnordico.github.io allows you to completely customize the look and feel of your resume through its theming system.

## Themes

Themes are located in the `templates/` directory. The default theme is `default`.

To create a new theme:

1.  Create a new directory: `templates/mytheme/`
2.  Copy the default templates:
    ```bash
    cp templates/default/resume.yaml.tmpl templates/mytheme/
    cp templates/default/index.html.tmpl templates/mytheme/
    ```
3.  Use your theme: `odinnordico.github.io pdf --theme mytheme`

## PDF Templates (`resume.yaml.tmpl`)

The PDF generation is driven by a YAML template that defines the structure and layout. This template is processed by Go's `text/template` engine before being parsed by the PDF generator (powered by `gofpdf`).

### Structure

The template is a list of rows, where each row contains columns.

```yaml
rows:
  - height: 10
    cols:
      - width: 12
        text:
          content: "{{.Basic.Name}}"
          size: 24
          style: bold
          align: center
```

### Components

*   **Text**: Renders text content.
    *   `content`: The text string (supports template variables).
    *   `size`: Font size.
    *   `style`: `bold`, `italic`, or `normal`.
    *   `align`: `left`, `center`, or `right`.
    *   `color`: RGB color object (`red`, `green`, `blue`).
    *   `hyperlink`: URL to open when clicked.

*   **Image**: Renders an image.
    *   `path`: Absolute path to the image file.
    *   `width`: Width percentage (0-100).
    *   `center`: Boolean to center the image.

*   **Line**: Renders a horizontal line.
    *   `thickness`: Line thickness.
    *   `color`: RGB color.

### Template Functions

odinnordico.github.io provides several helper functions for use in templates:

*   `formatDate`: Formats a date string (e.g., `{{ formatDate .Date }}`).
*   `formatEndDate`: Formats an end date, handling `null` as "Present".
*   `formatYear`: Extracts the year from a date.
*   `hasSocials`: Checks if the social list is not empty.
*   `assetPath`: Resolves the absolute path to an asset.
*   `calculateHeight`: Estimates the height required for a block of text.

## Website Templates (`index.html.tmpl`)

The website template is a standard HTML file that uses the [Scriggo](https://github.com/open2b/scriggo) template engine (which is very similar to Go's `html/template`).

You have full control over the HTML and CSS. You can include scripts, styles, and any other HTML elements.

### Variables

The entire data structure is available in the template via the `.` (dot) variable.

*   `{{.Basic}}`: Basic info.
*   `{{.Professional}}`: Work experience.
*   `{{.Education}}`: Education.
*   ...and so on.

### Example

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Basic.Name}} - Resume</title>
</head>
<body>
    <h1>{{.Basic.Name}}</h1>
    <p>{{.Basic.Summary}}</p>
    
    <h2>Experience</h2>
    {{range .Professional.Jobs}}
        <div class="job">
            <h3>{{.Position}} at {{.Company.Name}}</h3>
            <p>{{formatDate .StartDate}} - {{formatEndDate .EndDate}}</p>
        </div>
    {{end}}
</body>
</html>
```
