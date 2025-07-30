package entity

type Embed struct {
	Box struct {
		Probability float64 `json:"probability"`
		XMax        int     `json:"x_max"`
		XMin        int     `json:"x_min"`
		YMax        int     `json:"y_max"`
		YMin        int     `json:"y_min"`
	} `json:"box"`
	Embedding []float64 `json:"embedding"`
}
