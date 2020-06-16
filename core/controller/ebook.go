package controller

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/core/repository"
)

// Ebook controller
type Ebook struct {
	ebookRepo *repository.Ebook
}

// NewEbookController constructor
func NewEbookController() *Ebook {
	return &Ebook{
		ebookRepo: repository.NewEbookRepository(),
	}
}

// SaveEbook save ebook
func (c *Ebook) SaveEbook(ebook *model.Ebook) error {
	ebook.ResolveHash()
	ebook.Converted = false
	dirty, err := c.ebookRepo.Save(ebook, false)
	if err != nil {
		return err
	}

	_ebook, err := c.ebookRepo.FindByID(ebook.ID)
	if err != nil {
		return err
	}
	ebook.Pupil = _ebook.Pupil

	// if dirty
	if dirty {
		// upload to storage
		if err = c.uploadToStorage(ebook); err != nil {
			return err
		}
	}

	return nil
}

// RemoveFromStorage remove from storage
func (c *Ebook) RemoveFromStorage(ebook *model.Ebook) error {
	htmllocaldir := path.Join(config.Ebook.OriginDir, ebook.Pupil.Class.Year, ebook.Pupil.Class.Name, ebook.Pupil.Name, ebook.Date)
	return os.RemoveAll(htmllocaldir)
}

// uploadToCloudStorage upload css folder and index.html to local (or aliyun oss later)
func (c *Ebook) uploadToStorage(ebook *model.Ebook) error {
	htmllocaldir := path.Join(config.Ebook.OriginDir, ebook.Pupil.Class.Year, ebook.Pupil.Class.Name, ebook.Pupil.Name, ebook.Date)
	csslocaldir := path.Join(htmllocaldir, "css")

	_, err := os.Stat(csslocaldir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(csslocaldir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	csslocalfile := path.Join(csslocaldir, "style.css")
	err = ioutil.WriteFile(csslocalfile, []byte(ebook.ResolveCloudCSS()), os.ModePerm)
	if err != nil {
		return err
	}

	htmllocalfile := path.Join(htmllocaldir, "index.html")
	err = ioutil.WriteFile(htmllocalfile, []byte(ebook.ResolveCloudHTML()), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// SaveTemplatePreview save template preview
func (c *Ebook) SaveTemplatePreview(templatePreview *model.TemplatePreview) error {
	templatePreview.ResolveHash()
	templatePreview.Converted = false
	dirty, err := c.ebookRepo.SaveTemplatePreview(templatePreview, false)
	if err != nil {
		return err
	}

	// if dirty
	if dirty {
		// upload to storage
		if err = c.uploadTemplatePreviewToStorage(templatePreview); err != nil {
			return err
		}
	}

	return nil
}

// RemoveTemplatePreviewFromStorage remove from storage
func (c *Ebook) RemoveTemplatePreviewFromStorage(templatePreview *model.TemplatePreview) error {
	htmllocaldir := path.Join(config.Ebook.OriginDir, "templatePreview", templatePreview.Name)
	return os.RemoveAll(htmllocaldir)
}

// uploadTemplatePreviewToStorage upload css folder and index.html to local (or aliyun oss later)
func (c *Ebook) uploadTemplatePreviewToStorage(templatePreview *model.TemplatePreview) error {
	htmllocaldir := path.Join(config.Ebook.OriginDir, "templatePreview", templatePreview.Name)
	csslocaldir := path.Join(htmllocaldir, "css")

	_, err := os.Stat(csslocaldir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(csslocaldir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	csslocalfile := path.Join(csslocaldir, "style.css")
	err = ioutil.WriteFile(csslocalfile, []byte(templatePreview.ResolveCloudCSS()), os.ModePerm)
	if err != nil {
		return err
	}

	htmllocalfile := path.Join(htmllocaldir, "index.html")
	err = ioutil.WriteFile(htmllocalfile, []byte(templatePreview.ResolveCloudHTML()), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
