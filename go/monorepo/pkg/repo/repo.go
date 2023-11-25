// Package repo provides models and functionalities for managing the monorepo.
package repo

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const specFile = "repo.yaml"

// Spec describes the monorepo constructs, structures, and constraints.
type Spec struct {
	Name    string
	Domains []Domain
	Teams   []Team
}

// Domain is a coarse-grained construct for organizing and managing a group of related systems.
type Domain struct {
	Name       string
	Subdomains []Subdomain
}

// Subdomain is a fine-grained construct for organizing and managing a group of related subsystems.
type Subdomain struct {
	Name string
}

// Team represents a team owning a subdomain and its projects.
type Team struct {
	Name string
}

// Read reads the specifications of the monorepo.
func Read() (Spec, error) {
	filename, err := findRepoFile(".")
	if err != nil {
		return Spec{}, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return Spec{}, err
	}
	defer f.Close()

	var spec Spec
	if err := yaml.NewDecoder(f).Decode(&spec); err != nil {
		return Spec{}, err
	}

	return spec, nil
}

func findRepoFile(path string) (string, error) {
	path, _ = filepath.Abs(path)
	filename := filepath.Join(path, specFile)

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			if parent := filepath.Dir(path); parent != "/" {
				return findRepoFile(parent)
			}
		}
		return "", err
	}

	return filename, nil
}
