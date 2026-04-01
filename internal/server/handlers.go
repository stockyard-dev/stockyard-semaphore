package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-semaphore/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Member{}};writeJSON(w,200,list)}
func(s *Server)handleSet(w http.ResponseWriter,r *http.Request){var m store.Member;json.NewDecoder(r.Body).Decode(&m);if m.Name==""{writeError(w,400,"name required");return};if m.Status==""{m.Status="available"};if m.Emoji==""{m.Emoji="🟢"};s.db.Upsert(&m);writeJSON(w,200,m)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
