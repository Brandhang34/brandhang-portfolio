package handler

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

//go:embed portfolio.json
var portfolioJSON []byte

const imgPath = "/assets/imgs/portfolio_imgs/"

// Portfolio is a single portfolio entry: a project, certification or wiki page.
type Portfolio struct {
	Title       string   `json:"title"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
}

// ImagePath returns the servable path for the entry's image, or "" when the
// entry has no image (the template renders a placeholder instead).
func (p Portfolio) ImagePath() string {
	if p.Image == "" {
		return ""
	}
	return imgPath + p.Image
}

// filterTags maps the values used by the portfolio page's <select> to the tag
// strings stored on each entry.
var filterTags = map[string]string{
	"projects":       "Project",
	"certifications": "Certification",
	"homelab":        "Homelab",
}

var portfolioList = mustLoadPortfolio()

func mustLoadPortfolio() []Portfolio {
	items, err := loadPortfolio(portfolioJSON)
	if err != nil {
		panic(fmt.Sprintf("loading portfolio.json: %v", err))
	}
	return items
}

func loadPortfolio(data []byte) ([]Portfolio, error) {
	var items []Portfolio
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	for i, item := range items {
		if item.Title == "" {
			return nil, fmt.Errorf("entry %d: title is required", i)
		}
	}
	return items, nil
}

// LoadAllPortfolioItems returns every portfolio entry.
func LoadAllPortfolioItems() []Portfolio {
	return portfolioList
}

// GetByTag returns the entries carrying the given tag.
func GetByTag(tag string) []Portfolio {
	tagged := []Portfolio{}
	for _, item := range portfolioList {
		if slices.Contains(item.Tags, tag) {
			tagged = append(tagged, item)
		}
	}
	return tagged
}

// SearchPortfolio filters by tag and then by a case-insensitive substring match
// on the title, description and tags.
func SearchPortfolio(searchQuery, tagsQuery string) []Portfolio {
	list := portfolioList
	if tag, ok := filterTags[tagsQuery]; ok {
		list = GetByTag(tag)
	}

	searchQuery = strings.ToLower(strings.TrimSpace(searchQuery))
	if searchQuery == "" {
		return list
	}

	results := []Portfolio{}
	for _, item := range list {
		if matches(item, searchQuery) {
			results = append(results, item)
		}
	}
	return results
}

func matches(item Portfolio, query string) bool {
	if strings.Contains(strings.ToLower(item.Title), query) ||
		strings.Contains(strings.ToLower(item.Description), query) {
		return true
	}
	return slices.ContainsFunc(item.Tags, func(tag string) bool {
		return strings.Contains(strings.ToLower(tag), query)
	})
}
