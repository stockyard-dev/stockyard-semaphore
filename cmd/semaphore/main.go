package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-semaphore/internal/server";"github.com/stockyard-dev/stockyard-semaphore/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./semaphore-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("semaphore: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Semaphore — Self-hosted team availability board\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("semaphore: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
