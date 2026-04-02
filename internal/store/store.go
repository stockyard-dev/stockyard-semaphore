package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type TeamMember struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	Availability string `json:"availability"`
	StatusMessage string `json:"status_message"`
	Timezone string `json:"timezone"`
	UntilTime string `json:"until_time"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"semaphore.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS team_members(id TEXT PRIMARY KEY,name TEXT NOT NULL,email TEXT DEFAULT '',role TEXT DEFAULT '',availability TEXT DEFAULT 'available',status_message TEXT DEFAULT '',timezone TEXT DEFAULT '',until_time TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *TeamMember)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO team_members(id,name,email,role,availability,status_message,timezone,until_time,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Email,e.Role,e.Availability,e.StatusMessage,e.Timezone,e.UntilTime,e.CreatedAt);return err}
func(d *DB)Get(id string)*TeamMember{var e TeamMember;if d.db.QueryRow(`SELECT id,name,email,role,availability,status_message,timezone,until_time,created_at FROM team_members WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Email,&e.Role,&e.Availability,&e.StatusMessage,&e.Timezone,&e.UntilTime,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]TeamMember{rows,_:=d.db.Query(`SELECT id,name,email,role,availability,status_message,timezone,until_time,created_at FROM team_members ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []TeamMember;for rows.Next(){var e TeamMember;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Role,&e.Availability,&e.StatusMessage,&e.Timezone,&e.UntilTime,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *TeamMember)error{_,err:=d.db.Exec(`UPDATE team_members SET name=?,email=?,role=?,availability=?,status_message=?,timezone=?,until_time=? WHERE id=?`,e.Name,e.Email,e.Role,e.Availability,e.StatusMessage,e.Timezone,e.UntilTime,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM team_members WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM team_members`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]TeamMember{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR email LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["availability"];ok&&v!=""{where+=" AND availability=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,email,role,availability,status_message,timezone,until_time,created_at FROM team_members WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []TeamMember;for rows.Next(){var e TeamMember;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Role,&e.Availability,&e.StatusMessage,&e.Timezone,&e.UntilTime,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
