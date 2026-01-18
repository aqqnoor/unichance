package scoring

import "math"

type Profile struct {
  GPA *float64
  GPAScale *float64
  IELTS *float64
  TOEFL *int
  SAT *int
  BudgetYear *float64
}

type Requirements struct {
  MinGPA *float64
  MinIELTS *float64
  MinTOEFL *int
  MinSAT *int
}

type Result struct {
  Score int
  Reasons []string
}

func Compute(p Profile, r Requirements) Result {
  score := 0
  reasons := []string{}

  // GPA (0-40)
  if p.GPA != nil && p.GPAScale != nil && *p.GPAScale > 0 {
    gpaNorm := (*p.GPA) / (*p.GPAScale) // 0..1
    part := int(math.Round(40 * clamp01(gpaNorm)))
    score += part
  } else {
    reasons = append(reasons, "GPA көрсетілмеген, баға дәлдігі төмен")
  }
  if r.MinGPA != nil && p.GPA != nil && p.GPAScale != nil {
    // compare normalized (assume req is on same scale if stored so; MVP)
    if *p.GPA < *r.MinGPA {
      reasons = append(reasons, "GPA талаптан төмен")
    }
  }

  // Language (0-30)
  langPart := 0
  if p.IELTS != nil {
    langPart = int(math.Round(30 * clamp01(*p.IELTS/9.0)))
    if r.MinIELTS != nil && *p.IELTS < *r.MinIELTS {
      reasons = append(reasons, "IELTS талаптан төмен")
    }
  } else if p.TOEFL != nil {
    langPart = int(math.Round(30 * clamp01(float64(*p.TOEFL)/120.0)))
    if r.MinTOEFL != nil && *p.TOEFL < *r.MinTOEFL {
      reasons = append(reasons, "TOEFL талаптан төмен")
    }
  } else {
    reasons = append(reasons, "Тіл сертификаты жоқ (IELTS/TOEFL)")
  }
  score += langPart

  // SAT (0-20)
  if p.SAT != nil {
    satPart := int(math.Round(20 * clamp01(float64(*p.SAT)/1600.0)))
    score += satPart
    if r.MinSAT != nil && *p.SAT < *r.MinSAT {
      reasons = append(reasons, "SAT талаптан төмен")
    }
  } else {
    reasons = append(reasons, "SAT көрсетілмеген (кей универге міндетті)")
  }

  // Extra (0-10) — MVP-де achievements_summary бар болса бонус
  if true { // handler-де қосымша flag береміз
    // placeholder
  }

  if score > 100 { score = 100 }
  if score < 0 { score = 0 }

  return Result{Score: score, Reasons: reasons}
}

func clamp01(x float64) float64 {
  if x < 0 { return 0 }
  if x > 1 { return 1 }
  return x
}
