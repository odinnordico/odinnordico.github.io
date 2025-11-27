// Package models defines the data structures for resume information.
package models

import "time"

// ResumeData represents the complete resume data structure containing all sections.
type ResumeData struct {
	Basic        BasicData        `yaml:"basic,omitempty"`
	Professional ProfessionalData `yaml:"professional,omitempty"`
	Certificates []Certificate    `yaml:"certificates,flow"`
	Education    []Education      `yaml:"education,flow"`
	Skills       []Skill          `yaml:"skills,flow"`
	Social       []Entity         `yaml:"social,flow"`
}

// BasicData contains basic personal information.
type BasicData struct {
	Name          string `yaml:"name"`
	DisplayName   string `yaml:"display_name,omitempty"`
	Location      string `yaml:"location,omitempty"`
	Pronunciation string `yaml:"pronunciation,omitempty"`
	Phrase        string `yaml:"phrase,omitempty"`
	Summary       string `yaml:"summary,omitempty"`
	Website       string `yaml:"website,omitempty"`
}

// ProfessionalData contains professional experience information.
type ProfessionalData struct {
	Title             string  `yaml:"title,omitempty"`
	YearsOfExperience float64 `yaml:"years_of_experience,omitempty"`
	Jobs              []Job   `yaml:"jobs,flow"`
}

// Job represents a work experience entry.
type Job struct {
	Position       string     `yaml:"position,omitempty"`
	StartDate      time.Time  `yaml:"start_date,omitempty"`
	EndDate        *time.Time `yaml:"end_date,omitempty"` // nil means current position
	JobDescription string     `yaml:"job_description,omitempty"`
	Company        Entity     `yaml:"company,omitempty"`
}

// Entity represents an organization or social media account with optional logo.
type Entity struct {
	Name string `yaml:"name,omitempty"`
	URL  string `yaml:"url,omitempty"`
	Logo Logo   `yaml:"logo,omitempty"`
}

// Logo represents branding information for an entity.
type Logo struct {
	Library string `yaml:"library,omitempty"` // Icon library (e.g., "brands", "solid")
	Image   string `yaml:"image,omitempty"`   // Icon name
}

// Certificate represents a professional certification.
type Certificate struct {
	Name           string    `yaml:"name,omitempty"`
	Description    string    `yaml:"description,omitempty"`
	Date           time.Time `yaml:"date,omitempty"`
	CertificateURL string    `yaml:"certificate_url,omitempty"`
	URL            string    `yaml:"url,omitempty,omitempty"`
	Provider       Entity    `yaml:"provider,omitempty"`
	Topics         []string  `yaml:"topics,flow"`
}

// Education represents an education entry.
type Education struct {
	Title       string    `yaml:"title,omitempty"`
	Date        time.Time `yaml:"date,omitempty"`
	URL         string    `yaml:"url,omitempty"`
	Level       string    `yaml:"level,omitempty"`
	Provider    Entity    `yaml:"provider,omitempty"`
	Description string    `yaml:"description,omitempty"`
}

// Skill represents a skill entry with proficiency level.
type Skill struct {
	Name        string   `yaml:"name,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Level       int      `yaml:"level,omitempty" validate:"min=1,max=10"`
	Logo        Logo     `yaml:"logo,omitempty"`
	Tags        []string `yaml:"tags,flow"`
}
