package model

import (
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
)

type Atm struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
}

func (Atm) TableName() string {
	return "atm"
}

func (s Atm) GetLocation() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

func AllAtms() []Atm {
	var services []Atm
	Db.Find(&services)

	return services
}

func AtmById(id int64) Atm {
	var service Atm
	Db.Find(&service, id)

	return service
}

func AllAtmsInRange(circle spatial.Circle) []Atm {
	var services []Atm
	Db.Find(&services)

	return services
}
