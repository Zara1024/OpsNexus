package model

// GlobalSearchQuery defines the supported search filters.
type GlobalSearchQuery struct {
	Keyword string `form:"keyword"`
	Types   string `form:"types"`
	Limit   int    `form:"limit"`
}

// GlobalSearchResult is one cross-resource search item.
type GlobalSearchResult struct {
	Type        string   `json:"type"`
	TypeLabel   string   `json:"typeLabel"`
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Route       string   `json:"route"`
	Tags        []string `json:"tags"`
	AccessPath  string   `json:"-"`
}

// GlobalSearchGroup groups search items by resource type.
type GlobalSearchGroup struct {
	Type      string               `json:"type"`
	TypeLabel string               `json:"typeLabel"`
	Count     int                  `json:"count"`
	Items     []GlobalSearchResult `json:"items"`
}

// GlobalSearchResponse is the top-level search response.
type GlobalSearchResponse struct {
	Keyword string               `json:"keyword"`
	Total   int                  `json:"total"`
	Groups  []GlobalSearchGroup  `json:"groups"`
	Results []GlobalSearchResult `json:"results"`
}
