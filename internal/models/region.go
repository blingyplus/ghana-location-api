package models

type Region struct {
	ID       string  `json:"id"`
	CountryID string `json:"country_id"`
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	Capital  *string `json:"capital,omitempty"`
}
