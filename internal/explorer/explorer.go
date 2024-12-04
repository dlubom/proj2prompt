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

		// Add a descriptive message and separator before each path
		if info.IsDir() {
			builder.WriteString(fmt.Sprintf("\n#######\nDirectory: %s\n#######\n", relativePath))
		} else {
			builder.WriteString(fmt.Sprintf("\n-----\nFile: %s\n-----\n", relativePath))
		}

		// Include file content for files
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

	// Read up to 8000 bytes or the full file if smaller
	maxSampleSize := 8000
	sampleSize := int(info.Size())
	if sampleSize > maxSampleSize {
		sampleSize = maxSampleSize
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, sampleSize)
	bytesRead, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	buffer = buffer[:bytesRead]

	// Check if the file is binary
	if isBinary(buffer) {
		// Display the first 10 bytes in hex
		hexBytes := ""
		for i := 0; i < len(buffer) && i < 10; i++ {
			hexBytes += fmt.Sprintf("%02x ", buffer[i])
		}
		return fmt.Sprintf("[Binary file: %s, Size: %d bytes, First 10 bytes: %s]", info.Name(), info.Size(), hexBytes), nil
	} else {
		// Treat as text file
		return string(buffer), nil
	}
}

// Helper function to determine if a file is binary
func isBinary(data []byte) bool {
	// Define a threshold for non-text bytes
	// If more than 30% of the sample are non-text, consider it binary
	nonTextThreshold := 0.3
	nonTextCount := 0

	for _, b := range data {
		// Check for null bytes or non-printable characters
		if (b == 0) || (b > 0 && b < 9) || (b > 13 && b < 32) || b == 127 {
			nonTextCount++
		}
	}

	return float64(nonTextCount)/float64(len(data)) > nonTextThreshold
}
