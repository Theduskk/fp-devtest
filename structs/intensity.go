package structs

type ResponseWrapper struct {
	Data []Intensity `json:"data"`
}

type Intensity struct {
	From      string           `json:"from"`
	To        string           `json:"to"`
	Intensity IntensityDetails `json:"intensity"`
}

type IntensityDetails struct {
	Forecast *int    `json:"forecast"`
	Actual   *int    `json:"actual"`
	Index    *string `json:"index"`
}
