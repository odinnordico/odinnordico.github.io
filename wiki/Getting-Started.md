# Getting Started

This guide will walk you through creating your first resume with odinnordico.github.io.

## 1. Initialize Your Workspace

If you haven't already, clone the repository to get the example data and templates. Even if you installed the binary globally, having the project structure is helpful for starting out.

```bash
git clone https://github.com/odinnordico/odinnordico.github.io.git my-resume
cd my-resume
```

The directory structure looks like this:

```text
my-resume/
├── data/               # Your resume content (YAML files)
├── templates/          # Visual templates
├── assets/             # Images and icons
└── public/             # Generated output (PDFs and Website)
```

## 2. Edit Your Data

Navigate to the `data/` directory. You will see several YAML files. These contain the content of your resume.

Open `data/basic.yml` and update it with your information:

```yaml
name: "Jane Doe"
display_name: "Jane Doe"
location: "New York, NY"
summary: |
  Passionate developer building amazing tools...
website: "https://janedoe.dev"
```

Edit `data/professional.yml` to add your work experience:

```yaml
title: "Software Engineer"
years_of_experience: 5
jobs:
  - position: "Backend Developer"
    company:
      name: "Startup Inc"
      url: "https://startup.inc"
    start_date: 2021-05-01
    end_date: null
    job_description: |
      - Built scalable APIs using Go
      - Optimized database queries
```

Explore the other files (`education.yml`, `skills.yml`, `social.yml`) and update them as needed. See [Data Structure](Data-Structure) for a detailed reference.

## 3. Generate Your Resume

Now that your data is ready, let's generate the resume.

### Generate PDF

Run the following command to generate the PDF version:

```bash
odinnordico.github.io pdf
```

This will create `public/assets/files/resume.pdf`.

### Generate Website

To generate the static website:

```bash
odinnordico.github.io website
```

This will create `public/index.html` and copy necessary assets.

## 4. Live Development

For the best experience, use the `serve` command. This starts a local web server and watches your files for changes.

```bash
odinnordico.github.io serve --watch
```

Open your browser to `http://localhost:8080`. You will see your resume website.

Now, try editing `data/basic.yml` and saving the file. The server will detect the change, regenerate the website and PDF, and reload the page automatically!

## Next Steps

*   **[Configuration](Configuration)**: Learn how to customize paths and settings.
*   **[Templates](Templates-and-Customization)**: Change the look of your resume.
*   **[Multi-language](Multi-language-Support)**: Add a second language.
