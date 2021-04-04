package handler

import (
	"net/http"
	"strconv"

	"github.com/amiwx/p4_runner_golang/app/model"
	"github.com/amiwx/p4_runner_golang/app/process"
	"github.com/amiwx/p4_runner_golang/config"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetAllRuns(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	runs := []model.Run{}
	db.Scopes(Paginate(r)).Find(&runs)
	respondJSON(w, http.StatusOK, runs)
}

func GetRun(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	runid := vars["runid"]
	run := getRunOr404(db, runid, w, r)
	if run == nil {
		return
	}
	respondJSON(w, http.StatusOK, run)
}

func StartRun(db *gorm.DB, w http.ResponseWriter, r *http.Request, p4Config *config.P4Config) {

	vars := r.URL.Query()

	rundir := vars.Get("rundir")
	if rundir == "" {
		respondError(w, http.StatusBadRequest, "param rundir is not set or empty")
		return
	}
	run, err := process.ProcessPosits(rundir, db, p4Config)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// fmt.Printf("%v+", run)
	db.Create(&run)
	respondJSON(w, http.StatusOK, run)

}

func getRunOr404(db *gorm.DB, runid string, w http.ResponseWriter, r *http.Request) *model.Run {
	run := model.Run{}
	runid_int, err := strconv.Atoi(runid)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	if err := db.First(&run, gorm.Model{ID: uint(runid_int)}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &run
}
