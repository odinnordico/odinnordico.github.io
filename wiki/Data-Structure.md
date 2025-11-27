# Data Structure

odinnordico.github.io uses a set of YAML files in the `data/` directory to define your resume content. This separation of data and presentation allows you to easily update your information without touching the design.

## File Overview

| File | Purpose |
|------|---------|
| `basic.yml` | Personal information (Name, Contact, Summary) |
| `professional.yml` | Work experience |
| `education.yml` | Academic history |
| `skills.yml` | Technical skills and proficiency |
| `certificates.yml` | Certifications and awards |
| `social.yml` | Social media profiles |
| `projects.yml` | (Optional) Personal projects |

## Detailed Schemas

### `basic.yml`

Contains your core identity information.

```yaml
name: "John Doe"          # Full name
display_name: "John Doe"  # Name to display (can be same as name)
location: "City, Country" # Current location
summary: |                # Professional summary (multiline)
  Experienced developer...
website: "https://..."    # Personal website URL
```

### `professional.yml`

Describes your career history.

```yaml
title: "Senior Engineer"  # Current professional title
years_of_experience: 8    # Total years of experience
jobs:
  - position: "Role Title"
    company:
      name: "Company Name"
      url: "https://..."
      logo:               # Optional company logo
        library: "brands" # 'brands' or 'solid' or 'regular'
        image: "google"   # icon name (without extension)
    start_date: 2020-01-01 # YYYY-MM-DD
    end_date: null        # null for current job
    job_description: |    # Bullet points or text
      - Achievement 1
      - Achievement 2
```

### `education.yml`

Lists your academic background.

```yaml
- title: "Degree Name"
  level: "Bachelor's/Master's"
  provider:
    name: "University Name"
    url: "https://..."
  date: 2015-05-20        # Graduation date
  description: "Notes..." # Honors, GPA, etc.
```

### `skills.yml`

Showcases your technical abilities.

```yaml
- name: "Skill Name"      # e.g., "Go"
  description: "Details"  # e.g., "Backend, CLI tools"
  level: 9                # Proficiency (1-10)
  logo:                   # Icon for the skill
    library: "devicon"
    image: "go"
  tags:                   # Categories
    - "backend"
```

### `social.yml`

Links to your online presence.

```yaml
- name: "GitHub"
  url: "https://github.com/..."
  logo:
    library: "brands"
    image: "github"
```

### `certificates.yml`

Professional certifications.

```yaml
- name: "Cert Name"
  description: "Details"
  date: 2023-01-01
  certificate_url: "https://..." # Link to credential
  provider:
    name: "Issuer Name"
    url: "https://..."
  topics:
    - "Topic 1"
```

## Tips

*   **Dates**: Use `YYYY-MM-DD` format.
*   **Multiline Strings**: Use `|` for multiline strings in YAML to preserve formatting.
*   **Null Values**: Use `null` or omit fields that are not applicable.
