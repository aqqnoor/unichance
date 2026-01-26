package universities

import "time"

type University struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	CountryCode   string     `json:"country_code"`
	City          *string    `json:"city,omitempty"`
	Website       *string    `json:"website,omitempty"`
	QSRank        *int       `json:"qs_rank,omitempty"`
	THERank       *int       `json:"the_rank,omitempty"`
	DataUpdatedAt *time.Time `json:"data_updated_at,omitempty"`

	Links    []UniversityLink `json:"links"`
	Programs []ProgramLite    `json:"programs"`
}

type UniversityLink struct {
	ID             string     `json:"id"`
	LinkType       string     `json:"link_type"`
	URL            string     `json:"url"`
	Title          *string    `json:"title,omitempty"`
	IsOfficial     bool       `json:"is_official"`
	Priority       int        `json:"priority"`
	SourceCode     *string    `json:"source_code,omitempty"`
	LastVerifiedAt *time.Time `json:"last_verified_at,omitempty"`
}

type ProgramLite struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	DegreeLevel     string   `json:"degree_level"`
	Field           string   `json:"field"`
	Language        string   `json:"language"`
	TuitionAmount   *float64 `json:"tuition_amount,omitempty"`
	TuitionCurrency *string  `json:"tuition_currency,omitempty"`
	HasScholarship  bool     `json:"has_scholarship"`
}
