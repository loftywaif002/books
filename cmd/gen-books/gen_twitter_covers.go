package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kjk/u"
)

func loadImageMust(path string) image.Image {
	r, err := os.Open(path)
	u.PanicIfErr(err)
	defer r.Close()
	img, _, err := image.Decode(r)
	u.PanicIfErr(err)
	return img
}

func saveImageMust(path string, img image.Image) {
	w, err := os.Create(path)
	u.PanicIfErr(err)
	defer func() {
		err = w.Close()
		u.PanicIfErr(err)
	}()
	err = png.Encode(w, img)
	u.PanicIfErr(err)
}

func optiImageMust(path string) {
	cmd := exec.Command("optipng", "-o7", path)
	_, err := cmd.CombinedOutput()
	u.PanicIfErr(err)
}

func saveTwitterImage(path string, img image.Image) {
	saveImageMust(path, img)
	optiImageMust(path)
	fmt.Printf("Saved and optimized '%s'\n", path)
}

func getSubimage(img image.Image, r image.Rectangle) image.Image {
	switch im := img.(type) {
	case *image.NRGBA:
		return im.SubImage(r)
	case *image.Paletted:
		return im.SubImage(r)
	}
	u.PanicIf(true, "unsupported image type %T", img)
	return img
}

func genTwitterImage(img image.Image) image.Image {
	r := img.Bounds()
	dx := r.Dx()
	r = image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{dx, dx},
	}
	return getSubimage(img, r)
}

func printImageInfo(path string, img image.Image) {
	r := img.Bounds()
	fmt.Printf("%s: %v %T\n", path, r, img)
}

func getExistingTwitterImagesMust(dir string) map[string]bool {
	m := make(map[string]bool)
	fileInfos, err := ioutil.ReadDir(dir)
	u.PanicIfErr(err)
	for _, fi := range fileInfos {
		name := fi.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext == ".png" {
			m[name] = true
		}
	}
	return m
}

func getCoversListMust(dir string) []string {
	var res []string
	fileInfos, err := ioutil.ReadDir(dir)
	u.PanicIfErr(err)
	for _, fi := range fileInfos {
		name := fi.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".png" {
			continue
		}
		if strings.Contains(name, "@2x") {
			continue
		}
		res = append(res, name)
	}
	return res
}

func genTwitterImagesAndExit() {
	srcDir := "covers"
	dstDir := filepath.Join(srcDir, "twitter")
	createDirMust(dstDir)
	covers := getCoversListMust(srcDir)
	existingTwitter := getExistingTwitterImagesMust(dstDir)
	// fmt.Printf("covers: %v\n", covers)
	for _, coverName := range covers {
		dstPath := filepath.Join(dstDir, coverName)
		if _, ok := existingTwitter[coverName]; false && ok {
			fmt.Printf("%s already exists as %s\n", coverName, dstPath)
			continue
		}
		path := filepath.Join(srcDir, coverName)
		img := loadImageMust(path)
		printImageInfo(path, img)
		sub := genTwitterImage(img)
		saveTwitterImage(dstPath, sub)
	}
	os.Exit(0)
}
