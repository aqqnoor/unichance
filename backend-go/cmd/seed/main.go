package main

import (
  "context"
  "encoding/csv"
  "log"
  "os"
  "strconv"
  "strings"
  "time"

  "github.com/joho/godotenv"
  "github.com/jackc/pgx/v5/pgxpool"
)

func parseTime(s string) *time.Time {
  s = strings.TrimSpace(s)
  if s == "" { return nil }
  t, err := time.Parse(time.RFC3339, s)
  if err != nil { return nil }
  return &t
}

func main() {
  _ = godotenv.Load(".env")
  dbURL := os.Getenv("DATABASE_URL")
  if dbURL == "" { log.Fatal("DATABASE_URL required") }

  pool, err := pgxpool.New(context.Background(), dbURL)
  if err != nil { log.Fatal(err) }
  defer pool.Close()

  ctx := context.Background()

  // universities
  uniMap := map[string]string{} // name -> id

  uf, err := os.Open("seed/universities.csv")
  if err != nil { log.Fatal(err) }
  defer uf.Close()

  ur := csv.NewReader(uf)
  ur.FieldsPerRecord = -1
  uRows, err := ur.ReadAll()
  if err != nil { log.Fatal(err) }
  if len(uRows) < 2 { log.Fatal("universities.csv empty") }

  // header: name,country_code,city,website,qs_rank,the_rank,data_updated_at
  for i := 1; i < len(uRows); i++ {
    r := uRows[i]
    if len(r) < 7 { continue }
    name := strings.TrimSpace(r[0])
    if name == "" { continue }

    country := strings.TrimSpace(r[1])
    city := strings.TrimSpace(r[2])
    website := strings.TrimSpace(r[3])

    var qs *int
    if strings.TrimSpace(r[4]) != "" {
      v, _ := strconv.Atoi(r[4]); qs = &v
    }
    var the *int
    if strings.TrimSpace(r[5]) != "" {
      v, _ := strconv.Atoi(r[5]); the = &v
    }
    updated := parseTime(r[6])

    var id string
    err := pool.QueryRow(ctx, `
      INSERT INTO universities(name,country_code,city,website,qs_rank,the_rank,data_updated_at)
      VALUES ($1,$2,NULLIF($3,''),NULLIF($4,''),$5,$6,$7)
      RETURNING id
    `, name, country, city, website, qs, the, updated).Scan(&id)
    if err != nil { log.Fatal(err) }
    uniMap[name] = id
  }

  // programs
  pf, err := os.Open("seed/programs.csv")
  if err != nil { log.Fatal(err) }
  defer pf.Close()

  pr := csv.NewReader(pf)
  pr.FieldsPerRecord = -1
  pRows, err := pr.ReadAll()
  if err != nil { log.Fatal(err) }
  if len(pRows) < 2 { log.Fatal("programs.csv empty") }

  // header:
  // university_name,title,degree_level,field,language,tuition_amount,tuition_currency,has_scholarship,scholarship_type,scholarship_percent_min,scholarship_percent_max,description,data_updated_at
  for i := 1; i < len(pRows); i++ {
    r := pRows[i]
    if len(r) < 13 { continue }

    uniName := strings.TrimSpace(r[0])
    uniID := uniMap[uniName]
    if uniID == "" { continue }

    title := strings.TrimSpace(r[1])
    level := strings.TrimSpace(r[2]) // bachelor/master
    field := strings.TrimSpace(r[3])
    lang := strings.TrimSpace(r[4])

    var tuition *float64
    if strings.TrimSpace(r[5]) != "" {
      v, _ := strconv.ParseFloat(r[5], 64); tuition = &v
    }
    currency := strings.TrimSpace(r[6])

    hasSch := strings.TrimSpace(r[7]) == "true"
    schType := strings.TrimSpace(r[8])

    var schMin *int
    if strings.TrimSpace(r[9]) != "" {
      v, _ := strconv.Atoi(r[9]); schMin = &v
    }
    var schMax *int
    if strings.TrimSpace(r[10]) != "" {
      v, _ := strconv.Atoi(r[10]); schMax = &v
    }

    desc := strings.TrimSpace(r[11])
    updated := parseTime(r[12])

    _, err := pool.Exec(ctx, `
      INSERT INTO programs(
        university_id,title,degree_level,field,language,
        tuition_amount,tuition_currency,has_scholarship,scholarship_type,scholarship_percent_min,scholarship_percent_max,
        description,data_updated_at
      ) VALUES (
        $1,$2,$3,$4,$5,
        $6,NULLIF($7,'')::tuition_currency,$8,NULLIF($9,''),$10,$11,
        NULLIF($12,''),$13
      )
    `, uniID, title, level, field, lang, tuition, currency, hasSch, schType, schMin, schMax, desc, updated)
    if err != nil { log.Fatal(err) }
  }

  log.Printf("seed done: universities=%d programs=%d\n", len(uniMap), len(pRows)-1)
}
