package bookmark

import (
	"encoding/json"
	"os"
	"sort"
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
	sort.Strings(cats)
	return cats
}

func AddCategory(bm Bookmarks, name string) {
	bm[name] = []Bookmark{}
}

func DeleteCategory(bm Bookmarks, name string) {
	delete(bm, name)
}

func AddBookmark(bm Bookmarks, category string, b Bookmark) {
	bm[category] = append(bm[category], b)
}

func DeleteBookmark(bm Bookmarks, category string, index int) {
	items := bm[category]
	bm[category] = append(items[:index], items[index+1:]...)
}

func UpdateBookmark(bm Bookmarks, category string, index int, b Bookmark) {
	bm[category][index] = b
}
