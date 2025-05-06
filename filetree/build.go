package filetree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flowtemplates/flow-go/parser"
)

func parseNameAndCondition(raw string) (parser.AST, parser.Expr) {
	// FIXME:
	if strings.HasPrefix(raw, "{%") {
		end := strings.Index(raw, "%}")
		if end > 2 {
			// cond := raw[2:end]
			// name := raw[end+2:]
			return parser.AST{}, &parser.BadExpr{}
		}
	}

	return parser.AST{}, nil
}

func BuildFileTree(rootPath string) (*FileTree, error) {
	root, err := buildDir(rootPath)
	if err != nil {
		return nil, err
	}

	return (*FileTree)(root), nil
}

func buildDir(path string) (*Dir, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %w", path, err)
	}

	nameNode, condNode := parseNameAndCondition(filepath.Base(path))

	dir := &Dir{
		Path:      path,
		Name:      nameNode,
		Condition: condNode,
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		nameNode, condNode := parseNameAndCondition(entry.Name())

		if entry.IsDir() {
			subDir, err := buildDir(entryPath)
			if err != nil {
				return nil, err
			}

			dir.Dirs = append(dir.Dirs, *subDir)
		} else {
			dir.Files = append(dir.Files, File{
				Path:      entryPath,
				Name:      nameNode,
				Condition: condNode,
				Content:   nil,
			})
		}
	}

	return dir, nil
}
