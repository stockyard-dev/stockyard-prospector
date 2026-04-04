package server

import (
	"encoding/json"
	"net/http"

	"github.com/stockyard-dev/stockyard-prospector/internal/store"
)

type Server struct{ db *store.DB; mux *http.ServeMux; limits Limits }

func New(db *store.DB, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits}
	s.mux.HandleFunc("GET /api/deals", s.list)
	s.mux.HandleFunc("POST /api/deals", s.create)
	s.mux.HandleFunc("GET /api/deals/{id}", s.get)
	s.mux.HandleFunc("PUT /api/deals/{id}", s.update)
	s.mux.HandleFunc("PATCH /api/deals/{id}/stage", s.setStage)
	s.mux.HandleFunc("DELETE /api/deals/{id}", s.del)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"tier": s.limits.Tier, "upgrade_url": "https://stockyard.dev/prospector/"}) })
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
	return s
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func wj(w http.ResponseWriter, c int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(c); json.NewEncoder(w).Encode(v) }
func we(w http.ResponseWriter, c int, m string) { wj(w, c, map[string]string{"error": m}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) { if r.URL.Path != "/" { http.NotFound(w, r); return }; http.Redirect(w, r, "/ui", 302) }
func od(d []store.Deal) []store.Deal { if d == nil { return []store.Deal{} }; return d }

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	stage := r.URL.Query().Get("stage"); q := r.URL.Query().Get("q")
	filters := map[string]string{}; if stage != "" { filters["stage"] = stage }
	if q != "" || len(filters) > 0 { wj(w, 200, map[string]any{"deals": od(s.db.Search(q, filters))}); return }
	wj(w, 200, map[string]any{"deals": od(s.db.List())})
}
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 && s.db.Count() >= s.limits.MaxItems { we(w, 402, "Free tier limit reached"); return }
	var d store.Deal; json.NewDecoder(r.Body).Decode(&d); if d.Name == "" { we(w, 400, "name required"); return }
	if d.Stage == "" { d.Stage = "lead" }; s.db.Create(&d); wj(w, 201, s.db.Get(d.ID))
}
func (s *Server) get(w http.ResponseWriter, r *http.Request) { d := s.db.Get(r.PathValue("id")); if d == nil { we(w, 404, "not found"); return }; wj(w, 200, d) }
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	existing := s.db.Get(r.PathValue("id")); if existing == nil { we(w, 404, "not found"); return }
	var d store.Deal; json.NewDecoder(r.Body).Decode(&d); d.ID = existing.ID; d.CreatedAt = existing.CreatedAt
	if d.Name == "" { d.Name = existing.Name }; if d.Stage == "" { d.Stage = existing.Stage }
	s.db.Update(&d); wj(w, 200, s.db.Get(d.ID))
}
func (s *Server) setStage(w http.ResponseWriter, r *http.Request) {
	d := s.db.Get(r.PathValue("id")); if d == nil { we(w, 404, "not found"); return }
	var body struct { Stage string `json:"stage"` }; json.NewDecoder(r.Body).Decode(&body)
	d.Stage = body.Stage; s.db.Update(d); wj(w, 200, s.db.Get(d.ID))
}
func (s *Server) del(w http.ResponseWriter, r *http.Request) { s.db.Delete(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }
func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	deals := s.db.List(); if deals == nil { deals = []store.Deal{} }
	pipeline := map[string]int{}; totalValue := 0; weighted := 0
	for _, d := range deals { pipeline[d.Stage]++; totalValue += d.Value; weighted += d.Value * d.Probability / 100 }
	wj(w, 200, map[string]any{"total": len(deals), "pipeline": pipeline, "total_value": totalValue, "weighted_value": weighted})
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"service": "prospector", "status": "ok", "deals": s.db.Count()}) }
