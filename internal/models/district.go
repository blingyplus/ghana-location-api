package models

type District struct {
	ID       string  `json:"id"`
	RegionID string `json:"region_id"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Type     string  `json:"type"` // metro, municipal, district
	Capital  *string `json:"capital,omitempty"`
}
