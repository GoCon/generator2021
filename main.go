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
	SessionID      string `csv:"session_id"`
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

type SessionContent struct {
	Key          string
	Title        string
	ID           string
	Format       string
	TalkType     string `yaml:"talkType"`
	Level        string
	Tags         []string
	Speakers     []string
	VideoID      *string `yaml:"videoId"`
	Presentation *string
	Draft        bool
}

func createSpeaker(s *Speaker) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(wd, "speakers", fmt.Sprintf("s.md", s.ID)), os.O_WRONLY|os.O_CREATE, os.ModePerm)
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

func createSession(s *Speaker) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(filepath.Join(wd, "sessions", fmt.Sprintf("session-%s.md", s.SessionID)), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var talkType string
	switch s.TalkFormat {
	case "Lightning Talk: 5 minutes":
		talkType = "lightning_talk"
	default:
		// 他はやらない
		panic(s.TalkFormat)
	}

	sc := &SessionContent{
		Key:      s.SessionID,
		Title:    s.Title,
		ID:       s.SessionID,
		Format:   "conference",
		TalkType: talkType,
		Level:    strings.ToLower(s.AudienceLevel),       // e.g. `Advanced` => `advanced`
		Tags:     []string{strings.ToUpper(s.SessionID)}, // e.g. `lt1` => `LT1`
		Speakers: []string{s.ID},
		Draft:    false,
	}

	out, err := yaml.Marshal(sc)
	if err != nil {
		panic(err)
	}
	body := fmt.Sprintf(`---
%s---
%s
---
%s`, string(out), s.Abstract, s.Description)

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
		if s.State != "accepted" {
			continue
		}
		createSpeaker(s)

		if s.SessionID == "" {
			continue
		}
		createSession(s)
	}
}
