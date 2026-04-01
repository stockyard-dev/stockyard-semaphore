package store
import("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{*sql.DB}
type Member struct{ID int64 `json:"id"`;Name string `json:"name"`;Status string `json:"status"`;Emoji string `json:"emoji"`;Note string `json:"note"`;AutoReturn string `json:"auto_return"`;UpdatedAt time.Time `json:"updated_at"`}
func Open(d string)(*DB,error){os.MkdirAll(d,0755);dsn:=filepath.Join(d,"semaphore.db")+"?_journal_mode=WAL&_busy_timeout=5000";db,err:=sql.Open("sqlite",dsn);if err!=nil{return nil,fmt.Errorf("open: %w",err)};db.SetMaxOpenConns(1);migrate(db);return &DB{db},nil}
func migrate(db *sql.DB){db.Exec(`CREATE TABLE IF NOT EXISTS members(id INTEGER PRIMARY KEY AUTOINCREMENT,name TEXT NOT NULL UNIQUE,status TEXT DEFAULT 'available',emoji TEXT DEFAULT '🟢',note TEXT DEFAULT '',auto_return TEXT DEFAULT '',updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)}
func(db *DB)Upsert(m *Member)error{_,err:=db.Exec(`INSERT INTO members(name,status,emoji,note,auto_return,updated_at)VALUES(?,?,?,?,?,CURRENT_TIMESTAMP) ON CONFLICT(name) DO UPDATE SET status=excluded.status,emoji=excluded.emoji,note=excluded.note,auto_return=excluded.auto_return,updated_at=CURRENT_TIMESTAMP`,m.Name,m.Status,m.Emoji,m.Note,m.AutoReturn);return err}
func(db *DB)List()([]Member,error){rows,_:=db.Query(`SELECT id,name,status,emoji,note,auto_return,updated_at FROM members ORDER BY name`);defer rows.Close();var out[]Member;for rows.Next(){var m Member;rows.Scan(&m.ID,&m.Name,&m.Status,&m.Emoji,&m.Note,&m.AutoReturn,&m.UpdatedAt);out=append(out,m)};return out,nil}
func(db *DB)Delete(id int64){db.Exec(`DELETE FROM members WHERE id=?`,id)}
func(db *DB)Stats()(map[string]interface{},error){var total,available int;db.QueryRow(`SELECT COUNT(*) FROM members`).Scan(&total);db.QueryRow(`SELECT COUNT(*) FROM members WHERE status='available'`).Scan(&available);return map[string]interface{}{"members":total,"available":available},nil}
