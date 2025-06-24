package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	// Get extensions from command line args (without the dot)
	exts := os.Args[1:]
	for i, e := range exts {
		exts[i] = strings.TrimPrefix(strings.ToLower(e), ".")
	}
	startDir, _ := os.Getwd()

	var buf bytes.Buffer
	var fileList []string
	var fileCount, lineCount, charCount int

	// Walk the directory tree
	filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip error files
		}
		if d.IsDir() {
			return nil
		}
		// Filter by extension if provided
		if len(exts) > 0 {
			ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(d.Name())), ".")
			found := false
			for _, e := range exts {
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
	buf.WriteString(buildTreeString(tree, "", true))
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

	// Copy to clipboard
	clipboard.WriteAll(buf.String())

	// Print summary
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

	// Print colored summary (with reset code at the end)
	fmt.Printf("%sTotal lines %d Total Characters %d\033[0m\n", colorCode, lineCount, charCount)
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
