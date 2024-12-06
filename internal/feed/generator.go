package feed

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/feeds"

	"github.com/tinyzimmer/bg3mods-feed/internal/config"
	"github.com/tinyzimmer/bg3mods-feed/internal/mods"
)

// Generator is an interface for generating feeds of mods.
type Generator interface {
	// GetFeed generates a feed of mods based on the given options.
	GetFeed(context.Context, GeneratorOptions) (*Feed, error)
}

// Feed represents a feed of mods.
type Feed struct {
	// Content is the raw feed content.
	Content []byte
	// Format is the format of the feed.
	Format config.FeedFormat
	// SyncedAt is the time the feed was last syncedm
	SyncedAt time.Time
}

type generator struct {
	api      mods.Fetcher
	defaults GeneratorOptions

	cachedData    map[cacheKey]*cachedFeed
	cachedDataMux sync.Mutex
}

type cacheKey struct {
	maxItems int
	sort     string
	tags     string
	platform config.Platform
}

type cachedFeed struct {
	feed *feeds.Feed
	at   time.Time
}

// NewGenerator creates a new feed generator using the given fetcher and default options.
func NewGenerator(fetcher mods.Fetcher, defaults GeneratorOptions) Generator {
	return &generator{
		api:        fetcher,
		defaults:   defaults,
		cachedData: make(map[cacheKey]*cachedFeed),
	}
}

func (g *generator) GetFeed(ctx context.Context, overrides GeneratorOptions) (*Feed, error) {
	g.cachedDataMux.Lock()
	defer g.cachedDataMux.Unlock()

	opts := g.defaults.Merge(overrides)
	key := cacheKey{
		maxItems: opts.MaxItems,
		sort:     opts.GetSort(),
		tags:     strings.Join(opts.Tags, ","),
		platform: opts.Platform,
	}
	current := g.cachedData[key]
	if current == nil || time.Since(current.at) > opts.FetchInterval {
		feed, err := g.generate(ctx, opts)
		if err != nil {
			return nil, err
		}
		g.cachedData[key] = &cachedFeed{
			feed: feed,
			at:   time.Now().UTC(),
		}
		current = g.cachedData[key]
	} else {
		log.Println("Using cached feed data from", current.at)
	}

	var data string
	var err error
	switch opts.Format {
	case config.FormatRSS:
		data, err = current.feed.ToRss()
	case config.FormatAtom:
		data, err = current.feed.ToAtom()
	case config.FormatJSON:
		data, err = current.feed.ToJSON()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to render feed: %w", err)
	}
	return &Feed{
		Content:  []byte(data),
		Format:   opts.Format,
		SyncedAt: current.at,
	}, nil
}

func (g *generator) generate(ctx context.Context, opts GeneratorOptions) (*feeds.Feed, error) {
	limit := 100
	if opts.MaxItems > 0 && opts.MaxItems < limit {
		limit = opts.MaxItems
	}

	var offset int
	feed := &feeds.Feed{
		Title:       "BG3 Mods Feed",
		Link:        &feeds.Link{Href: ""},
		Description: "A feed of the latest mods for Baldur's Gate 3",
	}

	for {
		res, err := g.api.Fetch(ctx, mods.FetchOptions{
			Limit:  limit,
			Offset: offset,
			Tags:   opts.Tags,
			Sort:   opts.GetSort(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch mods: %w", err)
		}
		for _, mod := range res.Data {
			if opts.MaxItems > 0 && len(feed.Items) >= opts.MaxItems {
				return feed, nil
			}
			if opts.Platform.IsValid() && !mod.SupportsPlatform(opts.Platform) {
				continue
			}
			feed.Items = append(feed.Items, &feeds.Item{
				Id:          mod.NameID,
				Title:       mod.Name,
				Link:        &feeds.Link{Href: mod.ProfileURL},
				Description: mod.Summary,
				Created:     mod.DateAdded(),
				Updated:     mod.DateUpdated(),
				Content:     mod.Description,
			})
		}
		offset += limit
		if len(res.Data) < limit {
			return feed, nil
		}
	}
}
