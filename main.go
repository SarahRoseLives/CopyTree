package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

const (
	chatgptSectionSize = 20000 // green output length in characters
)

func main() {
	// Handle --chatgpt flag and file extensions
	var chatgptMode bool
	extensions := []string{}
	for _, arg := range os.Args[1:] {
		if arg == "--chatgpt" {
			chatgptMode = true
		} else if strings.HasPrefix(arg, "--") {
			// ignore unknown flags for now
			continue
		} else {
			extensions = append(extensions, strings.TrimPrefix(strings.ToLower(arg), "."))
		}
	}

	startDir, _ := os.Getwd()

	var fileList []string
	var fileCount, lineCount, charCount int
	var buf bytes.Buffer

	// Walk the directory tree
	filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip error files
		}
		if d.IsDir() {
			return nil
		}
		// Filter by extension if provided
		if len(extensions) > 0 {
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(d.Name())), ".")
			found := false
			for _, e := range extensions {
				if ext == e {
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}
		fileList = append(fileList, path)
		return nil
	})

	// Build tree for both printing and clipboard
	tree := buildTree(startDir, fileList)

	// Print tree to stdout
	printTreeRec(tree, "", true)

	// Add tree string to clipboard buffer
	treeString := buildTreeString(tree, "", true)
	buf.WriteString(treeString)
	buf.WriteString("\n\n")

	// Copy files to clipboard buffer
	for _, path := range fileList {
		rel, _ := filepath.Rel(startDir, path)
		buf.WriteString(fmt.Sprintf("====%s====\n", filepath.Join(".", rel)))
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		buf.Write(data)
		if len(data) > 0 && data[len(data)-1] != '\n' {
			buf.WriteByte('\n')
		}
		fileCount++
		lineCount += bytes.Count(data, []byte{'\n'})
		charCount += len(data)
	}

	clipboardData := buf.String()

	if !chatgptMode {
		// Normal mode: Copy all at once
		clipboard.WriteAll(clipboardData)
		fmt.Printf("Copied Dir Tree and %d files to clipboard\n", fileCount)
		// Color the summary based on size
		var colorCode string
		switch {
		case charCount <= 20000: // Green - should work with any AI
			colorCode = "\033[32m"
		case charCount <= 50000: // Yellow - should work with some like DeepSeek or Gemini
			colorCode = "\033[33m"
		default: // Red - highly unlikely to work with AI
			colorCode = "\033[31m"
		}
		fmt.Printf("%sTotal lines %d Total Characters %d\033[0m\n", colorCode, lineCount, charCount)
		return
	}

	// ChatGPT mode
	// Section the output into <= 20000 char sections, each split at a boundary if possible
	sections := splitIntoSections(clipboardData, chatgptSectionSize)
	// Add prompt to first section
	chatgptPrompt := `I have a lot of files to show you, I'm going to send you each section separately.
Tell me when you're ready for the first file.
Then continue to ask for the next file until we have completed all copying.

`
	sections[0] = chatgptPrompt + sections[0]

	fmt.Printf("Entering ChatGPT mode: splitting output into %d sections (~%d chars each)\n", len(sections), chatgptSectionSize)
	fmt.Printf("Copied section 1 of %d to clipboard. Paste into ChatGPT, then press [Enter] for next section.\n", len(sections))
	clipboard.WriteAll(sections[0])

	reader := bufio.NewReader(os.Stdin)
	for i := 1; i < len(sections); i++ {
		fmt.Printf("[Section %d/%d] Press [Enter] to copy next section to clipboard...", i+1, len(sections))
		reader.ReadString('\n')
		clipboard.WriteAll(sections[i])
		fmt.Printf("Section %d copied to clipboard!\n", i+1)
	}
	fmt.Printf("All sections copied! (Total files: %d, Total lines: %d, Total Characters: %d)\n", fileCount, lineCount, charCount)
}

// Helper: splits text into sections not exceeding maxLen, tries to split on file header
func splitIntoSections(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}
	lines := strings.Split(text, "\n")
	var sections []string
	var buf strings.Builder
	for _, line := range lines {
		// If adding this line would exceed maxLen, start a new section
		if buf.Len()+len(line)+1 > maxLen && buf.Len() > 0 {
			sections = append(sections, buf.String())
			buf.Reset()
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(line)
	}
	if buf.Len() > 0 {
		sections = append(sections, buf.String())
	}
	return sections
}

// Helper function to print a filtered tree
func printTreeRec(node *TreeNode, prefix string, last bool) {
	if node.Name != "." {
		fmt.Printf("%s", prefix)
		if last {
			fmt.Print("└── ")
		} else {
			fmt.Print("├── ")
		}
		fmt.Println(node.Name)
	}
	keys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		keys = append(keys, k)
	}
	sortStrings(keys)
	for i, k := range keys {
		child := node.Children[k]
		isLast := i == len(keys)-1
		var newPrefix string
		if node.Name == "." {
			newPrefix = ""
		} else if last {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}
		printTreeRec(child, newPrefix, isLast)
	}
}

// Helper function to build a filtered tree as a string (for clipboard)
func buildTreeString(node *TreeNode, prefix string, last bool) string {
	var out strings.Builder
	if node.Name != "." {
		out.WriteString(prefix)
		if last {
			out.WriteString("└── ")
		} else {
			out.WriteString("├── ")
		}
		out.WriteString(node.Name)
		out.WriteString("\n")
	}
	keys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		keys = append(keys, k)
	}
	sortStrings(keys)
	for i, k := range keys {
		child := node.Children[k]
		isLast := i == len(keys)-1
		var newPrefix string
		if node.Name == "." {
			newPrefix = ""
		} else if last {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}
		out.WriteString(buildTreeString(child, newPrefix, isLast))
	}
	return out.String()
}

// Converts file list into a tree structure
type TreeNode struct {
	Name     string
	Children map[string]*TreeNode
	IsFile   bool
}

func buildTree(base string, files []string) *TreeNode {
	root := &TreeNode{Name: ".", Children: make(map[string]*TreeNode)}
	for _, f := range files {
		rel, _ := filepath.Rel(base, f)
		parts := strings.Split(rel, string(os.PathSeparator))
		node := root
		for i, part := range parts {
			if _, ok := node.Children[part]; !ok {
				node.Children[part] = &TreeNode{
					Name:     part,
					Children: make(map[string]*TreeNode),
					IsFile:   i == len(parts)-1,
				}
			}
			node = node.Children[part]
		}
	}
	return root
}

// Sort strings (for consistent tree output)
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
