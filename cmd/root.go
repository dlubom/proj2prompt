package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/dlubom/proj2prompt/internal/explorer"
    "github.com/dlubom/proj2prompt/internal/clipboard"
)

var (
    outputPath      string
    excludePatterns []string
    toClipboard     bool
    version         = "1.0.0"
)

// Root command definition
var rootCmd = &cobra.Command{
    Use:   "proj2prompt [directory]",
    Short: "Proj2Prompt generates a project structure for LLM prompts",
    Run: func(cmd *cobra.Command, args []string) {
        // Set the root directory; defaults to current directory
        rootPath := "."
        if len(args) > 0 {
            rootPath = args[0]
        }

        // Initialize the directory explorer
        exp := explorer.NewExplorer(rootPath, excludePatterns)
        result, err := exp.Explore()
        if err != nil {
            fmt.Println("Error exploring directories:", err)
            os.Exit(1)
        }

        // Handle output to file if specified
        if outputPath != "" {
            err := os.WriteFile(outputPath, []byte(result), 0644)
            if err != nil {
                fmt.Println("Error writing to file:", err)
                os.Exit(1)
            }
            fmt.Printf("Output written to file: %s\n", outputPath)
        }

        // Handle clipboard copy if requested
        if toClipboard {
            err := clipboard.CopyToClipboard(result)
            if err != nil {
                fmt.Println("Error copying to clipboard:", err)
                os.Exit(1)
            }
            fmt.Println("Output copied to clipboard.")
        }

        // Default to printing in the terminal
        if outputPath == "" && !toClipboard {
            fmt.Println(result)
        }
    },
}

// Initialize CLI flags and commands
func init() {
    rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Save the output to a file")
    rootCmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{}, "Add file/directory exclusion rules")
    rootCmd.Flags().BoolVarP(&toClipboard, "clipboard", "c", false, "Copy output to clipboard")
    rootCmd.Flags().BoolP("version", "v", false, "Display the application version")

    rootCmd.SetVersionTemplate("Proj2Prompt version {{.Version}}\n")
    rootCmd.Version = version
}

// Execute the root command
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
