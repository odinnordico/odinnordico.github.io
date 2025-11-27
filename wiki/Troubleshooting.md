# Troubleshooting

Common issues and their solutions.

## PDF Generation

### Images not appearing in PDF
*   **Cause**: Incorrect path or unsupported format.
*   **Solution**:
    *   Ensure the image path in your YAML or template is correct.
    *   Use PNG or JPG formats.
    *   Check if the file exists in `assets/`.

### Text overflow or overlapping
*   **Cause**: Text is too long for the allocated column width or height.
*   **Solution**:
    *   Increase the `height` of the row in `resume.yaml.tmpl`.
    *   Decrease the font size.
    *   Shorten the text content.

## Website Generation

### "Template file not found"
*   **Cause**: The specified theme does not exist or is missing files.
*   **Solution**:
    *   Check the `--theme` flag.
    *   Ensure `templates/<theme>/index.html.tmpl` exists.

### Changes not reflecting in browser
*   **Cause**: Browser caching.
*   **Solution**:
    *   Hard refresh the page (Ctrl+F5 or Cmd+Shift+R).
    *   Disable cache in browser developer tools.

## General

### "command not found: odinnordico.github.io"
*   **Cause**: The binary is not in your system PATH.
*   **Solution**:
    *   If installed via `go install`, ensure `$(go env GOPATH)/bin` is in your PATH.
    *   If built from source, move the binary to `/usr/local/bin` or add the build directory to PATH.

### YAML Parsing Errors
*   **Cause**: Invalid YAML syntax (indentation is often the culprit).
*   **Solution**:
    *   Use a YAML validator or linter.
    *   Check that indentation uses spaces, not tabs.
