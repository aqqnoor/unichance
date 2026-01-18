package shortlist

type Item struct {
  ProgramID string `json:"program_id"`
  Title string `json:"title"`
  University string `json:"university"`
  Country string `json:"country"`
  Degree string `json:"degree"`
  Field string `json:"field"`
  TuitionAmount *float64 `json:"tuition_amount"`
  TuitionCurrency *string `json:"tuition_currency"`
  HasScholarship bool `json:"has_scholarship"`
  QSRank *int `json:"qs_rank"`
  THERank *int `json:"the_rank"`
}
