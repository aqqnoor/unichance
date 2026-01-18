package shortlist

import (
  "context"
  "github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct { DB *pgxpool.Pool }

func (r Repo) ensureDefault(ctx context.Context, userID string) (string, error) {
  var id string
  err := r.DB.QueryRow(ctx,
    `SELECT id FROM shortlists WHERE user_id=$1 ORDER BY created_at LIMIT 1`, userID,
  ).Scan(&id)
  if err == nil { return id, nil }

  err = r.DB.QueryRow(ctx,
    `INSERT INTO shortlists(user_id) VALUES ($1) RETURNING id`, userID,
  ).Scan(&id)
  return id, err
}

func (r Repo) Add(ctx context.Context, userID, programID string) error {
  sid, err := r.ensureDefault(ctx, userID)
  if err != nil { return err }
  _, err = r.DB.Exec(ctx,
    `INSERT INTO shortlist_items(shortlist_id, program_id)
     VALUES ($1,$2) ON CONFLICT DO NOTHING`, sid, programID)
  return err
}

func (r Repo) Remove(ctx context.Context, userID, programID string) error {
  sid, err := r.ensureDefault(ctx, userID)
  if err != nil { return err }
  _, err = r.DB.Exec(ctx,
    `DELETE FROM shortlist_items WHERE shortlist_id=$1 AND program_id=$2`, sid, programID)
  return err
}

func (r Repo) List(ctx context.Context, userID string) ([]Item, error) {
  sid, err := r.ensureDefault(ctx, userID)
  if err != nil { return nil, err }

  rows, err := r.DB.Query(ctx, `
    SELECT
      p.id, p.title, u.name, u.country_code,
      p.degree_level::text, p.field,
      p.tuition_amount, p.tuition_currency::text,
      p.has_scholarship, u.qs_rank, u.the_rank
    FROM shortlist_items si
    JOIN programs p ON p.id = si.program_id
    JOIN universities u ON u.id = p.university_id
    WHERE si.shortlist_id=$1
    ORDER BY u.qs_rank ASC NULLS LAST
  `, sid)
  if err != nil { return nil, err }
  defer rows.Close()

  var out []Item
  for rows.Next() {
    var it Item
    if err := rows.Scan(
      &it.ProgramID, &it.Title, &it.University, &it.Country,
      &it.Degree, &it.Field,
      &it.TuitionAmount, &it.TuitionCurrency,
      &it.HasScholarship, &it.QSRank, &it.THERank,
    ); err != nil { return nil, err }
    out = append(out, it)
  }
  return out, rows.Err()
}
