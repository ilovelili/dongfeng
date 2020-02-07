package model

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

// Ebook ebook entity
type Ebook struct {
	BaseModel
	Pupil     *Pupil   `json:"pupil"`
	PupilID   uint     `gorm:"unique_index:idx_pupil_date" json:"pupil_id"`
	Date      string   `gorm:"unique_index:idx_pupil_date" json:"date"`
	Hash      string   `gorm:"unique_index" json:"-"`
	Converted bool     `json:"-"`
	Images    []string `gorm:"-" json:"images"`
	HTML      string   `gorm:"-" json:"html"`
	CSS       string   `gorm:"-" json:"css"`
}

// ResolveHash resolve content md5 hash
func (e *Ebook) ResolveHash() {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(int(e.PupilID)))
	sb.WriteString(e.Date)
	for _, img := range e.Images {
		sb.WriteString(img)
	}
	sb.WriteString(e.HTML)
	sb.WriteString(e.CSS)

	str := sb.String()
	hash := md5.Sum([]byte(str))
	e.Hash = hex.EncodeToString(hash[:])
}

// ResolveCloudCSS replace image link
func (e *Ebook) ResolveCloudCSS() string {
	return strings.Replace(e.CSS, "../img/", "../../../../../../img/", -1)
}

// ResolveCloudHTML replace style link with oss css
func (e *Ebook) ResolveCloudHTML() string {
	return strings.Replace(e.HTML, "./img/", "../../../../../img/", -1)
}
