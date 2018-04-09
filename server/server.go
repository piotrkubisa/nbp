package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/karolgorecki/nbp/svc"

	"github.com/julienschmidt/httprouter"
)

// RegisterHandlers does something
func RegisterHandlers() *httprouter.Router {
	rt := httprouter.New()
	rt.GET("/:date/:type/:code", errorHandler(IndexHandler))

	rt.NotFound = ntHandler

	fmt.Println("Running on: http://localhost:" + os.Getenv("PORT"))
	return rt
}

// IndexHandler Does something
func IndexHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	rDate := p.ByName("date")
	rType := p.ByName("type")
	rCode := p.ByName("code")

	var (
		res svc.Query
		err error
	)

	// Is the type OK?
	switch rType {
	case "avg":
		res, err = svc.Average(rDate, rCode)
		if err != nil {
			handleOutput(w, http.StatusBadRequest, err.Error)
			return nil
		}
	case "both":
		res, err = svc.Both(rDate, rCode)
		if err != nil {
			handleOutput(w, http.StatusBadRequest, err.Error)
			return nil
		}
	default:
		handleOutput(w, http.StatusBadRequest, "Given type is wrong. Use 'avg' or 'both'")
		return nil
	}

	handleOutput(w, http.StatusOK, res)
	return nil
}

// handleOutput handles the response for each endpoint.
// It follows the JSEND standard for JSON response.
// See https://labs.omniti.com/labs/jsend
func handleOutput(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)

	success := false
	if code == 200 {
		success = true
	}

	// JSend has three possible statuses: success, fail and error
	// In case of error, there is no data sent, only an error message.
	status := "success"
	msgType := "data"
	if !success {
		status = "error"
		msgType = "message"
	}

	res := map[string]interface{}{"status": status}
	if data != nil {
		res[msgType] = data
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
}

func ntHandler(w http.ResponseWriter, r *http.Request) {
	handleOutput(w, http.StatusNotFound, "The resource you're looking was not found")
}
