package htmlhouse

import (
	"fmt"
	"time"
)

type (
	PublicHouse struct {
		ID       string    `json:"id"`
		Title    string    `json:"title"`
		URL      string    `json:"url"`
		ThumbURL string    `json:"thumb_url"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
		Views    int       `json:"views"`
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
