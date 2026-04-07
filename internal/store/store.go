package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// TeamMember represents a person whose current availability is tracked.
// Availability is one of: available, busy, away, in_meeting, off.
// UntilTime, when set, marks an automatic return to "available" once
// the time has passed (caller is responsible for invoking ClearExpired).
type TeamMember struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Availability  string `json:"availability"`
	StatusMessage string `json:"status_message"`
	Timezone      string `json:"timezone"`
	UntilTime     string `json:"until_time"`
	CreatedAt     string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "semaphore.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS team_members(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT DEFAULT '',
		role TEXT DEFAULT '',
		availability TEXT DEFAULT 'available',
		status_message TEXT DEFAULT '',
		timezone TEXT DEFAULT '',
		until_time TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_members_availability ON team_members(availability)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_members_role ON team_members(role)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS extras(
		resource TEXT NOT NULL,
		record_id TEXT NOT NULL,
		data TEXT NOT NULL DEFAULT '{}',
		PRIMARY KEY(resource, record_id)
	)`)
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

func (d *DB) Create(e *TeamMember) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Availability == "" {
		e.Availability = "available"
	}
	_, err := d.db.Exec(
		`INSERT INTO team_members(id, name, email, role, availability, status_message, timezone, until_time, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Email, e.Role, e.Availability, e.StatusMessage, e.Timezone, e.UntilTime, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *TeamMember {
	var e TeamMember
	err := d.db.QueryRow(
		`SELECT id, name, email, role, availability, status_message, timezone, until_time, created_at
		 FROM team_members WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Availability, &e.StatusMessage, &e.Timezone, &e.UntilTime, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []TeamMember {
	rows, _ := d.db.Query(
		`SELECT id, name, email, role, availability, status_message, timezone, until_time, created_at
		 FROM team_members ORDER BY name ASC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []TeamMember
	for rows.Next() {
		var e TeamMember
		rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Availability, &e.StatusMessage, &e.Timezone, &e.UntilTime, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *TeamMember) error {
	_, err := d.db.Exec(
		`UPDATE team_members SET name=?, email=?, role=?, availability=?, status_message=?, timezone=?, until_time=?
		 WHERE id=?`,
		e.Name, e.Email, e.Role, e.Availability, e.StatusMessage, e.Timezone, e.UntilTime, e.ID,
	)
	return err
}

// SetAvailability is the dedicated path for the most common operation:
// quickly flipping someone's availability (with an optional message and
// auto-return time) without touching their other fields.
func (d *DB) SetAvailability(id, availability, message, untilTime string) error {
	_, err := d.db.Exec(
		`UPDATE team_members SET availability=?, status_message=?, until_time=? WHERE id=?`,
		availability, message, untilTime, id,
	)
	return err
}

// ClearExpired finds members whose until_time is in the past and resets
// them to "available". Called periodically by the server.
func (d *DB) ClearExpired() int {
	res, err := d.db.Exec(
		`UPDATE team_members
		 SET availability='available', status_message='', until_time=''
		 WHERE until_time != '' AND until_time < ?`,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return 0
	}
	n, _ := res.RowsAffected()
	return int(n)
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM team_members WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM team_members`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []TeamMember {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR email LIKE ? OR role LIKE ?)"
		s := "%" + q + "%"
		args = append(args, s, s, s)
	}
	if v, ok := filters["availability"]; ok && v != "" {
		where += " AND availability=?"
		args = append(args, v)
	}
	if v, ok := filters["role"]; ok && v != "" {
		where += " AND role=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, email, role, availability, status_message, timezone, until_time, created_at
		 FROM team_members WHERE `+where+`
		 ORDER BY name ASC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []TeamMember
	for rows.Next() {
		var e TeamMember
		rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Availability, &e.StatusMessage, &e.Timezone, &e.UntilTime, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns total members and a count by availability state plus
// the number currently marked available.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":           d.Count(),
		"available":       0,
		"by_availability": map[string]int{},
		"by_role":         map[string]int{},
	}

	var available int
	d.db.QueryRow(`SELECT COUNT(*) FROM team_members WHERE availability='available'`).Scan(&available)
	m["available"] = available

	if rows, _ := d.db.Query(`SELECT availability, COUNT(*) FROM team_members GROUP BY availability`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_availability"] = by
	}

	if rows, _ := d.db.Query(`SELECT role, COUNT(*) FROM team_members WHERE role != '' GROUP BY role`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_role"] = by
	}

	return m
}

// ─── Extras ───────────────────────────────────────────────────────

func (d *DB) GetExtras(resource, recordID string) string {
	var data string
	err := d.db.QueryRow(
		`SELECT data FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	).Scan(&data)
	if err != nil || data == "" {
		return "{}"
	}
	return data
}

func (d *DB) SetExtras(resource, recordID, data string) error {
	if data == "" {
		data = "{}"
	}
	_, err := d.db.Exec(
		`INSERT INTO extras(resource, record_id, data) VALUES(?, ?, ?)
		 ON CONFLICT(resource, record_id) DO UPDATE SET data=excluded.data`,
		resource, recordID, data,
	)
	return err
}

func (d *DB) DeleteExtras(resource, recordID string) error {
	_, err := d.db.Exec(
		`DELETE FROM extras WHERE resource=? AND record_id=?`,
		resource, recordID,
	)
	return err
}

func (d *DB) AllExtras(resource string) map[string]string {
	out := make(map[string]string)
	rows, _ := d.db.Query(
		`SELECT record_id, data FROM extras WHERE resource=?`,
		resource,
	)
	if rows == nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id, data string
		rows.Scan(&id, &data)
		out[id] = data
	}
	return out
}
