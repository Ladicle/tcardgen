package hugo

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/gohugoio/hugo/parser/pageparser"
)

type FrontMatter struct {
	Title    string
	Author   string
	Category string
	Tags     []string
	LastMod  time.Time
}

// ParseFrontMatter parses the frontmatter of the specified Hugo content.
func ParseFrontMatter(filename string) (*FrontMatter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfm, err := pageparser.ParseFrontMatterAndContent(file)
	if err != nil {
		return nil, err
	}

	var ok bool
	fm := &FrontMatter{}

	if fm.Title, ok = cfm.FrontMatter["title"].(string); !ok {
		return nil, fmt.Errorf("can not convert title to string: %+v", cfm.FrontMatter)
	} else if fm.Title == "" {
		return nil, fmt.Errorf("title is empty: %+v", cfm.FrontMatter)
	}

	if fm.Author, err = getFirstFMItem(cfm, "author"); err != nil {
		return nil, err
	}

	if fm.Category, err = getFirstFMItem(cfm, "categories"); err != nil {
		return nil, err
	}

	tags, ok := cfm.FrontMatter["tags"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("can not convert tags to interface{} array: %+v", cfm.FrontMatter)
	} else if len(tags) == 0 {
		return nil, fmt.Errorf("there is no tags: %+v", tags)
	}
	for _, t := range tags {
		tag := t.(string)
		if !isUpper(tag) {
			tag = strings.Title(tag)
		}
		fm.Tags = append(fm.Tags, tag)
	}

	fm.LastMod = cfm.FrontMatter["lastmod"].(time.Time)

	return fm, nil
}

func getFirstFMItem(cfm pageparser.ContentFrontMatter, key string) (string, error) {
	categoriesitems, ok := cfm.FrontMatter[key].([]interface{})
	if !ok {
		return "", fmt.Errorf("can not convert %s to interface{} array: %+v", key, cfm.FrontMatter)
	}
	if len(categoriesitems) < 1 {
		return "", fmt.Errorf("can not get %s from front matter: %+v", key, cfm.FrontMatter)
	}
	return categoriesitems[0].(string), nil
}

func isUpper(text string) bool {
	for _, r := range []rune(text) {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
