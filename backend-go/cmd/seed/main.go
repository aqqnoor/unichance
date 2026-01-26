package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func parseTime(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func parseBool(s string) *bool {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return nil
	}
	if s == "true" || s == "1" {
		b := true
		return &b
	}
	if s == "false" || s == "0" {
		b := false
		return &b
	}
	return nil
}

func parseIntPtr(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func main() {
	_ = godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL required")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ctx := context.Background()

	// universities
	uniMap := map[string]string{} // name -> id

	uf, err := os.Open("seed/universities.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer uf.Close()

	ur := csv.NewReader(uf)
	ur.FieldsPerRecord = -1
	uRows, err := ur.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(uRows) < 2 {
		log.Fatal("universities.csv empty")
	}

	// header: name,country_code,city,website,qs_rank,the_rank,data_updated_at
	for i := 1; i < len(uRows); i++ {
		r := uRows[i]
		if len(r) < 7 {
			continue
		}
		name := strings.TrimSpace(r[0])
		if name == "" {
			continue
		}

		country := strings.TrimSpace(r[1])
		city := strings.TrimSpace(r[2])
		website := strings.TrimSpace(r[3])

		var qs *int
		if strings.TrimSpace(r[4]) != "" {
			v, _ := strconv.Atoi(r[4])
			qs = &v
		}
		var the *int
		if strings.TrimSpace(r[5]) != "" {
			v, _ := strconv.Atoi(r[5])
			the = &v
		}
		updated := parseTime(r[6])

		var id string
		err := pool.QueryRow(ctx, `
      INSERT INTO universities(name,country_code,city,website,qs_rank,the_rank,data_updated_at)
      VALUES ($1,$2,NULLIF($3,''),NULLIF($4,''),$5,$6,$7)
      RETURNING id
    `, name, country, city, website, qs, the, updated).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		uniMap[name] = id
	}

	// programs
	pf, err := os.Open("seed/programs.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer pf.Close()

	pr := csv.NewReader(pf)
	pr.FieldsPerRecord = -1
	pRows, err := pr.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(pRows) < 2 {
		log.Fatal("programs.csv empty")
	}

	// sources
	srcMap := map[string]string{} // code -> id

	sf, err := os.Open("seed/sources.csv")
	if err != nil {
		log.Printf("sources.csv not found, skip: %v", err)
	} else {
		defer sf.Close()

		sr := csv.NewReader(sf)
		sr.FieldsPerRecord = -1
		sRows, err := sr.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		if len(sRows) >= 2 {
			// header: code,name,kind,base_url,docs_url,license,reliability,is_active,refresh_interval_hours,last_fetched_at
			for i := 1; i < len(sRows); i++ {
				r := sRows[i]
				if len(r) < 10 {
					continue
				}

				code := strings.TrimSpace(r[0])
				name := strings.TrimSpace(r[1])
				kind := strings.TrimSpace(r[2])
				baseURL := strings.TrimSpace(r[3])
				docsURL := strings.TrimSpace(r[4])
				license := strings.TrimSpace(r[5])

				reliability := 3
				if strings.TrimSpace(r[6]) != "" {
					v, _ := strconv.Atoi(r[6])
					if v >= 1 && v <= 5 {
						reliability = v
					}
				}

				isActive := true
				if pb := parseBool(r[7]); pb != nil {
					isActive = *pb
				}

				var refreshHours *int
				if v := parseIntPtr(r[8]); v != nil {
					if *v > 0 {
						refreshHours = v
					} else {
						// 0 allowed as "no schedule"
						zero := 0
						refreshHours = &zero
					}
				}

				lastFetched := parseTime(r[9])

				var id string
				err := pool.QueryRow(ctx, `
          INSERT INTO sources(
            code,name,kind,base_url,docs_url,license,reliability,is_active,refresh_interval_hours,last_fetched_at
          ) VALUES (
            $1,$2,$3,NULLIF($4,''),NULLIF($5,''),NULLIF($6,''),$7,$8,$9,$10
          )
          ON CONFLICT (code) DO UPDATE SET
            name=EXCLUDED.name,
            kind=EXCLUDED.kind,
            base_url=EXCLUDED.base_url,
            docs_url=EXCLUDED.docs_url,
            license=EXCLUDED.license,
            reliability=EXCLUDED.reliability,
            is_active=EXCLUDED.is_active,
            refresh_interval_hours=EXCLUDED.refresh_interval_hours,
            last_fetched_at=EXCLUDED.last_fetched_at,
            updated_at=NOW()
          RETURNING id
        `, code, name, kind, baseURL, docsURL, license, reliability, isActive, refreshHours, lastFetched).Scan(&id)
				if err != nil {
					log.Fatal(err)
				}

				srcMap[code] = id
			}
			log.Printf("seed sources done: %d\n", len(srcMap))
		}
	}

	// university_links
	lf, err := os.Open("seed/university_links.csv")
	if err != nil {
		log.Printf("university_links.csv not found, skip: %v", err)
	} else {
		defer lf.Close()

		lr := csv.NewReader(lf)
		lr.FieldsPerRecord = -1
		lRows, err := lr.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		insertedLinks := 0

		if len(lRows) >= 2 {
			// header: university_name,link_type,url,title,is_official,priority,source_code,last_verified_at
			for i := 1; i < len(lRows); i++ {
				r := lRows[i]
				if len(r) < 8 {
					continue
				}

				uniName := strings.TrimSpace(r[0])
				uniID := uniMap[uniName]
				if uniID == "" {
					continue
				}

				linkType := strings.TrimSpace(r[1])
				url := strings.TrimSpace(r[2])
				if linkType == "" || url == "" {
					continue
				}

				title := strings.TrimSpace(r[3])

				isOfficial := true
				if pb := parseBool(r[4]); pb != nil {
					isOfficial = *pb
				}

				priority := 10
				if strings.TrimSpace(r[5]) != "" {
					v, _ := strconv.Atoi(r[5])
					if v >= 1 && v <= 100 {
						priority = v
					}
				}

				sourceCode := strings.TrimSpace(r[6])
				sourceID := ""
				if sourceCode != "" {
					sourceID = srcMap[sourceCode]
				}
				verified := parseTime(r[7])

				// ON CONFLICT uses your constraint: (university_id, link_type, url)
				if sourceID == "" {
					_, err = pool.Exec(ctx, `
            INSERT INTO university_links(
              university_id, source_id, link_type, url, title, is_official, priority, last_verified_at
            ) VALUES (
              $1, NULL, $2, $3, NULLIF($4,''), $5, $6, $7
            )
            ON CONFLICT (university_id, link_type, url) DO UPDATE SET
              title=EXCLUDED.title,
              is_official=EXCLUDED.is_official,
              priority=EXCLUDED.priority,
              last_verified_at=EXCLUDED.last_verified_at,
              updated_at=NOW()
          `, uniID, linkType, url, title, isOfficial, priority, verified)
				} else {
					_, err = pool.Exec(ctx, `
            INSERT INTO university_links(
              university_id, source_id, link_type, url, title, is_official, priority, last_verified_at
            ) VALUES (
              $1, $2, $3, $4, NULLIF($5,''), $6, $7, $8
            )
            ON CONFLICT (university_id, link_type, url) DO UPDATE SET
              source_id=EXCLUDED.source_id,
              title=EXCLUDED.title,
              is_official=EXCLUDED.is_official,
              priority=EXCLUDED.priority,
              last_verified_at=EXCLUDED.last_verified_at,
              updated_at=NOW()
          `, uniID, sourceID, linkType, url, title, isOfficial, priority, verified)
				}

				if err != nil {
					log.Fatal(err)
				}
				insertedLinks++
			}
		}

		log.Printf("seed university_links done: %d\n", insertedLinks)
    log.Printf("seed done: universities=%d programs=%d\n", len(uniMap), len(pRows)-1)

	}

	// header:
	// university_name,title,degree_level,field,language,tuition_amount,tuition_currency,has_scholarship,scholarship_type,scholarship_percent_min,scholarship_percent_max,description,data_updated_at
	for i := 1; i < len(pRows); i++ {
		r := pRows[i]
		if len(r) < 13 {
			continue
		}

		uniName := strings.TrimSpace(r[0])
		uniID := uniMap[uniName]
		if uniID == "" {
			continue
		}

		title := strings.TrimSpace(r[1])
		level := strings.TrimSpace(r[2]) // bachelor/master
		field := strings.TrimSpace(r[3])
		lang := strings.TrimSpace(r[4])

		var tuition *float64
		if strings.TrimSpace(r[5]) != "" {
			v, _ := strconv.ParseFloat(r[5], 64)
			tuition = &v
		}
		currency := strings.TrimSpace(r[6])

		hasSch := strings.TrimSpace(r[7]) == "true"
		schType := strings.TrimSpace(r[8])

		var schMin *int
		if strings.TrimSpace(r[9]) != "" {
			v, _ := strconv.Atoi(r[9])
			schMin = &v
		}
		var schMax *int
		if strings.TrimSpace(r[10]) != "" {
			v, _ := strconv.Atoi(r[10])
			schMax = &v
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
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("seed done: universities=%d programs=%d\n", len(uniMap), len(pRows)-1)
}
