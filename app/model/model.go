package model

import (
	"gorm.io/gorm"
)

type Run struct {
	gorm.Model
	// RunID            uint `gorm:"primaryKey"`
	Version string
	// DateStarted      time.Time
	RunDir           string
	TotalFiles       int64
	APass            int64
	BPass            int64
	CPass            int64
	DPass            int64
	TotalPass        int64
	CompleteFailures int64
	TimeTaken        int64
}

type Posit struct {
	gorm.Model
	// PositID   uint   `gorm:"unique"`
	PositHash string `gorm:"uniqueIndex" gorm:"not null"`
	Filepath  string
	APass     bool
	BPass     bool
	CPass     bool
	DPass     bool
	Time      string
	Date      string
	Lat       float64
	Lon       float64
	LatString string `json:"latstr"`
	LonString string `json:"lonstr"`
	Brobs     []BROB //`gorm:"foreignKey:PositID;references:PositID"`
	Course    float64
	Speed     float64
	RPM       float64
	Slip      float64
	Distance  float64
	Weather   Weather       `gorm:"-"` //gorm:"embedded;type:text"`
	Stoppages ProcessedText `gorm:"-"` //gorm:"embedded;type:text"`
	Remarks   ProcessedText `gorm:"-"` //`gorm:"embedded;type:text"`
}

// additional structs for Posit

type BROB struct {
	gorm.Model
	PositID uint    //`gorm:"foreignKey:PositID"`
	Name    string  `json:"name"`
	Amount  float64 `json:"amount"`
}

type Weather struct {
	Text  string   `json:"text"`
	Wind  Wind     `gorm:"embedded;type:text" json:"wind"`
	Sea   SeaSwell `gorm:"embedded;type:text" json:"sea"`
	Swell SeaSwell `gorm:"embedded;type:text" json:"swell"`
}

type Wind struct {
	Beaufort  int      `json:"beaufort"`
	Direction string   `json:"direction"`
	Keywords  []string `gorm:"type:text" json:"keywords"`
}

type SeaSwell struct {
	Direction string   `json:"direction"`
	Height    float64  `json:"height"`
	Keywords  []string `gorm:"type:text" json:"keywords"`
}

type ProcessedText struct {
	Text     string   `json:"text"`
	Keywords []string `gorm:"type:text" json:"keywords"`
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&Run{})
	db.AutoMigrate(&Posit{})
	db.AutoMigrate(&BROB{})
	return db
}
