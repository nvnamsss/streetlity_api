package stages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

/*Containing commom stage for pipeline, in order to reduce the effort and lines of code to implement
the pipeline in request handle*/

//ServiceValidateStage create the validated stage for adding a new service
func AddingServiceValidateStage(req *http.Request) *pipeline.Stage {
	s := pipeline.NewStage(func() (str struct {
		Address string
		Note    string
		Images  []string
	}, e error) {
		form := req.PostForm
		location, locationOk := form["location"]
		_, addressOk := form["address"]
		if !locationOk {
			return str, errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return str, errors.New("location param must have 2 values")
			}
		}

		if !addressOk {
			return str, errors.New("address param is missing")
		} else {
			str.Address = form["address"][0]
		}

		_, ok := form["note"]
		if ok {
			str.Note = form["note"][0]
		}

		str.Images = form["images"]

		return
	})

	return s
}

//AddingServiceParsingStage create the parsing stage for adding a new service
func AddingServiceParsingStage(req *http.Request) *pipeline.Stage {
	s := pipeline.NewStage(func() (str struct {
		Lat float64
		Lon float64
	}, e error) {
		var (
			latErr error
			lonErr error
		)
		form := req.PostForm

		str.Lat, latErr = strconv.ParseFloat(form["location"][0], 64)
		str.Lon, lonErr = strconv.ParseFloat(form["location"][1], 64)
		if latErr != nil {
			return str, errors.New("cannot parse location[0] to float")
		}
		if lonErr != nil {
			return str, errors.New("cannot parse location[1] to float")
		}
		return str, nil
	})

	return s
}

func InRangeServiceValidateStage(req *http.Request) *pipeline.Stage {
	validateParamsStage := pipeline.NewStage(func() error {
		query := req.URL.Query()

		location, locationOk := query["location"]
		if !locationOk {
			return errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return errors.New("location param must have 2 values")
			}
		}

		_, rangeOk := query["range"]
		if !rangeOk {
			return errors.New("range param is missing")
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() (str struct {
		Lat   float64
		Lon   float64
		Range float64
	}, e error) {
		var (
			latErr   error
			lonErr   error
			rangeErr error
		)
		query := req.URL.Query()

		str.Lat, latErr = strconv.ParseFloat(query["location"][0], 64)
		str.Lon, lonErr = strconv.ParseFloat(query["location"][1], 64)
		str.Range, rangeErr = strconv.ParseFloat(query["range"][0], 64)
		if latErr != nil {
			return str, errors.New("cannot parse location[0] to float")
		}
		if lonErr != nil {
			return str, errors.New("cannot parse location[1] to float")
		}
		if rangeErr != nil {
			return str, errors.New("cannot parse range to float")
		}
		return str, nil
	})

	validateParamsStage.NextStage(parseValueStage)

	return validateParamsStage
}

func QueryServiceValidateStage(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		Id      int64
		Lat     float64
		Lon     float64
		Address string
		Case    int
	}, e error) {
		query := req.URL.Query()

		id, ok := query["id"]
		if ok {
			if id, e := strconv.ParseInt(id[0], 10, 64); e != nil {
				return str, errors.New("cannot parse id to int")
			} else {
				str.Id = id
				str.Case = 1
				return str, nil
			}
		}

		lat, latOk := query["lat"]
		lon, lonOk := query["lon"]
		if latOk && lonOk {
			lat, latOk := strconv.ParseFloat(lat[0], 64)
			lon, lonOk := strconv.ParseFloat(lon[0], 64)

			if latOk == nil && lonOk == nil {
				str.Lat = lat
				str.Lon = lon
				str.Case = 2
				return str, nil
			}
		}

		addresses, addressOk := query["address"]

		if !addressOk {
			return str, errors.New("required at least one param id / lat - lon / address")
		}

		str.Address = addresses[0]
		str.Case = 3
		return
	})

	return stage
}

func ReviewValidateStage(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		ServiceId int64
		Commenter int64
		Score     float32
		Body      string
	}, e error) {
		form := req.PostForm

		if _, ok := form["service_id"]; !ok {
			return str, errors.New("service_id param is missing")
		}

		if _, ok := form["commenter"]; !ok {
			return str, errors.New("commenter param is missing")
		}

		if _, ok := form["score"]; !ok {
			return str, errors.New("score param is missing")
		}

		if _, ok := form["body"]; !ok {
			return str, errors.New("body param is missing")
		}

		if i, e := strconv.ParseInt(form["service_id"][0], 10, 64); e == nil {
			str.ServiceId = i
		} else {
			return str, errors.New("service_id cannot parse to int")
		}

		if i, e := strconv.ParseInt(form["commenter"][0], 10, 64); e == nil {
			str.Commenter = i
		} else {
			return str, errors.New("commenter cannot parse to int")
		}

		if f, e := strconv.ParseFloat(form["score"][0], 32); e == nil {
			str.Score = float32(f)
		} else {
			return str, errors.New("score cannot parse to float")
		}

		str.Body = form["body"][0]

		return
	})

	return stage
}

func ReviewByOrderValidate(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		ServiceId int64
		Order     int64
	}, e error) {
		query := req.URL.Query()

		_, ok := query["review_id"]
		if !ok {
			return str, errors.New("review_id param is missing")
		}

		_, ok = query["order"]
		if !ok {
			return str, errors.New("order param is missing")
		}

		review_id, e := strconv.ParseInt(query["review_id"][0], 10, 64)
		if e != nil {
			return str, errors.New("review_id cannot parse to int64")
		}

		order, e := strconv.ParseInt(query["order"][0], 10, 64)
		if e != nil {
			return str, errors.New("order cannot parse to int64")
		}

		str.ServiceId = review_id
		str.Order = order
		return
	})
	return stage
}

func UpdateReviewValidateStage(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		ReviewId int64
		NewBody  string
	}, e error) {
		form := req.PostForm
		review_ids, ok := form["review_id"]
		if !ok {
			return str, errors.New("review_id param is missing")
		}

		bodies, ok := form["new_body"]
		if !ok {
			return str, errors.New("new_body param is missing")
		}

		review_id, e := strconv.ParseInt(review_ids[0], 10, 64)
		if e != nil {
			return str, errors.New("review_id param cannot parse to int64")
		}

		str.ReviewId = review_id
		str.NewBody = bodies[0]
		return
	})

	return stage
}

func ReviewIdValidate(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		ReviewId int64
	}, e error) {
		form := req.PostForm
		review_ids, ok := form["review_id"]
		if !ok {
			return str, errors.New("review_id param is missing")
		}

		review_id, e := strconv.ParseInt(review_ids[0], 10, 64)
		if e != nil {
			return str, errors.New("review_id param cannot parse to int64")
		}

		str.ReviewId = review_id
		return
	})

	return stage
}