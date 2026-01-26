package programs

type ProgramCard struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	DegreeLevel string `json:"degree_level"`
	Field       string `json:"field"`
	Language    string `json:"language"`

	TuitionAmount   *float64 `json:"tuition_amount"`
	TuitionCurrency *string  `json:"tuition_currency"`

	HasScholarship        bool    `json:"has_scholarship"`
	ScholarshipType       *string `json:"scholarship_type"`
	ScholarshipPercentMin *int    `json:"scholarship_percent_min"`
	ScholarshipPercentMax *int    `json:"scholarship_percent_max"`

	UniversityName string  `json:"university_name"`
	CountryCode    string  `json:"country_code"`
	City           *string `json:"city"`
	QSRank         *int    `json:"qs_rank"`
	THERank        *int    `json:"the_rank"`

	UniversityID string `json:"university_id"`
}
