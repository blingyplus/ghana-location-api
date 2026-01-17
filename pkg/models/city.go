package models

type City struct {
	ID         string   `json:"id"`
	DistrictID string   `json:"district_id"`
	Name       string   `json:"name"`
	Lat        *float64 `json:"lat,omitempty"`
	Lng        *float64 `json:"lng,omitempty"`
}
