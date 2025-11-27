# Multi-language Support

odinnordico.github.io is designed to support multiple languages out of the box. You can generate your resume in as many languages as you need from a single project.

## Directory Structure

The default language (usually English) resides in the root of the `data/` directory. Additional languages are stored in `data/lang/<language_code>/`.

```text
data/
├── basic.yml          # Default language (e.g., English)
├── professional.yml
└── lang/
    ├── es/            # Spanish
    │   ├── basic.yml
    │   └── professional.yml
    └── fr/            # French
        ├── basic.yml
        └── professional.yml
```

## Adding a New Language

1.  **Create the directory**:
    Create a new directory inside `data/lang/` with your language code (e.g., `de` for German).
    ```bash
    mkdir -p data/lang/de
    ```

2.  **Copy data files**:
    Copy your existing YAML files to the new directory.
    ```bash
    cp data/*.yml data/lang/de/
    ```

3.  **Translate**:
    Edit the YAML files in the new directory and translate the content.

## Generating Output

### PDF

When you run `odinnordico.github.io pdf`, it automatically detects all available languages and generates a PDF for each one.

*   Default language: `public/assets/files/resume.pdf`
*   Spanish: `public/assets/files/resume-es.pdf`
*   French: `public/assets/files/resume-fr.pdf`

You can also generate for a specific language:

```bash
odinnordico.github.io pdf --lang es
```

### Website

When you run `odinnordico.github.io website`, it generates a complete site structure:

*   Default language: `public/index.html`
*   Spanish: `public/es/index.html`
*   French: `public/fr/index.html`

The website template usually includes links to switch between languages.
