# odinnordico.github.io

A powerful, flexible resume generator that creates beautiful PDFs and static websites from YAML data files. Built with Go and designed for developers who want version-controlled, multi-language resumes.

## Features

- 📄 **PDF Generation** - Create professional PDF resumes using the [Maroto v2](https://github.com/johnfercher/maroto) library
- 🌐 **Static Website** - Generate a responsive HTML website from the same data
- 🌍 **Multi-language Support** - Easily maintain resumes in multiple languages
- 🎨 **Themeable** - Customize the look and feel with templates
- 🔄 **Live Reload** - Development server with automatic regeneration on file changes
- 📝 **YAML-based** - Simple, readable data format
- 🎯 **Type-safe** - Strongly typed Go models ensure data consistency

## Technology Stack

- **PDF Generation**: [Maroto v2](https://github.com/johnfercher/maroto) - Pure Go PDF library
- **Template Engine**: [Scriggo](https://github.com/open2b/scriggo) - Fast Go template engine
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) - Modern CLI framework
- **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management
- **File Watching**: [fsnotify](https://github.com/fsnotify/fsnotify) - Cross-platform file system notifications

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/odinnordico/odinnordico.github.io.git
cd odinnordico.github.io

# Build the binary
go build -o odinnordico.github.io .

# Or install globally
go install
```

## Quick Start

### 1. Set Up Your Data

The resume data is stored in YAML files in the `data/` directory:

```
data/
├── basic.yml          # Personal information
├── professional.yml   # Work experience
├── education.yml      # Education history
├── certificates.yml   # Certifications
├── skills.yml        # Skills and expertise
├── social.yml        # Social media links
└── lang/             # Translations
    └── es/           # Spanish translations
        ├── basic.yml
        ├── professional.yml
        └── ...
```

### 2. Generate Your Resume

#### Generate PDF

```bash
# Generate PDF for all languages
go run . pdf

# Generate PDF for specific language
go run . pdf --lang es

# Use a custom theme
go run . pdf --theme modern
```

Output: `public/assets/files/resume.pdf` (and `resume-es.pdf` for Spanish)

#### Generate Website

```bash
# Generate static website for all languages
go run . website

# Use a custom theme
go run . website --theme modern
```

Output: `public/index.html` (English) and `public/es/index.html` (Spanish)

#### Development Server

```bash
# Start server (default: http://localhost:8080)
go run . serve

# Start with live reload (auto-regenerates on file changes)
go run . serve --watch

# Custom host and port
go run . serve --host 0.0.0.0 --port 3000 --watch
```

## Data Structure

### Basic Information (`data/basic.yml`)

```yaml
name: "John Doe"
display_name: "John Doe"
location: "San Francisco, CA"
summary: |
  Experienced software engineer with expertise in...
website: "https://johndoe.com"
```

### Professional Experience (`data/professional.yml`)

```yaml
title: "Senior Software Engineer"
years_of_experience: 8
jobs:
  - position: "Senior Software Engineer"
    company:
      name: "Tech Corp"
      url: "https://techcorp.com"
      logo:
        library: "brands"
        image: "company-logo"
    start_date: 2020-01-15
    end_date: null  # null means current position
    job_description: |
      - Led development of microservices architecture
      - Mentored junior developers
      - Improved system performance by 40%
```

### Social Links (`data/social.yml`)

```yaml
- name: "Email"
  url: "mailto:john@example.com"
  logo:
    library: "solid"
    image: "envelope"

- name: "GitHub"
  url: "https://github.com/johndoe"
  logo:
    library: "brands"
    image: "github"

- name: "LinkedIn"
  url: "https://linkedin.com/in/johndoe"
  logo:
    library: "brands"
    image: "linkedin"
```

### Skills (`data/skills.yml`)

```yaml
- name: "Go"
  description: "Backend development, microservices"
  level: 9
  logo:
    library: "devicon"
    image: "go"
  tags:
    - "backend"
    - "systems"

- name: "Python"
  description: "Data processing, automation"
  level: 8
  tags:
    - "backend"
    - "scripting"
```

### Education (`data/education.yml`)

```yaml
- title: "Bachelor of Science in Computer Science"
  level: "Bachelor's Degree"
  provider:
    name: "University of Technology"
    url: "https://university.edu"
  date: 2015-05-20
  description: "Graduated with honors"
```

### Certificates (`data/certificates.yml`)

```yaml
- name: "AWS Certified Solutions Architect"
  description: "Professional level certification"
  date: 2023-06-15
  certificate_url: "https://aws.amazon.com/certification/"
  url: "https://verify.cert.com/12345"
  provider:
    name: "Amazon Web Services"
    url: "https://aws.amazon.com"
  topics:
    - "Cloud Architecture"
    - "AWS Services"
```

## Multi-language Support

### Adding a New Language

1. Create a language directory:
```bash
mkdir -p data/lang/fr
```

2. Copy and translate your YAML files:
```bash
cp data/basic.yml data/lang/fr/basic.yml
# Edit data/lang/fr/basic.yml with French translations
```

3. Generate outputs:
```bash
# Generates resume.pdf, resume-es.pdf, resume-fr.pdf
go run . pdf

# Generates index.html, es/index.html, fr/index.html
go run . website
```

### Language Detection

The system automatically:
- Detects available languages from `data/lang/` directories
- Generates PDFs with language suffixes (e.g., `resume-es.pdf`)
- Creates language-specific website subdirectories (e.g., `public/es/`)
- Uses English as the default/fallback language

## Templates

Templates are located in `templates/default/`:

- `resume.yaml.tmpl` - PDF template (YAML-based, rendered with Maroto)
- `index.html.tmpl` - Website template (HTML with Scriggo)

### PDF Template Structure

The PDF template uses a YAML structure that defines rows and columns:

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

### Available Template Functions

- `formatDate` - Format dates as YYYY-MM
- `formatEndDate` - Format end dates or show "Present"
- `formatYear` - Extract year from date
- `getEmail` - Extract email from social links
- `getPhone` - Extract phone from social links
- `hasSocials` - Check if social media links exist
- `splitLines` - Split multiline text
- `calculateHeight` - Calculate required height for text
- `assetPath` - Resolve absolute asset paths
- `lastURLPart` - Extract username from URL

## Assets

Place your assets in the `assets/` directory:

```
assets/
├── media/
│   ├── brands/          # Brand logos (GitHub, LinkedIn, etc.)
│   │   ├── github.png
│   │   └── linkedin.png
│   └── solid/           # Solid icons
│       └── envelope.png
├── images/
│   └── profile.jpg
└── files/
    └── custom.pdf
```

Assets are automatically copied to `public/assets/` during website generation.

## CLI Commands

### Global Flags

```bash
--config string       # Config file (default: $HOME/.odinnordico.github.io.yaml)
--data-dir string     # Data directory (default: "data")
--output-dir string   # Output directory (default: "public")
--lang string         # Language code (default: "en")
```

### PDF Command

```bash
go run . pdf [flags]

Flags:
  --theme string   # Theme name (default: "default")
```

### Website Command

```bash
go run . website [flags]

Flags:
  --theme string   # Theme name (default: "default")
```

### Serve Command

```bash
go run . serve [flags]

Flags:
  --host string    # Host to serve on (default: "localhost")
  --port string    # Port to serve on (default: "8080")
  --watch          # Enable live reload (default: false)
  --theme string   # Theme name (default: "default")
```

## Development

### Project Structure

```
odinnordico.github.io/
├── cmd/                    # CLI commands
│   ├── pdf.go             # PDF generation command
│   ├── website.go         # Website generation command
│   └── serve.go           # Development server command
├── internal/
│   ├── generator/         # PDF and website generators
│   │   ├── pdf.go        # PDF generation logic
│   │   ├── template.go   # Template parsing and rendering
│   │   └── website.go    # Website generation logic
│   ├── loader/           # YAML data loading
│   ├── logger/           # Logging utilities
│   ├── models/           # Data models
│   └── utils/            # Utility functions
├── templates/            # Templates
│   └── default/
│       ├── resume.yaml.tmpl
│       └── index.html.tmpl
├── data/                 # Resume data
├── assets/              # Static assets
└── public/              # Generated output
```

### Building from Source

```bash
# Build for current platform
go build -o odinnordico.github.io .

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o odinnordico.github.io-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o odinnordico.github.io-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o odinnordico.github.io-windows-amd64.exe .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/generator/...
```

## Customization

### Creating a Custom Theme

1. Create a new theme directory:
```bash
mkdir -p templates/mytheme
```

2. Copy and modify templates:
```bash
cp templates/default/resume.yaml.tmpl templates/mytheme/
cp templates/default/index.html.tmpl templates/mytheme/
```

3. Use your theme:
```bash
go run . pdf --theme mytheme
go run . website --theme mytheme
```

### PDF Customization

The PDF template supports:
- **Text**: Content, size, style (normal/bold/italic), alignment, color, hyperlinks
- **Lines**: Thickness, color
- **Images**: Path, size percentage, centering

Example:
```yaml
- height: 5
  cols:
    - width: 6
      text:
        content: "Hello World"
        size: 12
        style: bold
        align: left
        hyperlink: "https://example.com"
        color:
          red: 0
          green: 0
          blue: 255
```

## Troubleshooting

### PDF Generation Issues

**Problem**: Images not loading in PDF
- **Solution**: Ensure images are PNG format and paths are correct
- Convert images: `mogrify -format png *.jpg`

**Problem**: Text overflow in PDF
- **Solution**: Adjust `calculateHeight` in template or reduce font size

### Website Generation Issues

**Problem**: Assets not copying
- **Solution**: Check that assets exist in `assets/` directory
- Verify file permissions

### Common Errors

**Error**: `template file not found`
- Check that templates exist in `templates/default/`
- Verify `--theme` flag is correct

**Error**: `failed to load resume data`
- Ensure YAML files are valid
- Check file paths and permissions

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Maroto](https://github.com/johnfercher/maroto) - Excellent PDF generation library
- [Scriggo](https://github.com/open2b/scriggo) - Fast template engine
- [Cobra](https://github.com/spf13/cobra) - Powerful CLI framework

## Support

- 📧 Email: [your-email@example.com](mailto:your-email@example.com)
- 🐛 Issues: [GitHub Issues](https://github.com/odinnordico/odinnordico.github.io/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/odinnordico/odinnordico.github.io/discussions)
