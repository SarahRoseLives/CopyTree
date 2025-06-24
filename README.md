# 📋 CopyTree

CopyTree is a command-line tool that recursively scans your current directory, collects source code files (filtered by extension), prints a visual tree of matched files, and copies the content of those files into your clipboard. This makes it easy to share structured code with AI tools or collaborators—just copy, paste, and go!

---

## 🚀 Features

- 🗂 **Recursively scans the current directory and subdirectories**
- 🧠 **Filters files by one or more extensions** (e.g. `go`, `js`, `py`)
- 🌲 **Prints a formatted tree** of included files for easy visualization
- 📋 **Copies all matching files' content** into your clipboard in a structured format
- 📊 **Outputs a summary**: file count, line count, and character count
- 🤖 **ChatGPT mode**: Splits output into manageable sections for easier LLM pasting

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

### Filter by extensions

```bash
./copytree go js py
```

- Only includes files ending with `.go`, `.js`, or `.py`

### ChatGPT sectioning mode

```bash
./copytree --chatgpt
```

- Splits the output into sections of up to 20,000 characters, ideal for pasting into ChatGPT or other LLMs with context limits.

### Example output

```plaintext
.
├── main.go
├── utils.go
└── README.md
Copied Dir Tree and 3 files to clipboard
Total lines 200 Total Characters 5274
```

---

## 📋 Clipboard format

After running, your clipboard will contain **both the directory tree and the file contents**, in a format like this:

```
.
├── main.go
├── utils.go
└── README.md

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
- On Linux, make sure you have `xclip` or `xsel` installed for clipboard support.
- Unknown CLI flags are ignored (except `--chatgpt`).

---

## 💡 Example Scenarios

- **Share code with ChatGPT or other LLMs**: Instantly copy all relevant code files for context.
- **Send snippets to a collaborator**: Grab only your `go` and `yaml` files for a quick review.
- **Create code archives**: Easily gather source files by extension for documentation or backup.

---

## 🧑‍💻 Advanced

### Sectioning and AI limits

- The `--chatgpt` flag splits the output into 20,000 character sections, with prompts to make it easy to paste sequentially into LLMs with limited context windows.
- Summaries are color-coded:  
  - **Green** (≤20,000 chars): Safe for most AIs  
  - **Yellow** (≤50,000 chars): May work with larger context AIs  
  - **Red** (>50,000 chars): Unlikely to work in one paste  
