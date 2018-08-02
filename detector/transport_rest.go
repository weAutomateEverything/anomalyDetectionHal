package detector

import (
	"net/http"
	"github.com/weAutomateEverything/go2hal/gokit"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"io/ioutil"
		"context"
	"strconv"
)

func NewTransport(s Service, logger kitlog.Logger, ) http.Handler {
	opts := gokit.GetServerOpts(logger, nil)

	addDataEndpoint := kithttp.NewServer(makeAddDataEndpoint(s), decodeAddData, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api/anomaly/{key} learning AddData
	//
	// Adds a data point for the key and returns an anomaly score
	//
	//
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: key
	//   in: path
	//   description: key for the data
	//   required: true
	//   type: integer
	// - name: message
	//   in: body
	//   description: floating point value
	//   required: true
	//   schema:
	//     type: float64
	// responses:
	//   '200':
	//     description: Success
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/anomaly/{key}", addDataEndpoint).Methods("POST")

	return r

}

func decodeAddData(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	req, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	val, err := strconv.ParseFloat(string(req),64)
	if err != nil {
		return nil, err
	}

	return &addDataRequest{
		value:val,
		key: vars["key"],
	}, nil





}
