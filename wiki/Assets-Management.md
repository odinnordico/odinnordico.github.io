# Assets Management

odinnordico.github.io has a dedicated `assets/` directory for managing all static files used in your resume, such as images, logos, and downloadable files.

## Directory Structure

```text
assets/
├── media/
│   ├── brands/          # Brand logos (e.g., GitHub, LinkedIn)
│   │   ├── github.png
│   │   └── linkedin.png
│   ├── solid/           # Solid icons (e.g., envelope, phone)
│   │   └── envelope.png
│   └── regular/         # Regular icons
├── images/              # Personal images
│   └── profile.jpg
└── files/               # Files to be included in the website
    └── custom.pdf
```

## Using Assets

### In YAML Data

When referencing logos in your YAML data files (like `social.yml` or `skills.yml`), you specify the library and the image name (without extension).

```yaml
logo:
  library: "brands"  # Looks in assets/media/brands/
  image: "github"    # Looks for github.png
```

### In Templates

In your templates, you can use the `assetPath` function to get the absolute path to an asset file.

**PDF Template:**
```yaml
image:
  path: '{{ assetPath "images/profile.jpg" }}'
```

**Website Template:**
Assets are automatically copied to `public/assets/` during website generation. You can reference them relatively.

```html
<img src="assets/images/profile.jpg" alt="Profile Photo">
```

## Supported Formats

*   **Images**: PNG and JPG are supported. PNG is recommended for logos and icons to preserve transparency.
*   **PDF**: The generated resume is saved to `assets/files/resume.pdf` so it can be linked for download.

## Adding New Icons

To add a new icon:

1.  Find a PNG version of the icon (e.g., from [FontAwesome](https://fontawesome.com/) or [Devicon](https://devicon.dev/)).
2.  Place it in the appropriate subdirectory under `assets/media/`.
3.  Reference it in your YAML data.
