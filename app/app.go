package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amiwx/p4_runner_golang/app/handler"
	"github.com/amiwx/p4_runner_golang/app/model"
	"github.com/amiwx/p4_runner_golang/config"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// main app stuff
type App struct {
	Router   *mux.Router
	DB       *gorm.DB
	P4Config *config.P4Config
}

func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Name,
		config.DB.Port,
	)

	var db *gorm.DB
	var err error
	switch config.DB.Dialect {
	case "postgresql":
		db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
		if err != nil {
			log.Fatal("Could not connect to database... ", err)
		}
	default:
		log.Println("dialect/driver incorrect or not found... using sqlite default DB")
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			log.Fatal("Could not connect to database... ", err)
		}
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.P4Config = config.P4
	a.setRouters()
}

func (a *App) setRouters() {
	// Runs Routers
	// a.Get("/start_run", a.StartRun)
	// a.Get("/update_run", a.UpdateRun)
	a.Get("/runs", a.GetAllRuns)
	a.Get("/runs/{runid}", a.GetRun)
	a.Get("/start_run", a.StartRun)
	a.Get("/show_runs", a.ShowRuns)

	// Posits Routers
	// a.Get("/posits", a.GetAllPosits)
	// a.Get("/posits/{positid}", a.GetPosit)
	// a.Get("/serve_posit", a.ServePosit)

}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) GetAllRuns(w http.ResponseWriter, r *http.Request) {
	handler.GetAllRuns(a.DB, w, r)
}

func (a *App) GetRun(w http.ResponseWriter, r *http.Request) {
	handler.GetRun(a.DB, w, r)
}

func (a *App) StartRun(w http.ResponseWriter, r *http.Request) {
	handler.StartRun(a.DB, w, r, a.P4Config)
}
func (a *App) ShowRuns(w http.ResponseWriter, r *http.Request) {
	handler.ShowRuns(a.DB, w, r)
}

// not related to other Run
func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
}
