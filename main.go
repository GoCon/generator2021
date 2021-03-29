package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v2"
)

const GitHubURLPrefix = "https://github.com/"

type Speaker struct {
	Name           string `csv:"name"`
	Avatar         string `csv:"avatar"`
	ID             string `csv:"id"`
	Location       string `csv:"location"`
	Bio            string `csv:"bio"`
	Twitter        string `csv:"twitter"`
	URL            string `csv:"url"`
	Organization   string `csv:"organization"`
	ShirtSize      string `csv:"shirt_size"`
	TalkFormat     string `csv:"talk_format"`
	Title          string `csv:"title"`
	Abstract       string `csv:"abstract"`
	Description    string `csv:"description"`
	Notes          string `csv:"notes"`
	AudienceLevel  string `csv:"audience_level"`
	Tags           string `csv:"tags"`
	Rating         string `csv:"rating"`
	State          string `csv:"state"`
	Confirmed      string `csv:"confirmed"`
	CreatedAt      string `csv:"created_at"`
	AdditionalInfo string `csv:"additional_info"`
}

type SpeakerContent struct {
	Key      string
	Name     string
	ID       string
	Company  string
	Feature  bool
	PhotoURL string `yaml:"photoURL"`
	Socials  []*Social
}

type Social struct {
	Icon string
	Link string
	Name string
}

func createSpeaker(s *Speaker) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(wd, "speakers", s.ID+".md"), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := &SpeakerContent{
		Key:      s.ID,
		Name:     s.Name,
		ID:       s.ID,
		Company:  s.Organization,
		Feature:  false,
		PhotoURL: fmt.Sprintf("/images/speakers/speaker-%s.jpg", s.ID),
	}

	if s.Twitter != "" {
		sc.Socials = append(sc.Socials, &Social{
			Icon: "twitter",
			Link: fmt.Sprintf("https://twitter.com/%s", s.Twitter),
			Name: s.Twitter,
		})
	}

	if s.URL != "" {
		icon := "link"
		name := "website"
		if strings.HasPrefix(s.URL, GitHubURLPrefix) {
			icon = "github"
			name = strings.TrimPrefix(s.URL, GitHubURLPrefix)
		}
		sc.Socials = append(sc.Socials, &Social{
			Icon: icon,
			Link: s.URL,
			Name: name,
		})
	}

	out, err := yaml.Marshal(sc)
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf(`---
%s---
%s`, string(out), s.Bio)

	_, err = f.WriteString(body)
	if err != nil {
		panic(err)
	}
}

func main() {
	f, err := os.Open("speakers.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var speakers []*Speaker
	if err := gocsv.UnmarshalFile(f, &speakers); err != nil { // Load clients from file
		panic(err)
	}
	for _, s := range speakers {
		//fmt.Println("%#v\n", speaker)
		if s.State != "accepted" {
			continue
		}
		createSpeaker(s)
	}
}
