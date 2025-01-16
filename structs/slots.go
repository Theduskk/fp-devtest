package structs

type Slot struct {
	ValidFrom string `json:"validFrom"`
	ValidTo   string `json:"validTo"`
	Carbon    Carbon `json:"carbon"`
}

type Carbon struct {
	Intensity int `json:"intensity"`
}

func ConvertIntensityItemsToSlot(i Intensity) Slot {
	return Slot{
		ValidFrom: i.From,
		ValidTo:   i.To,
		Carbon: Carbon{
			Intensity: *i.Intensity.Forecast,
		},
	}
}
