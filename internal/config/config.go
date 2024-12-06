package config

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DefaultAPIURL        = "https://embed.modhub.io/v1/games/6715/mods"
	DefaultListen        = ":8080"
	DefaultSort          = "recent"
	DefaultMaxItems      = 100
	DefaultFetchInterval = 5 * time.Minute
	DefaultFormat        = FormatAtom
)

type Platform string

const (
	PlatformWindows     Platform = "windows"
	PlatformMac         Platform = "mac"
	PlatformPS5         Platform = "ps5"
	PlatformXBoxSeriesX Platform = "xboxseriesx"
)

func (p Platform) IsValid() bool {
	switch p {
	case PlatformWindows, PlatformMac, PlatformPS5, PlatformXBoxSeriesX:
		return true
	}
	return false
}

type FeedFormat string

const (
	FormatRSS  FeedFormat = "rss"
	FormatAtom FeedFormat = "atom"
	FormatJSON FeedFormat = "json"
)

func (f FeedFormat) IsValid() bool {
	switch f {
	case FormatRSS, FormatAtom, FormatJSON:
		return true
	}
	return false
}

type Configuration struct {
	// Listen is the address to listen on. Defaults to :8080.
	Listen string `mapstructure:"listen"`
	// The API URL to fetch mods from. Defaults to the modhub.io API.
	APIURL string `mapstructure:"api-url"`
	// Tags to filter mods by
	Tags []string `mapstructure:"tags"`
	// Platforms to filter mods by
	Platform Platform `mapstructure:"platform"`
	// MaxFeedItems is the maximum number of feed items to render.
	// Defaults to 100 items.
	MaxFeedItems int `mapstructure:"max-feed-items"`
	// The field to sort the feed by. Defaults to -date_added.
	Sort string `mapstructure:"sort"`
	// FetchInterval is the interval to fetch mods at. Defaults to 5 minutes.
	FetchInterval time.Duration `mapstructure:"fetch-interval"`
	// Format is the format to render the feed in. Valid options are
	// "rss", "atom", and "json". Defaults to "atom".
	Format FeedFormat `mapstructure:"format"`
}

func (c Configuration) Log() {
	log.Println("Configuration:")
	log.Println("    Listen:", c.Listen)
	log.Println("    API URL:", c.APIURL)
	log.Println("    Tags:", strings.Join(c.Tags, ", "))
	log.Println("    Platform:", c.Platform)
	log.Println("    Max Feed Items:", c.MaxFeedItems)
	log.Println("    Sort:", c.Sort)
	log.Println("    Fetch Interval:", c.FetchInterval)
	log.Println("    Format:", c.Format)
}

var viperOnce sync.Once
var viperInstance *viper.Viper

func Load(filename string) (Configuration, error) {
	v := GetViper()
	var c Configuration
	if filename != "" {
		v.SetConfigFile(filename)
		if err := v.ReadInConfig(); err != nil {
			return c, err
		}
	}
	if err := v.Unmarshal(&c); err != nil {
		return c, err
	}
	if !c.Format.IsValid() {
		return c, fmt.Errorf("invalid feed format: %s", c.Format)
	}
	return c, nil
}

func GetViper() *viper.Viper {
	viperOnce.Do(func() {
		v := viper.New()
		v.SetEnvPrefix("BG3MODS")
		v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		v.AutomaticEnv()
		v.SetDefault("listen", DefaultListen)
		v.SetDefault("api-url", DefaultAPIURL)
		v.SetDefault("max-feed-items", DefaultMaxItems)
		v.SetDefault("sort", DefaultSort)
		v.SetDefault("fetch-interval", DefaultFetchInterval)
		v.SetDefault("format", string(DefaultFormat))
		viperInstance = v
	})
	return viperInstance
}

func BindPFlags(flags *pflag.FlagSet) {
	flags.String("listen", DefaultListen, "The address to listen on")
	flags.String("api-url", DefaultAPIURL, "The API URL to fetch mods from")
	flags.StringSlice("tags", nil, "Tags to filter mods by")
	flags.String("platform", "", "Platform to filter mods by (windows, mac, ps5, xboxseriesx)")
	flags.Int("max-feed-items", DefaultMaxItems, "The maximum number of feed items to render")
	flags.String("sort", DefaultSort, "The field to sort the feed by")
	flags.Duration("fetch-interval", DefaultFetchInterval, "The interval to fetch mods at")
	flags.String("format", string(DefaultFormat), "The format to render the feed in (rss, atom, json)")
	if err := GetViper().BindPFlags(flags); err != nil {
		panic(err)
	}
}
