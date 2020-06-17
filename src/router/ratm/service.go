package ratm

import (
	"net/http"
	"streelity/v1/model/atm"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service atm.Atm
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.URL.Query())
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		id := p.GetIntFirstOrDefault("Id")
		if service, e := atm.ServiceById(id); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

func AllServices(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []atm.Atm
	}
	res.Status = true

	if services, e := atm.AllServices(); e != nil {
		res.Error(e)
	} else {
		res.Services = services
	}

	sres.WriteJson(w, res)
}

func CreateService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service atm.AtmUcf
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.AddingServiceValidateStage(req)
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloatFirstOrDefault("Lat")
		lon := p.GetFloatFirstOrDefault("Lon")
		address := p.GetStringFirstOrDefault("Address")
		note := p.GetStringFirstOrDefault("Note")
		images := p.GetString("Images")
		var ucf atm.AtmUcf
		ucf.Lat = float32(lat)
		ucf.Lon = float32(lon)
		ucf.Address = address
		ucf.Note = note
		ucf.SetImages(images...)

		if service, e := atm.CreateUcf(ucf); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

func ServiceInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []atm.Atm
	}
	res.Status = true
	pipe := pipeline.NewPipeline()
	stage := stages.InRangeServiceValidateStage(req)
	pipe.First = stage

	res.Error(pipe.Run())

	if res.Status {
		lat := pipe.GetFloatFirstOrDefault("Lat")
		lon := pipe.GetFloatFirstOrDefault("Lon")
		max_range := pipe.GetFloatFirstOrDefault("Range")
		var location r2.Point = r2.Point{X: lat, Y: lon}

		res.Services = atm.ServicesInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func HandleService(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/atm").Subrouter()

	s.HandleFunc("/", CreateService).Methods("POST")
	s.HandleFunc("/", GetService).Methods("GET")
	s.HandleFunc("/all", AllServices).Methods("GET")
	s.HandleFunc("/upvote", UpvoteUnconfirmed).Methods("POST")
	s.HandleFunc("/create", CreateService).Methods("POST")

	return s
}