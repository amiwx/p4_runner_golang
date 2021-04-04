package process

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/amiwx/p4_runner_golang/app/model"
	"github.com/amiwx/p4_runner_golang/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ProcessPosits(rundir string, db *gorm.DB, p4Config *config.P4Config) (model.Run, error) {

	var run model.Run

	pathArray, err := findPosits(rundir)
	if err != nil {
		return run, err
	} else if len(pathArray) < 1 {
		return run, fmt.Errorf("No posits found in %s", rundir)
	}

	requestUrl := fmt.Sprintf("%s:%s%s", p4Config.Host, p4Config.Port, p4Config.ParseURL)
	// requestUrl := "http://127.0.0.1:8001/v1/parse"
	requestUrl = strings.TrimSpace(requestUrl)
	client := makeClient()
	positDataArray := make([]model.Posit, 0, len(pathArray))

	// initialize some run values
	run.TotalFiles = int64(len(pathArray))
	run.RunDir = rundir

	//var positData model.Posit
	start := time.Now()
	for _, positPath := range pathArray {
		// // posit, err := loadPosit(positPath)
		// if err != nil {
		// 	log.Printf("failed to open posit: %s", err)
		// } else {
		// }
		responseData, err := makeP4Request(positPath, requestUrl, client)
		if err != nil {
			log.Printf("failed to make p4_request: %v", err)
			continue
		}
		if run.Version == "" {
			run.Version = responseData.Version
		}

		positData := responseData.Data
		positData.Filepath = positPath
		didPositPass(&positData) // modifies reference of positData
		incrementRunData(&positData, &run)

		// creates hash for uniqueness purposes
		hash := md5.Sum([]byte(positPath))
		positData.PositHash = hex.EncodeToString(hash[:])

		positDataArray = append(positDataArray, positData)

	}

	// adds posits to database
	log.Printf("Inserting database with %d posit(s)", len(pathArray))
	err = db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "posit_hash"}},
			UpdateAll: true,
		},
	).CreateInBatches(&positDataArray, 1000).Error
	if err != nil {
		return run, err
	}

	run.TimeTaken = int64(time.Since(start).Seconds())

	return run, nil
}

func incrementRunData(posit *model.Posit, run *model.Run) {
	if posit.APass {
		run.APass++
	}
	if posit.BPass {
		run.BPass++
	}
	if posit.CPass {
		run.CPass++
	}
	if posit.DPass {
		run.DPass++
	}

	if posit.CPass && posit.BPass && posit.APass {
		run.TotalPass++
	} else if !posit.CPass && !posit.BPass && !posit.APass {
		// complete failures is posit has failed completely
		log.Printf("posit %s is a complete failure", posit.Filepath)
		run.CompleteFailures++
	}

	// switch {
	// case posit.APass:
	// 	run.APass++
	// 	fallthrough
	// case posit.BPass:
	// 	run.BPass++
	// 	fallthrough
	// case posit.CPass:
	// 	run.CPass++
	// 	fallthrough
	// case posit.DPass:
	// 	run.DPass++
	// }

}

func didPositPass(posit *model.Posit) {
	posit.APass = testSectionA(posit)
	posit.BPass = testSectionB(posit)
	posit.CPass = testSectionC(posit)
	posit.DPass = testSectionD(posit)
}

func testSectionA(posit *model.Posit) bool {

	keys := make([]string, 0, 3)

	if posit.Date == "" {
		keys = append(keys, "DATE")
	}

	if posit.Lat == 0.0 {
		keys = append(keys, "LAT")
	}

	if posit.Lon == 0.0 {
		keys = append(keys, "LON")
	}

	if len(keys) > 0 {
		log.Printf("Section A check failed, missing keys: %v", keys)
		return false
	} else {
		return true
	}

}

func testSectionB(posit *model.Posit) bool {

	if posit.Brobs != nil && len(posit.Brobs) >= 2 {
		return true
	} else {
		log.Printf("Section B check failed: %v", posit.Brobs)
		return false
	}
}

func testSectionC(posit *model.Posit) bool {
	keys := make([]string, 0, 3)

	if posit.Course == 0.0 {
		keys = append(keys, "COURSE")
	}

	if posit.Speed == 0.0 {
		keys = append(keys, "SPEED")
	}

	if posit.Distance == 0.0 {
		keys = append(keys, "DISTANCE")
	}

	if len(keys) > 0 {
		log.Printf("Section C check failed, missing keys: %v", keys)
		return false
	} else {
		return true
	}
}

func testSectionD(posit *model.Posit) bool {
	// var defaultWeather = model.Weather{}
	// if posit.Weather != defaultWeather {
	// 	return true
	// } else {
	// 	return false
	// }

	// not sure how to test this exactly
	return false
}
