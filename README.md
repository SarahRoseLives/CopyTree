# 📋 CopyTree

CopyTree is a command-line tool that walks through your current directory, collects source code files (filtered by extension), prints a visual tree of matched files, and copies the content of those files into your clipboard. This is especially useful when sharing code with AI tools or collaborators, allowing for instant copy-paste of structured project code.

---

## 🚀 Features

- 🗂 **Recursively scans the current directory**
- 🧠 **Filters files by one or more extensions** (e.g. `go`, `js`, `py`)
- 🌲 **Prints a formatted tree** of included files
- 📋 **Copies all matching files' content** into your clipboard in a structured format
- 📊 **Outputs a summary**: file count, line count, and character count

---

## 🛠 Installation

### 1. Clone the repository

```bash
git clone https://github.com/SarahRoseLives/CopyTree.git
cd CopyTree
```

### 2. Install Go (if you haven't)

Make sure you have [Go installed](https://golang.org/dl/).

### 3. Build the executable

```bash
go build -o copytree main.go
```

---

## ✨ Usage

### Basic usage (all files)

```bash
./copytree
```

- Scans the current directory and subdirectories
- Prints a tree of all files
- Copies all file contents to your clipboard

### Filter by extension

```bash
./copytree go js py
```

- Only includes files ending with `.go`, `.js`, or `.py`

### Example output

```plaintext
.
├── main.go
├── utils.go
└── README.md
Copied Dir Tree and 3 files to clipboard
Total lines 200 Total Characters 5274
```

- The above is printed to your terminal.
- Your clipboard now contains:

```
====./main.go====
<contents of main.go>

====./utils.go====
<contents of utils.go>

====./README.md====
<contents of README.md>
```

---

## 📝 Notes

- The tool works from your current directory. To copy from your project root, run it from there.
- Clipboard integration uses [atotto/clipboard](https://github.com/atotto/clipboard), which supports macOS, Windows, and Linux with X11.
- If you encounter clipboard issues on Linux, ensure you have `xclip` or `xsel` installed.

---

## 💡 Example Scenarios

- **Share code with ChatGPT or other LLMs**: Instantly copy all relevant code files for context.
- **Send snippets to a collaborator**: Grab only your `go` and `yaml` files for a quick review.
- **Create code archives**: Easily gather source files by extension for documentation or backup.
