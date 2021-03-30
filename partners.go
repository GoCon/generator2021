package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

type Partner struct {
	Key         string `csv:"key"`
	Title       string `csv:"title"`
	Category    string `csv:"category"`
	Order       int    `csv:"order"`
	LogoExists  string `csv:"logo_exists"`
	Description string `csv:"description"`
}

type PartnerContent struct {
	Key      string
	Title    string
	Category string
	Order    int
	Logo     string
	Lang     string
}

func createPartner(p *Partner) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dirPath := filepath.Join(wd, "partners", p.Category)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(dirPath, fmt.Sprintf("%s.md", p.Key)), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := &PartnerContent{
		Key:      p.Key,
		Title:    p.Title,
		Category: p.Category,
		Order:    p.Order,
		Logo:     fmt.Sprintf("/images/partners/logo-%s.png", p.Key),
		Lang:     "ja",
	}

	out, err := yaml.Marshal(sc)
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf(`---
%s---
%s`, string(out), p.Description)
	rd := strings.NewReader(body)

	_, err = io.Copy(f, rd)
	if err != nil {
		panic(err)
	}
}

func processPartners() {
	f, err := os.Open("partners.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var partners []*Partner
	if err := gocsv.UnmarshalFile(f, &partners); err != nil { // Load clients from file
		panic(err)
	}
	for _, p := range partners {
		if p.LogoExists != "yes" {
			continue
		}
		createPartner(p)
	}
}
