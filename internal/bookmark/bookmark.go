package bookmark

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Bookmark struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}

type Bookmarks map[string][]Bookmark

func Load(path string) (Bookmarks, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Bookmarks{}, nil
		}
		return nil, err
	}

	var bm Bookmarks
	if err := json.Unmarshal(data, &bm); err != nil {
		return nil, err
	}
	for _, items := range bm {
		sortBookmarks(items)
	}
	return bm, nil
}

func Save(path string, bm Bookmarks) error {
	data, err := json.MarshalIndent(bm, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Categories(bm Bookmarks) []string {
	cats := make([]string, 0, len(bm))
	for k := range bm {
		cats = append(cats, k)
	}
	sort.Slice(cats, func(i, j int) bool {
		return strings.ToLower(cats[i]) < strings.ToLower(cats[j])
	})
	return cats
}

func AddCategory(bm Bookmarks, name string) {
	bm[name] = []Bookmark{}
}

func DeleteCategory(bm Bookmarks, name string) {
	delete(bm, name)
}

func RenameCategory(bm Bookmarks, oldName, newName string) {
	if oldName == newName {
		return
	}
	bm[newName] = bm[oldName]
	delete(bm, oldName)
}

func AddBookmark(bm Bookmarks, category string, b Bookmark) {
	bm[category] = append(bm[category], b)
	sortBookmarks(bm[category])
}

func DeleteBookmark(bm Bookmarks, category string, index int) {
	items := bm[category]
	bm[category] = append(items[:index], items[index+1:]...)
}

func UpdateBookmark(bm Bookmarks, category string, index int, b Bookmark) {
	bm[category][index] = b
	sortBookmarks(bm[category])
}

func sortBookmarks(items []Bookmark) {
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})
}

func Import(path string, bm Bookmarks) error {
	imported, err := Load(path)
	if err != nil {
		return err
	}
	for cat, items := range imported {
		if existing, ok := bm[cat]; ok {
			bm[cat] = append(existing, items...)
			sortBookmarks(bm[cat])
		} else {
			bm[cat] = items
			sortBookmarks(bm[cat])
		}
	}
	return nil
}

// Param represents a named parameter extracted from a command template.
type Param struct {
	Name    string
	Default string
}

var ParamRegex = regexp.MustCompile(`\{\{(\w+)(?::([^}]*))?\}\}`)

// ParseParams extracts parameters from a command string containing {{name}} or
// {{name:default}} placeholders. Returns unique params in order of first appearance.
func ParseParams(cmd string) []Param {
	matches := ParamRegex.FindAllStringSubmatch(cmd, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var params []Param
	for _, m := range matches {
		name := m[1]
		if seen[name] {
			continue
		}
		seen[name] = true
		params = append(params, Param{Name: name, Default: m[2]})
	}
	return params
}

// ResolveParams replaces all {{name}} and {{name:default}} placeholders in cmd
// with the corresponding values from the provided map.
func ResolveParams(cmd string, values map[string]string) string {
	return ParamRegex.ReplaceAllStringFunc(cmd, func(match string) string {
		sub := ParamRegex.FindStringSubmatch(match)
		if sub == nil {
			return match
		}
		name := sub[1]
		if val, ok := values[name]; ok {
			return val
		}
		return match
	})
}

func Export(bm Bookmarks) (string, error) {
	filename := fmt.Sprintf("bookmarks-backup-%s.json", time.Now().Format("2006-01-02-150405"))
	err := Save(filename, bm)
	if err != nil {
		return "", err
	}
	return filename, nil
}
