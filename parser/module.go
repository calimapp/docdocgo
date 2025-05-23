package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/yuin/goldmark"
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

	doc, err := getPackageDocumentation(modulePath)
	if err != nil {
		return nil, err
	}
	moduleDoc := &goModule{
		Name:          moduleName,
		Documentation: *doc,
		Packages:      make([]goPackage, 0),
		Readme:        *moduleReadme,
		SourceFiles:   getSourceFiles(modulePath),
	}

	err = filepath.Walk(modulePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
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
		moduleDoc.Packages = append(moduleDoc.Packages, *pkg)
		return nil
	})
	return moduleDoc, err
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
