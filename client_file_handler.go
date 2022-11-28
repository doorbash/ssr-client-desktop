package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type ClientFileHandler struct {
	ctx      context.Context
	baseurl  string
	filename string
	dir      string
}

func (b *ClientFileHandler) Sha256FileExists() bool {
	stats, err := os.Stat(fmt.Sprintf("%s/%s.sha256", b.dir, b.filename))
	if err == nil {
		if stats.IsDir() {
			return false
		}
		if stats.Size() != 64 {
			return false
		}
		return true
	}
	return false
}

func (b *ClientFileHandler) ClientFileExists() bool {
	if !b.Sha256FileExists() {
		return false
	}
	stats, err := os.Stat(fmt.Sprintf("%s/%s", b.dir, b.filename))
	if err == nil {
		if stats.IsDir() {
			return false
		}
		if stats.Size() > 10*1024*1024 {
			return false
		}
		f, err := os.Open(fmt.Sprintf("%s/%s.sha256", b.dir, b.filename))
		if err != nil {
			log.Println(err)
			return false
		}
		defer f.Close()
		buf := make([]byte, 64)
		n, err := f.Read(buf)
		if err != nil || n != 64 {
			return false
		}
		c1 := string(buf)
		log.Println("checksum from .sha256 file: ", c1)
		f, err = os.Open(fmt.Sprintf("%s/%s", b.dir, b.filename))
		if err != nil {
			log.Println(err)
			return false
		}
		defer f.Close()
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Println(err)
			return false
		}
		c2 := fmt.Sprintf("%x", h.Sum(nil))
		log.Println("checksum from binary file: ", c2)
		return c1 == c2
	}
	return false
}

func (b *ClientFileHandler) DeleteSha256File() error {
	return os.Remove(fmt.Sprintf("%s/%s.sha256", b.dir, b.filename))
}

func (b *ClientFileHandler) DeleteClientFile() error {
	return os.Remove(fmt.Sprintf("%s/%s", b.dir, b.filename))
}

func (b *ClientFileHandler) downloadSha256File() error {
	resp, err := http.Get(fmt.Sprintf("%s/%s.sha256", b.baseurl, b.filename))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/%s.sha256", b.dir, b.filename))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	if !b.Sha256FileExists() {
		return errors.New("error while downloading sha256 file")
	}

	return nil
}

func (b *ClientFileHandler) DownloadClientFile(c context.Context, pf func(progress int)) error {
	b.DeleteSha256File()
	b.DeleteClientFile()

	err := b.downloadSha256File()
	if err != nil {
		return err
	}

	out, err := os.Create(fmt.Sprintf("%s/%s", b.dir, b.filename))
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(fmt.Sprintf("%s/%s", b.baseurl, b.filename))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cl := resp.Header.Get("Content-Length")
	if cl == "" {
		return errors.New("content length is empty")
	}
	size, err := strconv.ParseUint(cl, 10, 0)
	if err != nil {
		return err
	}
	sizeF := float64(size)
	reader := NewReaderCtx(c, resp.Body, func(total uint64) {
		progress := math.Round((float64(total * 100)) / sizeF)
		pf(int(progress))
	})
	if _, err = io.Copy(out, reader); err != nil {
		return err
	}
	return nil
}

func NewClientFileHandler(ctx context.Context, baseurl, name, dir string) *ClientFileHandler {
	return &ClientFileHandler{
		baseurl:  baseurl,
		filename: name,
		dir:      dir,
	}
}
