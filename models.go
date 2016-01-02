package htmlhouse

type (
	PublicHouse struct {
		ID       string
		Title    string
		ThumbURL string
	}

	HouseStats struct {
		ID    string
		Stats []Stat
	}

	Stat struct {
		Data  string
		Label string
	}
)
