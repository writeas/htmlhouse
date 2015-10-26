package htmlhouse

type (
	HouseStats struct {
		ID    string
		Stats []Stat
	}

	Stat struct {
		Data  string
		Label string
	}
)
