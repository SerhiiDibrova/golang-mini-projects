package initializers

import (
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

var styles = []string{"style1", "style2", "style3"}
var mutex = &sync.Mutex{}

func InitializeRoutes(router *mux.Router) {
	router.HandleFunc("/getStyles", GetStyles).Methods("GET")
}

func GetStyles(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	jsonResponse, err := json.Marshal(styles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	router := mux.NewRouter()
	InitializeRoutes(router)
	http.ListenAndServe(":8000", router)
}