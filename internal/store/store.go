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

// Deal represents a single sales pipeline opportunity.
// Value is stored as integer dollars (not cents) for backward compatibility
// with the original schema. Probability is 0-100.
type Deal struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Company      string `json:"company"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	Value        int    `json:"value"`
	Stage        string `json:"stage"`
	Probability  int    `json:"probability"`
	CloseDate    string `json:"close_date"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "prospector.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS deals(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		company TEXT DEFAULT '',
		contact_name TEXT DEFAULT '',
		contact_email TEXT DEFAULT '',
		value INTEGER DEFAULT 0,
		stage TEXT DEFAULT 'lead',
		probability INTEGER DEFAULT 0,
		close_date TEXT DEFAULT '',
		notes TEXT DEFAULT '',
		created_at TEXT DEFAULT(datetime('now'))
	)`)
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

func (d *DB) Create(e *Deal) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Stage == "" {
		e.Stage = "lead"
	}
	_, err := d.db.Exec(
		`INSERT INTO deals(id, name, company, contact_name, contact_email, value, stage, probability, close_date, notes, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Company, e.ContactName, e.ContactEmail, e.Value, e.Stage, e.Probability, e.CloseDate, e.Notes, e.CreatedAt,
	)
	return err
}

func (d *DB) Get(id string) *Deal {
	var e Deal
	err := d.db.QueryRow(
		`SELECT id, name, company, contact_name, contact_email, value, stage, probability, close_date, notes, created_at
		 FROM deals WHERE id=?`,
		id,
	).Scan(&e.ID, &e.Name, &e.Company, &e.ContactName, &e.ContactEmail, &e.Value, &e.Stage, &e.Probability, &e.CloseDate, &e.Notes, &e.CreatedAt)
	if err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Deal {
	rows, _ := d.db.Query(
		`SELECT id, name, company, contact_name, contact_email, value, stage, probability, close_date, notes, created_at
		 FROM deals
		 ORDER BY value DESC, created_at DESC`,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Deal
	for rows.Next() {
		var e Deal
		rows.Scan(&e.ID, &e.Name, &e.Company, &e.ContactName, &e.ContactEmail, &e.Value, &e.Stage, &e.Probability, &e.CloseDate, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

func (d *DB) Update(e *Deal) error {
	_, err := d.db.Exec(
		`UPDATE deals SET name=?, company=?, contact_name=?, contact_email=?, value=?, stage=?, probability=?, close_date=?, notes=?
		 WHERE id=?`,
		e.Name, e.Company, e.ContactName, e.ContactEmail, e.Value, e.Stage, e.Probability, e.CloseDate, e.Notes, e.ID,
	)
	return err
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM deals WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM deals`).Scan(&n)
	return n
}

func (d *DB) Search(q string, filters map[string]string) []Deal {
	where := "1=1"
	args := []any{}
	if q != "" {
		where += " AND (name LIKE ? OR company LIKE ? OR contact_name LIKE ? OR contact_email LIKE ?)"
		args = append(args, "%"+q+"%", "%"+q+"%", "%"+q+"%", "%"+q+"%")
	}
	if v, ok := filters["stage"]; ok && v != "" {
		where += " AND stage=?"
		args = append(args, v)
	}
	rows, _ := d.db.Query(
		`SELECT id, name, company, contact_name, contact_email, value, stage, probability, close_date, notes, created_at
		 FROM deals WHERE `+where+`
		 ORDER BY value DESC, created_at DESC`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Deal
	for rows.Next() {
		var e Deal
		rows.Scan(&e.ID, &e.Name, &e.Company, &e.ContactName, &e.ContactEmail, &e.Value, &e.Stage, &e.Probability, &e.CloseDate, &e.Notes, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// Stats returns deal counts and pipeline values, broken down by stage.
// pipeline_value excludes won and lost deals (the active pipeline).
// won_value sums all won deals (closed-won historical).
// weighted_value applies probability to active pipeline deals.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":          d.Count(),
		"pipeline_value": 0,
		"won_value":      0,
		"weighted_value": 0,
		"by_stage":       map[string]int{},
	}

	var pipeline int
	d.db.QueryRow(`SELECT COALESCE(SUM(value), 0) FROM deals WHERE stage NOT IN ('won', 'lost')`).Scan(&pipeline)
	m["pipeline_value"] = pipeline

	var won int
	d.db.QueryRow(`SELECT COALESCE(SUM(value), 0) FROM deals WHERE stage='won'`).Scan(&won)
	m["won_value"] = won

	// Weighted value: sum of (value * probability/100) for active deals.
	// Useful for forecast accuracy.
	var weighted float64
	d.db.QueryRow(`SELECT COALESCE(SUM(value * probability) / 100.0, 0) FROM deals WHERE stage NOT IN ('won', 'lost')`).Scan(&weighted)
	m["weighted_value"] = int(weighted)

	if rows, _ := d.db.Query(`SELECT stage, COUNT(*) FROM deals GROUP BY stage`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var s string
			var c int
			rows.Scan(&s, &c)
			by[s] = c
		}
		m["by_stage"] = by
	}

	return m
}

// ─── Extras: generic key-value storage for personalization custom fields ───

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
