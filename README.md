# Gosnap

**Gosnap** is a command-line utility written in Go that generates a snapshot of a directory's structure and the textual contents of its files. It recursively scans a specified directory, lists its structure (files and folders), and appends the contents of text-based files to an output file, while excluding binary files and common development artifacts by default.

## Features

- **Directory Structure Snapshot**: Displays a tree-like structure of files and folders.
- **Text File Contents**: Extracts and includes the contents of text files (UTF-8 encoded) in the output.
- **Binary File Filtering**: Automatically skips binary and non-UTF-8 files.
- **Exclude Common Artifacts**: Ignores development-related files/folders (e.g., `.git`, `.idea`, `.venv`, `lib`, `test`, `log`, `etc`) by default.
- **Custom Exclusions**: Allows manual exclusion of specific files or folders.
- **File Extension Filtering**: Includes only files with specified extensions (e.g., `.py`, `.go`).
- **Flexible Output**: Saves the snapshot to a user-specified file (default: `snap.txt`).
- **Default Input Directory**: Uses the current directory (`.`) if no input directory is specified.

## Installation

### Prerequisites
- **Go**: Ensure you have Go 1.24.3 or later installed. Download it from [golang.org](https://golang.org/dl/).
- **Git**: Required to clone the repository.

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/Wingrull/gosnap.git
   cd gosnap
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build the binary:
   ```bash
   go build -o gosnap
   ```
4. (Optional) Move the binary to a directory in your `PATH` for global access:
   ```bash
   mv gosnap /usr/local/bin/
   ```

## Usage

Run `gosnap` with the following command-line flags:

- `-e, --exclude <name>`: Manually exclude specific files or folders (can be used multiple times).
- `-en, --exclude-noise`: Automatically exclude common development artifacts (e.g., `.git`, `.idea`, `.venv`, `lib`, `test`, `log`, `etc`) (default: `true`). Set to `false` to include them.
- `-ext, --extension <ext>`: Include only files with specified extensions (e.g., `.py`, `.go`) (can be used multiple times).
- `-o, --output <file>`: Specify the output file path (default: `snap.txt`).

The input directory is optional; if not provided, the current directory (`.`) is used.

### Default Exclusions
When `-en/--exclude-noise` is enabled (default: `true`), the following files and folders are automatically excluded:
- `.git`
- `.venv`
- `__pycache__`
- `node_modules`
- `.idea`
- `.DS_Store`
- `lib`
- `test`
- `log`
- `etc`

To include these artifacts in the snapshot, use `-en=false`.

### Examples

1. **Generate a snapshot of the current directory**:
   ```bash
   ./gosnap
   ```
   This creates `snap.txt` with the directory structure and contents of text files, excluding artifacts like `.git`, `.idea`, `lib`, `test`, `log`, and `etc`.

   **Output in `snap.txt`**:
   ```
   Directory Structure:
   build.bat
   go.mod
   go.sum
   main.go
   snap.txt

   === File Contents ===

   File: build.bat
   @echo off
   pushd "%~dp0"
   ...

   File: go.mod
   module gosnap
   ...

   File: go.sum
   golang.org/x/text v0.25.0 ...
   ...

   File: main.go
   package main
   ...
   ```

2. **Snapshot only Python files in a specific directory**:
   ```bash
   ./gosnap /path/to/project -ext .py
   ```
   This processes `/path/to/project`, including only `.py` files and excluding `lib`, `test`, `log`, `etc`, and other artifacts.

   **Output in `snap.txt`**:
   ```
   Directory Structure:
   main.py
   utils.py

   === File Contents ===

   File: main.py
   <contents of main.py>

   File: utils.py
   <contents of utils.py>
   ```

3. **Include development artifacts**:
   ```bash
   ./gosnap -en=false -ext .py .
   ```
   This includes folders like `.idea`, `lib`, `test`, `log`, and `etc`, but only processes `.py` files.

4. **Specify a custom output file and multiple extensions**:
   ```bash
   ./gosnap -o snapshot.txt -ext .py -ext .go .
   ```
   This writes the snapshot to `snapshot.txt`, including only `.py` and `.go` files.

## Output Format

The output file (e.g., `snap.txt`) contains:
1. **Directory Structure**: A tree-like representation of files and folders.
2. **File Contents**: The contents of each text file, prefixed with its relative path.

Example:
```
Directory Structure:
src/
  main.py
  utils.py

=== File Contents ===

File: src/main.py
<contents of main.py>

File: src/utils.py
<contents of utils.py>
```

## Dependencies
- `golang.org/x/text`: Used for UTF-8 encoding detection.

Install it with:
```bash
go get golang.org/x/text
```

## Contributing
Contributions are welcome! To contribute:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit (`git commit -m "Add feature"`).
4. Push to your fork (`git push origin feature-branch`).
5. Open a Pull Request.

Please ensure your code follows Go best practices and includes tests if applicable.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact
For questions or suggestions, open an issue on [GitHub](https://github.com/Wingrull/gosnap/issues) or contact [Wingrull](https://github.com/Wingrull).