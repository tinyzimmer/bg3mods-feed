package feed

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tinyzimmer/bg3mods-feed/internal/config"
)

// GeneratorOptions are the options for generating a feed.
type GeneratorOptions struct {
	// MaxItems is the maximum number of items to include in the feed.
	MaxItems int
	// Sort is the field to sort the feed by.
	Sort string
	// Tags are the tags to filter the feed by.
	Tags []string
	// Platform is the platform to filter the feed by.
	Platform config.Platform
	// FetchInterval is the interval to fetch mods at.
	FetchInterval time.Duration
	// Format is the format to render the feed in.
	Format config.FeedFormat
}

// OptionsFromQuery parses the query parameters from a URL into a GeneratorOptions struct.
func OptionsFromQuery(u *url.URL) GeneratorOptions {
	opts := GeneratorOptions{}
	if maxItems, err := strconv.Atoi(u.Query().Get("max_items")); err == nil {
		opts.MaxItems = maxItems
	}
	if sort := u.Query().Get("sort"); sort != "" {
		opts.Sort = sort
	}
	if tags := u.Query().Get("tags"); tags != "" {
		opts.Tags = strings.Split(tags, ",")
	}
	if platform := config.Platform(u.Query().Get("platform")); platform.IsValid() {
		opts.Platform = platform
	}
	if fetchInterval, err := time.ParseDuration(u.Query().Get("fetch_interval")); err == nil {
		opts.FetchInterval = fetchInterval
	}
	if format := config.FeedFormat(u.Query().Get("format")); format.IsValid() {
		opts.Format = format
	}
	return opts
}

// Merge merges the given GeneratorOptions into the current options and returns
// a new copy.
func (g GeneratorOptions) Merge(overrides GeneratorOptions) GeneratorOptions {
	if overrides.MaxItems > 0 {
		g.MaxItems = overrides.MaxItems
	}
	if overrides.Sort != "" {
		g.Sort = overrides.Sort
	}
	if overrides.Platform != "" {
		g.Platform = overrides.Platform
	}
	if len(overrides.Tags) > 0 {
		g.Tags = overrides.Tags
	}
	if overrides.FetchInterval > 0 {
		g.FetchInterval = overrides.FetchInterval
	}
	if overrides.Format.IsValid() {
		g.Format = overrides.Format
	}
	return g
}

func (g GeneratorOptions) GetSort() string {
	if g.Sort == "" {
		return sortAliases["recent"]
	}
	if alias, ok := sortAliases[g.Sort]; ok {
		return alias
	}
	return g.Sort
}

var sortAliases = map[string]string{
	"recent":        "-date_live",
	"last_updated":  "-date_updated",
	"trending":      "-downloads_today",
	"highest_rated": "-ratings_weighted_aggregate",
	"popular":       "-downloads_total",
	"subscribers":   "-subscribers_total",
	"alphabetical":  "name",
}
