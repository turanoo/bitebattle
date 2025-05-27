package restaurant

type Place struct {
	Name    string  `json:"name"`
	Address string  `json:"address"`
	PlaceID string  `json:"place_id"`
	Rating  float64 `json:"rating,omitempty"`
	Photos  []Photo `json:"photos,omitempty"`
}

type Photo struct {
	PhotoReference string `json:"photo_reference"`
}
