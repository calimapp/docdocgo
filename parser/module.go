package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/yuin/goldmark"
	"golang.org/x/mod/modfile"
)

func ParseModule(modulePath string) (*goModule, error) {
	moduleName, err := getModuleName(modulePath)
	if err != nil {
		return nil, err
	}

	moduleReadme, err := getModuleReadme(modulePath)
	if err != nil {
		return nil, err
	}

	moduleDoc, err := getPackageDocumentation(modulePath)
	if err != nil {
		return nil, err
	}

	moduleDependencies, err := getModuleDependencies(modulePath)
	if err != nil {
		return nil, err
	}

	moduleVersion, err := getVersion(moduleName)
	if err != nil {
		moduleVersion.Time = time.Now()
		moduleVersion.Version = ""
	}

	module := &goModule{
		Name:          moduleName,
		Version:       moduleVersion.Version,
		Date:          moduleVersion.Time.Format(time.DateOnly),
		Documentation: *moduleDoc,
		License:       resolveLicense(modulePath),
		Packages:      make([]goPackage, 0),
		Readme:        *moduleReadme,
		SourceFiles:   getSourceFiles(modulePath),
		Dependencies:  moduleDependencies,
	}

	err = filepath.Walk(modulePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || path == modulePath {
			return nil
		}

		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		hasGoFiles := false
		for _, file := range files {
			if filepath.Ext(file.Name()) == ".go" {
				hasGoFiles = true
				break
			}
		}

		if !hasGoFiles {
			return nil
		}
		pkg, err := parsePackage(path)
		if err != nil {
			return err
		}
		module.Packages = append(module.Packages, *pkg)
		return nil
	})
	return module, err
}

func getModuleName(modulePath string) (string, error) {
	goModPath := filepath.Join(modulePath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("error reading go.mod file %s: %w", goModPath, err)
	}

	re := regexp.MustCompile(`module\s+([^\s]+)`)
	matches := re.FindSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("module name not found in go.mod file %s", goModPath)
	}

	return string(matches[1]), nil
}

// getModuleReadme return the README.md file if exists or nil otherwise
func getModuleReadme(modulePath string) (*template.HTML, error) {
	readmePath := filepath.Join(modulePath, "README.md")
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, err
	}
	var readmeHtml bytes.Buffer
	if err := goldmark.Convert(readmeContent, &readmeHtml); err != nil {
		return nil, err
	}
	html := template.HTML(readmeHtml.String())
	return &html, nil
}

type moduleVersion struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}

func getVersion(moduleRef string) (*moduleVersion, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", moduleRef)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("module not found: " + resp.Status)
	}

	var mv moduleVersion
	if err := json.NewDecoder(resp.Body).Decode(&mv); err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	return &mv, nil
}

type Dependency struct {
	Path     string `json:"Path"`
	Version  string `json:"Version"`
	Indirect bool   `json:"Indirect"`
}

func getModuleDependencies(modulePath string) ([]Dependency, error) {
	data, err := os.ReadFile(filepath.Join(modulePath, "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("failed to read go.mod file: %w", err)
	}
	gomod, err := modfile.Parse("go.mod", data, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod file: %w", err)
	}
	var dependencies []Dependency
	for _, dep := range gomod.Require {
		dependencies = append(dependencies, Dependency{
			Path:     dep.Mod.Path,
			Version:  dep.Mod.Version,
			Indirect: dep.Indirect,
		})
	}
	return dependencies, nil
}
