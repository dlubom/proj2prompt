package explorer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Explorer is responsible for traversing directories and applying filters
type Explorer struct {
	RootPath        string
	ExcludePatterns []string
	gitIgnore       *ignore.GitIgnore
}

// NewExplorer initializes a new Explorer instance
func NewExplorer(rootPath string, excludePatterns []string) *Explorer {
	exp := &Explorer{
		RootPath:        rootPath,
		ExcludePatterns: excludePatterns,
	}
	exp.loadGitIgnore()
	return exp
}

// Load .gitignore file if present
func (e *Explorer) loadGitIgnore() {
	gitIgnorePath := filepath.Join(e.RootPath, ".gitignore")
	if _, err := os.Stat(gitIgnorePath); err == nil {
		gitIgnore, compileErr := ignore.CompileIgnoreFile(gitIgnorePath)
		if compileErr == nil {
			e.gitIgnore = gitIgnore
		} else {
			fmt.Printf("Warning: failed to compile .gitignore: %v\n", compileErr)
		}
	}
}

// Check if a file or directory should be excluded
func (e *Explorer) shouldExclude(path string, info fs.FileInfo) bool {
	// Exclude .git directory
	if info.IsDir() && info.Name() == ".git" {
		return true
	}

	// Exclude files/directories listed in .gitignore
	if e.gitIgnore != nil && e.gitIgnore.MatchesPath(path) {
		return true
	}

	// Apply user-defined exclusion patterns
	for _, pattern := range e.ExcludePatterns {
		matched, err := filepath.Match(pattern, info.Name())
		if err == nil && matched {
			return true
		}
	}

	return false
}

// Explore traverses the directory and collects its structure and content
func (e *Explorer) Explore() (string, error) {
	var builder strings.Builder

	err := filepath.Walk(e.RootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if e.shouldExclude(path, info) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relativePath, err := filepath.Rel(e.RootPath, path)
		if err != nil {
			return err
		}

		// Add the path to the result
		builder.WriteString(relativePath)
		builder.WriteString("\n")

		// Include file content for text files
		if !info.IsDir() {
			content, err := e.readFileContent(path, info)
			if err != nil {
				return err
			}
			builder.WriteString(content)
			builder.WriteString("\n")
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

// Read file content; handle text and binary files differently
func (e *Explorer) readFileContent(path string, info fs.FileInfo) (string, error) {
	if info.Size() == 0 {
		return "", nil
	}

	// Identify text files by extensions
	textExtensions := []string{".txt", ".md", ".go", ".py", ".java", ".js", ".html", ".css"}
	isTextFile := false
	for _, ext := range textExtensions {
		if strings.HasSuffix(info.Name(), ext) {
			isTextFile = true
			break
		}
	}

	if isTextFile {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	} else {
		// Handle binary files
		return fmt.Sprintf("[Binary file: %s, Size: %d bytes]", info.Name(), info.Size()), nil
	}
}
