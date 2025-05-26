package parser

import (
	"os"
	"path/filepath"

	"github.com/google/licensecheck"
)

var licenseFileNames []string = []string{
	"COPYING", "COPYING.md", "COPYING.markdown", "COPYING.txt",
	"LICENCE", "LICENCE.md", "LICENCE.markdown", "LICENCE.txt", "LICENCE-2.0.txt", "MIT-LICENCE",
	"LICENSE", "LICENSE.md", "LICENSE.markdown", "LICENSE.txt", "LICENSE-2.0.txt", "MIT-LICENSE",
}

func resolveLicense(modulePath string) string {
	for _, name := range licenseFileNames {
		fullPath := filepath.Join(modulePath, name)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			continue // Try next file
		}

		// Try to detect the license type
		result := licensecheck.Scan(data)

		if result.Match != nil {
			return result.Match[0].ID
		}
	}
	return "None"
}
