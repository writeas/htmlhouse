package htmlhouse

import (
	"fmt"
)

type (
	PublicHouse struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		URL      string `json:"url"`
		ThumbURL string `json:"thumb_url"`
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

func (h *PublicHouse) process(app *app) {
	h.URL = fmt.Sprintf("%s/%s.html", app.cfg.HostName, h.ID)
	if h.ThumbURL != "" {
		h.ThumbURL = "https://peeper.html.house/" + h.ThumbURL
	}
}
