# Proj2Prompt

Proj2Prompt is a command-line tool designed to generate project structures formatted for use as prompts in large language models (LLMs). It efficiently explores directories, applies exclusion rules, and outputs structured descriptions of file hierarchies and file contents.

## Features

- Generate structured representations of project directories.
- Apply `.gitignore` and custom exclusion patterns.
- Save output to a file or copy directly to the clipboard.
- Automatically detects and describes binary files.

## Installation

### Prerequisites
- Go 1.23 or higher.

### Build
Clone the repository and build the executable:
```bash
git clone https://github.com/dlubom/proj2prompt.git
cd proj2prompt
go build -o proj2prompt
```

### Run
Run the tool using the generated executable:
```bash
./proj2prompt [flags] [directory]
```

## Usage

### Syntax
```bash
proj2prompt [flags] [directory]
```

### Flags
- `-o, --output`: Save the output to a file.
- `-e, --exclude`: Specify exclusion patterns (e.g., `*.tmp`, `node_modules`).
- `-c, --clipboard`: Copy output to clipboard.
- `-v, --version`: Display the application version.

### Examples

#### Explore a directory and print the structure
```bash
proj2prompt ./my_project
```

#### Save the output to a file
```bash
proj2prompt ./my_project -o structure.txt
```

#### Exclude specific files or directories
```bash
proj2prompt ./my_project -e "*.log" -e "node_modules"
```

#### Copy output to clipboard
```bash
proj2prompt ./my_project -c
```