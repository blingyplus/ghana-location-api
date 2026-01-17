package models

type Constituency struct {
	ID         string  `json:"id"`
	DistrictID *string `json:"district_id,omitempty"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
}
