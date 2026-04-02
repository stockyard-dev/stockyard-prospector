package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-prospector/internal/server";"github.com/stockyard-dev/stockyard-prospector/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9700"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./prospector-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("prospector: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Prospector — Self-hosted sales pipeline tracker\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("prospector: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
