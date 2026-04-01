package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-prospector/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){stage:=r.URL.Query().Get("stage");list,_:=s.db.List(stage);if list==nil{list=[]store.Lead{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var l store.Lead;json.NewDecoder(r.Body).Decode(&l);if l.Name==""{writeError(w,400,"name required");return};if l.Stage==""{l.Stage="prospect"};s.db.Create(&l);writeJSON(w,201,l)}
func(s *Server)handleUpdate(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var l store.Lead;json.NewDecoder(r.Body).Decode(&l);l.ID=id;s.db.Update(&l);writeJSON(w,200,l)}
func(s *Server)handleClose(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Won bool `json:"won"`};json.NewDecoder(r.Body).Decode(&req);s.db.CloseLead(id,req.Won);writeJSON(w,200,map[string]string{"status":"closed"})}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
