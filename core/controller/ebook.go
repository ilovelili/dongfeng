package controller

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	oss "github.com/ilovelili/aliyun-client/oss"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/core/repository"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

// Ebook controller
type Ebook struct {
	ebookRepo *repository.Ebook
	svc       *oss.Service
}

// NewEbookController constructor
func NewEbookController() *Ebook {
	return &Ebook{
		ebookRepo: repository.NewEbookRepository(),
		svc: func() *oss.Service {
			_svc := oss.NewService(config.OSS.APIKey, config.OSS.APISecret)
			_svc.SetEndPoint(config.OSS.Endpoint)
			_svc.SetBucket(config.OSS.BucketName)
			return _svc
		}(),
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

		// delegate convert & merge to dongfeng jobs batch
		// // convert to pdf
		// if err = c.convert(ebook); err != nil {
		// 	return err
		// }

		// // merge to ebook
		// if err = c.merge(ebook); err != nil {
		// 	return err
		// }

		// set converted to true if everything goes smoothly
		// ebook.Converted = true
		// _, err = c.ebookRepo.Save(ebook, true)
		// if err != nil {
		// 	return err
		// }
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

// convert ebook to pdf and jpg using chrome headless
func (c *Ebook) convert(ebook *model.Ebook) (err error) {
	year, class, name, date := ebook.Pupil.Class.Year, ebook.Pupil.Class.Name, ebook.Pupil.Name, ebook.Date
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	chromeDevTool := os.Getenv("CHROME_DEV_TOOL")
	if chromeDevTool == "" {
		chromeDevTool = "http://127.0.0.1:9222"
	}

	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New(chromeDevTool)
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return
		}
	}
	defer devt.Close(ctx, pt)

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err := rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return
	}
	defer conn.Close() // Leaving connections open will leak memory.

	cli := cdp.NewClient(conn)
	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := cli.Page.DOMContentEventFired(ctx)
	if err != nil {
		return
	}
	defer domContent.Close()

	// Enable the runtime
	if err = cli.Runtime.Enable(ctx); err != nil {
		return
	}

	// Enable the network
	if err = cli.Network.Enable(ctx, network.NewEnableArgs()); err != nil {
		return
	}

	// Enable events on the Page domain, it's often preferrable to create
	// event clients before enabling events so that we don't miss any.
	if err = cli.Page.Enable(ctx); err != nil {
		return
	}

	htmllocaldir := path.Join(config.Ebook.OriginDir, year, class, name, date)
	// Create the Navigate arguments
	navArgs := page.NewNavigateArgs(fmt.Sprintf("file://%s", path.Join(htmllocaldir, "index.html")))
	nav, err := cli.Page.Navigate(ctx, navArgs)
	if err != nil {
		return
	}

	// wait till image loaded
	time.Sleep(time.Duration(config.Ebook.ImageLoadTimeout) * time.Second)

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return
	}

	fmt.Printf("Page loaded with frame ID: %s\n", nav.FrameID)

	imgOutput := path.Join(htmllocaldir, "output.jpg")
	// Capture a screenshot of the current page.
	screenshotArgs := page.NewCaptureScreenshotArgs().
		SetFormat("jpeg").
		SetQuality(100)

	screenshot, err := cli.Page.CaptureScreenshot(ctx, screenshotArgs)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(imgOutput, screenshot.Data, 0644); err != nil {
		return
	}

	fmt.Printf("Saved screenshot: %s\n", imgOutput)

	// Print to PDF
	printToPDFArgs := page.NewPrintToPDFArgs().
		SetLandscape(false).
		SetPrintBackground(true).
		SetMarginTop(0).
		SetMarginBottom(0).
		SetMarginLeft(0).
		SetMarginRight(0).
		SetPaperWidth(config.Ebook.Width).
		SetPaperHeight(config.Ebook.Height)

	print, _ := cli.Page.PrintToPDF(ctx, printToPDFArgs)
	pdfOutput := path.Join(htmllocaldir, "output.pdf")
	if err = ioutil.WriteFile(pdfOutput, print.Data, 0644); err != nil {
		return
	}

	fmt.Printf("Saved pdf: %s\n", pdfOutput)

	// move to dest dir
	pdfdestdir := path.Join(config.Ebook.PDFDestDir, year, class, name)
	_, err = os.Stat(pdfdestdir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(pdfdestdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	newPdfOutput := path.Join(pdfdestdir, fmt.Sprintf("%s.pdf", ebook.Date))
	fmt.Printf("moving pdf from %s to %s\n", pdfOutput, newPdfOutput)
	if err = os.Rename(pdfOutput, newPdfOutput); err != nil {
		return err
	}

	imgdestdir := path.Join(config.Ebook.ImageDestDir, year, class, name)
	_, err = os.Stat(imgdestdir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(imgdestdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	newImgOutput := path.Join(imgdestdir, fmt.Sprintf("%s.jpg", ebook.Date))
	fmt.Printf("moving img from %s to %s\n", pdfOutput, newPdfOutput)
	err = os.Rename(imgOutput, newImgOutput)
	return err
}

// merge merge pdf files into ebook
func (c *Ebook) merge(ebook *model.Ebook) (err error) {
	year, class, name := ebook.Pupil.Class.Year, ebook.Pupil.Class.Name, ebook.Pupil.Name
	// check if pdftk installed or not
	_, err = exec.LookPath("pdftk")
	if err != nil {
		return
	}

	filepathmap := make(map[string][]string)
	targetdir := config.Ebook.MergeTargetDir
	destdir := path.Join(config.Ebook.MergeDestDir, class, name)

	err = filepath.Walk(targetdir, func(filepath string, info os.FileInfo, err error) error {
		// target
		if !info.IsDir() && path.Ext(info.Name()) == ".pdf" {
			key := path.Dir(filepath)
			// select the target file with corresponding class and name
			if strings.Index(key, fmt.Sprintf("/%s/%s", class, name)) == -1 {
				return nil
			}

			// ignore the dest file
			if strings.Index(key, config.Ebook.MergeDestDir) > -1 {
				return nil
			}

			if paths, ok := filepathmap[key]; ok {
				filepathmap[key] = append(paths, filepath)
			} else {
				filepathmap[key] = []string{filepath}
			}
		}
		return nil
	})

	if err != nil {
		return
	}

	// first clear dest dir
	os.RemoveAll(destdir)
	_, err = os.Stat(destdir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(destdir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for dir, filepaths := range filepathmap {
		// sort pdf by date
		sort.Strings(filepaths)
		// https://stackoverflow.com/questions/31467153/golang-failed-exec-command-that-works-in-terminal
		// cmdline := fmt.Sprintf("pdftk %s cat output merge.pdf", path.Join(filepath, "*.pdf"))
		pdffiles := strings.Join(filepaths, " ")
		cmdline := fmt.Sprintf("pdftk %s cat output %s", pdffiles, path.Join(dir, "merge.pdf"))
		args := strings.Split(cmdline, " ")
		cmd := exec.Command(args[0], args[1:]...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}

		// move to dest
		// 电子书_${this.currentName}_${this.currentClass}_${this.currentYear}学年.pdf
		err = os.Rename(path.Join(dir, "merge.pdf"), path.Join(destdir, fmt.Sprintf("电子书_%s_%s_%s学年.pdf", name, class, year)))
		if err != nil {
			return
		}
	}

	return
}
