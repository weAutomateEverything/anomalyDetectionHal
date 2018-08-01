package main

import (
	"os"
	"github.com/weAutomateEverything/mockXray"
	"github.com/weAutomateEverything/go2hal/database"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/weAutomateEverything/AnomalyDetectionHal/detector"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
			"io/ioutil"
	"os/signal"
	"syscall"
	"fmt"
)

func main() {
	if os.Getenv("XRAY_URL") != "" {
		var ss *sampling.LocalizedStrategy
		if os.Getenv("XRAY_SAMPLING_RULES") != "" {
			ss = getSampleStratergyFromFile()
		} else {
			ss = getDefaultSampleStratergy()
		}

		//XRAY
		xray.Configure(xray.Config{
			DaemonAddr:       os.Getenv("XRAY_URL"), // default
			LogLevel:         "info",                // default
			ServiceVersion:   "1.2.3",
			SamplingStrategy: ss,
		})

	} else {
		mockXray.StartMockXrayServer()
	}

	db := database.NewConnection()

	store := detector.NewDataStore(db)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	service := detector.NewService(store)
	service = detector.NewXray(service)
	service = detector.NewLoggingService(logger, service)
	service = detector.NewPrometheus(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "anomaly_detector",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{"method"}),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "anomaly_detector",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "anomaly_detector",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}), service)

	mux := http.NewServeMux()
	httpLogger := log.With(logger, "component", "http")
	mux.Handle("/api/anomaly/",detector.NewTransport(service,httpLogger))
	mux.Handle("/api/metrics", promhttp.Handler())
	mux.Handle("/api/swagger.json", swagger{})

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", ":8005", "msg", "listening")
		errs <- http.ListenAndServe(":8005", xray.Handler(xray.NewFixedSegmentNamer("anomaly_detector"), accessControl(mux)))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <- errs)

}


func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func getSampleStratergyFromFile() *sampling.LocalizedStrategy {
	ss, err := sampling.NewLocalizedStrategyFromFilePath(os.Getenv("XRAY_SAMPLING_RULES"))
	if err != nil {
		panic(err)
	}
	return ss
}

func getDefaultSampleStratergy() *sampling.LocalizedStrategy {
	ss, err := sampling.NewLocalizedStrategy()
	if err != nil {
		panic(err)
	}
	return ss
}

type swagger struct {
	http.Handler
}

func (h swagger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("swagger.json")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write(b)
	}
}

