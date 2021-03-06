// Package atom provides an Atom plugin that can be used to generate
// Atom RSS feeds.
package atom

import (
	"fmt"
	"github.com/dominikbraun/espresso/filesystem"
	"github.com/dominikbraun/espresso/model"
	"github.com/dominikbraun/espresso/render"
	"github.com/gorilla/feeds"
	"path/filepath"
	"time"
)

// filename is the default filename for the Atom XML.
const filename string = "atom.xml"

// Meta holds all metadata for the RSS feed and should be populated
// using the values from site.yml automatically.
type Meta struct {
	Title       string
	BaseURL     string
	Description string
	Author      string
	Subtitle    string
	Copyright   string
}

// atom stores all metadata and generates the RSS feed.
type atom struct {
	meta *Meta
	feed *feeds.Feed
}

// New creates a new instance of the Atom plugin. It is fully initialized
// using the provided metadata and thus ready to use.
func New(meta *Meta) *atom {
	a := atom{
		meta: meta,
		feed: &feeds.Feed{
			Title:       meta.Title,
			Link:        &feeds.Link{Href: meta.BaseURL},
			Description: meta.Description,
			Author:      &feeds.Author{Name: meta.Author},
			Created:     time.Now(),
			Subtitle:    meta.Subtitle,
			Items:       make([]*feeds.Item, 0),
			Copyright:   meta.Copyright,
		},
	}
	return &a
}

// ProcessArticlePage implements render.Plugin.ProcessArticlePage and
// generates and adds a new RSS feed entry based on the article data.
func (a *atom) ProcessArticlePage(_ *render.Context, page *model.ArticlePage) error {
	if page.Article.Hide {
		return nil
	}

	absoluteURL := fmt.Sprintf("%s%s%s", a.meta.BaseURL, page.Path, page.Article.ID)

	item := &feeds.Item{
		Title:       page.Article.Title,
		Link:        &feeds.Link{Href: absoluteURL},
		Description: page.Article.Description,
		Id:          absoluteURL,
		Created:     page.Article.Date,
	}
	a.feed.Items = append(a.feed.Items, item)

	return nil
}

// Finalize implements render.Plugin.Finalize and writes the RSS feed XML
// file to the target directory stored in ctx.
func (a *atom) Finalize(ctx *render.Context) error {
	filePath := filepath.Join(ctx.TargetDir, filename)
	atomFile, err := filesystem.CreateFile(filePath)
	if err != nil {
		return err
	}

	return a.feed.WriteAtom(atomFile)
}
