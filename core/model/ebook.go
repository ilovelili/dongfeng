package model

// Ebook ebook entity
type Ebook struct {
	BaseModel
	Pupil     *Pupil
	PupilID   uint     `json:"pupil_id"`
	Date      string   `json:"date"`
	Hash      string   `json:"-"`
	Converted bool     `json:"-"`
	Images    []string `gorm:"-" json:"-"`
	HTML      string   `gorm:"-" json:"-"`
	CSS       string   `gorm:"-" json:"-"`
	Dates     []string `gorm:"-" json:"-"`
}
