package fuel

import (
	"errors"
	"log"
	"math"
	"streelity/v1/model"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

//FuelUcf representation the Fuel service which is confirmed
type Fuel struct {
	model.Service
	// Id  int64
	// Lat float32 `gorm:"column:lat"`
	// Lon float32 `gorm:"column:lon"`
}

var services spatial.RTree

//Determine table name
func (Fuel) TableName() string {
	return "fuel"
}

func (s Fuel) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllFuels query all fuel services
func AllFuels() []Fuel {
	var services []Fuel
	model.Db.Find(&services)

	return services
}

//AddFuel add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddFuel(s Fuel) (e error) {
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Fuel{}).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e := model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "add fuel", e.Error())
	}

	return
}

func queryFuel(s Fuel) (service Fuel, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query fuel", e.Error())
	}

	return
}

//FuelByService get fuel by provide model.Service
func FuelByService(s model.Service) (services Fuel, e error) {
	services.Service = s
	return queryFuel(services)
}

//FuelById query the fuel service by specific id
func FuelById(id int64) (service Fuel, e error) {
	if e = model.Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
		return service, errors.New("Problem occured when query")
	}

	return
}

//ToiletByIds query the toilets service by specific id
func FuelByIds(ids ...int64) (services []Fuel) {
	for _, id := range ids {
		s, e := FuelById(id)
		if e != nil {
			continue
		}

		services = append(services, s)
	}

	return
}

func distance(p1 r2.Point, p2 r2.Point) float64 {
	x := math.Pow(p1.X-p2.X, 2)
	y := math.Pow(p1.Y-p2.Y, 2)
	return math.Sqrt(x + y)
}

//FuelsInRange query the fuel services which is in the radius of a location
func FuelsInRange(p r2.Point, max_range float64) []Fuel {
	var result []Fuel = []Fuel{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isFuel := item.(Fuel)
			if isFuel && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

func (s Fuel) AfterCreate(scope *gorm.Scope) (e error) {
	if e = services.AddItem(s); e != nil {
		log.Println("[Database]", "After create fuel", e.Error())
	}

	return
}

func LoadService() {
	log.Println("[Fuel]", "Loading service")

	fuels := AllFuels()
	for _, atm := range fuels {
		services.AddItem(atm)
	}
}

func init() {
	model.OnConnected.Subscribe(LoadService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadService)
	})
}