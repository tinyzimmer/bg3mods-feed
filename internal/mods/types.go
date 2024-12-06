package mods

import (
	"time"

	"github.com/tinyzimmer/bg3mods-feed/internal/config"
)

type GetModsResponse struct {
	Data         []Mod `json:"data"`
	ResultCount  int   `json:"result_count"`
	ResultLimit  int   `json:"result_limit"`
	ResultOffset int   `json:"result_offset"`
	ResultTotal  int   `json:"result_total"`
}

type Mod struct {
	ID                   int           `json:"id"`
	GameID               int           `json:"game_id"`
	Status               int           `json:"status"`
	Visible              int           `json:"visible"`
	SubmittedBy          User          `json:"submitted_by"`
	DateAddedEpoch       uint64        `json:"date_added"`
	DateUpdatedEpoch     uint64        `json:"date_updated"`
	DateLiveEpoch        uint64        `json:"date_live"`
	MaturityOption       int           `json:"maturity_option"`
	CommunityOptions     int           `json:"community_options"`
	CreditOptions        int           `json:"credit_options"`
	MonetizationOptions  int           `json:"monetization_options"`
	Stock                float64       `json:"stock"`
	Price                float64       `json:"price"`
	Tax                  float64       `json:"tax"`
	Logo                 Image         `json:"logo"`
	HomepageURL          string        `json:"homepage_url"`
	Name                 string        `json:"name"`
	NameID               string        `json:"name_id"`
	Summary              string        `json:"summary"`
	Description          string        `json:"description"`
	DescriptionPlaintext string        `json:"description_plaintext"`
	MetadataBlob         string        `json:"metadata_blob"`
	ProfileURL           string        `json:"profile_url"`
	Dependencies         bool          `json:"dependencies"`
	Media                Media         `json:"media"`
	Modfile              Modfile       `json:"modfile"`
	Platforms            []Platform    `json:"platforms"`
	Tags                 []Tag         `json:"tags"`
	Stats                Stats         `json:"stats"`
	PlatformStats        PlatformStats `json:"platform_stats"`
	GameName             string        `json:"game_name"`
}

func (m Mod) String() string {
	return m.Name
}

func (m Mod) DateAdded() time.Time {
	return time.Unix(int64(m.DateAddedEpoch), 0)
}

func (m Mod) DateUpdated() time.Time {
	return time.Unix(int64(m.DateUpdatedEpoch), 0)
}

func (m Mod) DateLive() time.Time {
	return time.Unix(int64(m.DateLiveEpoch), 0)
}

func (m Mod) SupportsPlatform(platform config.Platform) bool {
	for _, p := range m.Modfile.Platforms {
		if p.Platform == string(platform) {
			return p.Status == 1
		}
	}
	return false
}

type User struct {
	ID                int    `json:"id"`
	NameID            string `json:"name_id"`
	Username          string `json:"username"`
	DisplayNamePortal string `json:"display_name_portal"`
	DateOnline        uint64 `json:"date_online"`
	DateJoined        uint64 `json:"date_joined"`
	Avatar            Avatar `json:"avatar"`
	Timezone          string `json:"timezone"`
	Language          string `json:"language"`
	ProfileURL        string `json:"profile_url"`
}

type Avatar struct {
	Filename     string `json:"filename"`
	Original     string `json:"original"`
	Thumb50x50   string `json:"thumb_50x50"`
	Thumb100x100 string `json:"thumb_100x100"`
}

type Image struct {
	Filename      string `json:"filename"`
	Original      string `json:"original"`
	Thumb320x180  string `json:"thumb_320x180"`
	Thumb640x360  string `json:"thumb_640x360"`
	Thumb1280x720 string `json:"thumb_1280x720"`
}

type Media struct {
	Images []Image `json:"images"`
}

type Modfile struct {
	ID                   int        `json:"id"`
	ModID                int        `json:"mod_id"`
	DateAdded            uint64     `json:"date_added"`
	DateUpdated          uint64     `json:"date_updated"`
	DateScanned          uint64     `json:"date_scanned"`
	VirusStatus          int        `json:"virus_status"`
	VirusPositive        int        `json:"virus_positive"`
	VirusTotalHash       string     `json:"virustotal_hash"`
	Filesize             uint64     `json:"filesize"`
	FilesizeUncompressed uint64     `json:"filesize_uncompressed"`
	Filehash             Hash       `json:"filehash"`
	Filename             string     `json:"filename"`
	Version              string     `json:"version"`
	Changelog            string     `json:"changelog"`
	MetadataBlob         string     `json:"metadata_blob"`
	Download             Download   `json:"download"`
	Platforms            []Platform `json:"platforms"`
}

type Hash struct {
	MD5 string `json:"md5"`
}

type Download struct {
	BinaryURL   string `json:"binary_url"`
	DateExpires uint64 `json:"date_expires"`
}

type Platform struct {
	Platform    string `json:"platform"`
	Status      int    `json:"status"`
	ModfileLive int    `json:"modfile_live"`
}

type Tag struct {
	Name          string `json:"name"`
	NameLocalized string `json:"name_localized"`
	DateAdded     uint64 `json:"date_added"`
}

type Stats struct {
	ModID                     int     `json:"mod_id"`
	PopularityRankPosition    int     `json:"popularity_rank_position"`
	PopularityRankTotalMods   int     `json:"popularity_rank_total_mods"`
	DownloadsToday            int     `json:"downloads_today"`
	DownloadsTotal            int     `json:"downloads_total"`
	SubscribersTotal          int     `json:"subscribers_total"`
	RatingsTotal              int     `json:"ratings_total"`
	RatingsPositive           int     `json:"ratings_positive"`
	RatingsNegative           int     `json:"ratings_negative"`
	RatingsPercentagePositive float64 `json:"ratings_percentage_positive"`
	RatingsWeightedAggregate  float64 `json:"ratings_weighted_aggregate"`
	RatingsDisplayText        string  `json:"ratings_display_text"`
	DateExpires               uint64  `json:"date_expires"`
}

type PlatformStats struct {
	LiveCount    int `json:"live_count"`
	PendingCount int `json:"pending_count"`
}
