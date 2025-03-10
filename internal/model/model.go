package model

type Song struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Group       string `json:"group"`
	SongTitle   string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
