-- shortlists (1 user = many shortlists, MVP-де 1 default қолданамыз)
CREATE TABLE IF NOT EXISTS shortlists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL DEFAULT 'My shortlist',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_shortlists_user ON shortlists(user_id);

-- shortlist items
CREATE TABLE IF NOT EXISTS shortlist_items (
  shortlist_id UUID NOT NULL REFERENCES shortlists(id) ON DELETE CASCADE,
  program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (shortlist_id, program_id)
);
