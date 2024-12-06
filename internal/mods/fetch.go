package mods

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FetchOptions are the options for fetching mods from the API.
type FetchOptions struct {
	Limit  int
	Offset int
	Sort   string
	Tags   []string
}

// Fetcher is the interface for fetching mods from the API.
type Fetcher interface {
	// Fetch fetches mods from the API based on the given options.
	Fetch(context.Context, FetchOptions) (*GetModsResponse, error)
}

type fetcher struct {
	apiURL string
}

// NewFetcher creates a new Fetcher using the given API URL.
func NewFetcher(apiURL string) Fetcher {
	return &fetcher{apiURL}
}

func (f *fetcher) Fetch(ctx context.Context, opts FetchOptions) (*GetModsResponse, error) {
	url, err := f.optsToURL(opts)
	if err != nil {
		return nil, err
	}
	log.Println("Fetching mods from", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Modio-Origin", "web")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, body)
	}
	var modResp GetModsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modResp); err != nil {
		return nil, err
	}
	return &modResp, nil
}

func (f *fetcher) optsToURL(opts FetchOptions) (string, error) {
	u, err := url.Parse(f.apiURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	if opts.Limit > 0 {
		q.Set("_limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		q.Set("_offset", strconv.Itoa(opts.Offset))
	}
	if opts.Sort != "" {
		q.Set("_sort", opts.Sort)
	}
	if len(opts.Tags) > 0 {
		q.Set("tags-in", strings.Join(opts.Tags, ","))
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
