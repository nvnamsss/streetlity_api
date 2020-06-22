package rmaintenance

import (
	"errors"
	"net/http"
	"streelity/v1/model/maintenance"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service maintenance.Maintenance
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.URL.Query())
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		id := p.GetIntFirstOrDefault("Id")
		if service, e := maintenance.ServiceById(id); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

func QueryService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service maintenance.Maintenance
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.QueryMaintenanceValidate(req)
	p.First = stage
	res.Error(p.Run())
	if res.Status {
		c := p.GetInt("Case")[0]
		switch c {
		case 1:
			id := p.GetInt("Id")[0]
			if service, e := maintenance.ServiceById(id); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 2:
			lat := p.GetFloat("Lat")[0]
			lon := p.GetFloat("Lon")[0]
			if service, e := maintenance.ServiceByLocation(lat, lon); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 3:
			address := p.GetString("Address")[0]
			if service, e := maintenance.ServiceByAddres(address); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		}
	}
	sres.WriteJson(w, res)
}

func AllServices(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []maintenance.Maintenance
	}
	res.Status = true

	if services, e := maintenance.AllServices(); e != nil {
		res.Error(e)
	} else {
		res.Services = services
	}

	sres.WriteJson(w, res)
}

func CreateService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service maintenance.MaintenanceUcf
	}
	res.Status = true

	p := pipeline.NewPipeline()
	req.ParseForm()
	stage := stages.AddingServiceValidateStage(req)
	nameStage := pipeline.NewStage(func() (str struct {
		Name string
	}, e error) {
		form := req.PostForm
		names, ok := form["name"]
		if !ok {
			return str, errors.New("name param is missing")
		}
		str.Name = names[0]
		return
	})
	stage.NextStage(nameStage)
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloatFirstOrDefault("Lat")
		lon := p.GetFloatFirstOrDefault("Lon")
		address := p.GetStringFirstOrDefault("Address")
		note := p.GetStringFirstOrDefault("Note")
		name := p.GetStringFirstOrDefault("Name")
		images := p.GetString("Images")
		var ucf maintenance.MaintenanceUcf
		ucf.Lat = float32(lat)
		ucf.Lon = float32(lon)
		ucf.Address = address
		ucf.Note = note
		ucf.Name = name
		ucf.SetImages(images...)

		if service, e := maintenance.CreateUcf(ucf); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

func SetOwner(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}
	p := pipeline.NewPipeline()
	stage := stages.SetOwnerValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		service_id := p.GetInt("ServiceId")[0]
		owner := p.GetString("Owner")[0]
		values := make(map[string]string)
		values["owner"] = owner
		maintenance.UpdateService(service_id, values)
	}
	sres.WriteJson(w, res)
}

func ServiceInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []maintenance.Maintenance
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

		res.Services = maintenance.ServicesInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func HandleService(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/maintenance").Subrouter()

	s.HandleFunc("/", CreateService).Methods("POST")
	s.HandleFunc("/", GetService).Methods("GET")
	s.HandleFunc("/query", QueryService).Methods("GET")
	s.HandleFunc("/owner", SetOwner).Methods("POST")
	s.HandleFunc("/all", AllServices).Methods("GET")
	s.HandleFunc("/create", CreateService).Methods("POST")
	s.HandleFunc("/range", ServiceInRange).Methods("GET")
	return s
}
