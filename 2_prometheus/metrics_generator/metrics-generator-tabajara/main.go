package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpURIRequests = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_requests_duration_seconds",
	Help:    "HTTP requests latency histogram",
	Buckets: []float64{0.3, 4, 35},
}, []string{
	"uri",    //200
	"method", //2
	"status", //3
	"server_name",
	"component_name",
	"component_version",
})

var httpAppRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_app_total",
	Help: "HTTP requests latency histogram",
}, []string{
	"status",             //3
	"device_app_version", //52
	"server_name",
	"component_name",
	"component_version",
})

var httpDeviceRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_device_total",
	Help: "HTTP requests latency histogram",
}, []string{
	"status",            //3
	"device_os_name",    //2
	"device_os_version", //98
	"server_name",
	"component_name",
	"component_version",
})

var httpRequestsReceived = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_received",
	Help: "HTTP requests received qty",
}, []string{
	"server_name",
	"component_name",
	"component_version",
})

var httpResponseSent = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_responses_sent",
	Help: "HTTP responses sent qty",
}, []string{
	"server_name",
	"component_name",
	"component_version",
})

var httpPendingRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "http_pending_requests",
	Help: "HTTP pending requests",
}, []string{
	"server_name",
	"component_name",
	"component_version",
})

//Accident accident
type Accident struct {
	ResourceName string `json:"resource",omitempty`
	Type         string `json:"type",omitempty`
	Value        string `json:"value",omitempty`
}

var serverName = ""
var componentName = ""
var componentVersion = ""
var activeAccidents = []Accident{}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	serverName0 := flag.String("server-name", "", "Emulated server name")
	componentName0 := flag.String("component-name", "mything", "Emulated component name")
	componentVersion0 := flag.String("component-version", "1.0.0", "Emulated component version")
	accidentResource0 := flag.String("accident-resource", "", "Accident resource names separated by space")
	accidentType0 := flag.String("accident-type", "none", "Accident type (none, latency, calls, errors)")
	accidentRatio0 := flag.String("accident-ratio", "1.0", "Accident proportion")
	flag.Parse()

	serverName = *serverName0
	componentName = *componentName0
	componentVersion = *componentVersion0

	accidentRatio := *accidentRatio0
	resourcesAffectedNames := strings.Split(*accidentResource0, " ")

	accidentType := *accidentType0

	for _, v := range resourcesAffectedNames {
		if accidentType != "none" {
			activeAccidents = append(activeAccidents, Accident{ResourceName: v, Type: accidentType, Value: accidentRatio})
		}
	}

	logrus.Infof("%s %s %s %s %s %s", componentName, componentVersion, resourcesAffectedNames, accidentType, accidentRatio)

	prometheus.MustRegister(httpURIRequests)
	prometheus.MustRegister(httpAppRequests)
	prometheus.MustRegister(httpDeviceRequests)
	prometheus.MustRegister(httpRequestsReceived)
	prometheus.MustRegister(httpResponseSent)
	prometheus.MustRegister(httpPendingRequests)

	go generateHTTPMetrics()

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/metrics-http", promhttp.Handler())
	router.Handle("/metrics-negocio", promhttp.Handler())
	router.HandleFunc("/surgery-accident", PostAccidents).Methods("POST")
	router.HandleFunc("/surgery-accident", DeleteAccidents).Methods("DELETE")
	err := http.ListenAndServe("0.0.0.0:3000", router)
	if err != nil {
		logrus.Errorf("Error while listening requests: %s", err)
		os.Exit(1)
	}
}

func generateHTTPMetrics() {
	logrus.Infof("Starting requests simulation to generate centralizador metrics...")
	var uris = generateItems("/resources/somegroup/item-", 30)
	var statuses = []string{"4xx", "2xx", "5xx"}
	var methods = []string{"POST", "GET"}
	var deviceOsName = []string{"ios", "android"}
	var deviceAppVersion = generateItems("v", 4)
	var deviceOsVersion = generateItems("v", 3)

	httpRequestsReceived.WithLabelValues(serverName, componentName, componentVersion).Inc()
	httpRequestsReceived.WithLabelValues(serverName, componentName, componentVersion).Inc()
	httpRequestsReceived.WithLabelValues(serverName, componentName, componentVersion).Inc()
	httpRequestsReceived.WithLabelValues(serverName, componentName, componentVersion).Inc()

	for {
		uri := getRandomElemNormal(uris)

		//change number of calls
		calls := int(getValueAccident("calls", 1.0, uri))
		//if calls < 1, call randomly proportional to probability ratio
		if calls < 1 {
			if calls > rand.Intn(100)/100.0 {
				calls = 0
				time.Sleep(100 * time.Millisecond) //avoid 100% CPU
			} else {
				calls = 1
			}
		}
		// logrus.Infof("calls http %d", calls)

		for i := 1; i <= calls; i++ {
			median := (hash(uri)%23 + 1) * 100
			requestTime := generateSample(float64(median), float64(median)/5) / 1000.0
			requestTime = math.Max(requestTime, 0)
			// logrus.Infof("calls http %s %d", uri, median)

			//change duration
			requestTime = getValueAccident("latency", requestTime, uri)

			status := getRandomElemNormal(statuses)

			statRatio := getValueAccident("errors", 1.0001, uri)
			if statRatio != 1.0001 {
				sr := rand.Intn(100)
				if sr <= int(statRatio)*100 {
					status = "5xx"
				}
			}

			//use POST in 30% of requests. the same URL will have always the same method
			m := 1
			if randomInt(int64(hash(uri)), 10) < 3 {
				m = 0
			}

			httpURIRequests.WithLabelValues(
				uri,
				methods[m],
				status,
				serverName,
				componentName,
				componentVersion).Observe(requestTime)

			httpAppRequests.WithLabelValues(
				status,
				getRandomElemNormal(deviceAppVersion),
				serverName,
				componentName,
				componentVersion).Inc()

			httpDeviceRequests.WithLabelValues(
				status,
				getRandomElem(deviceOsName),
				getRandomElemNormal(deviceOsVersion),
				serverName,
				componentName,
				componentVersion).Inc()

			httpRequestsReceived.WithLabelValues(serverName, componentName, componentVersion).Inc()
			httpResponseSent.WithLabelValues(serverName, componentName, componentVersion).Inc()
			httpPendingRequests.WithLabelValues(serverName, componentName, componentVersion).Set(float64(randomRangeNormal(0, 400)))

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func getValueAccident(accidentType string, defaultValue float64, resourceName string) float64 {
	if len(activeAccidents) > 0 {
		for _, v := range activeAccidents {
			re := regexp.MustCompile(v.ResourceName)
			if re.MatchString(resourceName) && accidentType == v.Type {
				// logrus.Infof("Using accident %s", v)
				ratio, err := strconv.ParseFloat(v.Value, 64)
				if err != nil {
					// logrus.Errorf("Error using accident %s", err.Error())
					return defaultValue
				}
				return float64(defaultValue) * ratio
			}
		}
	}
	return defaultValue
}

//DeleteAccidents handle reset accidents endpoint
func DeleteAccidents(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("DeleteAccidents r=%v", r)
	activeAccidents = []Accident{}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"reset ok. 😅\"}"))
}

//PostAccidents handle accidents endpoint
func PostAccidents(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("PostAccidents r=%v", r)
	accident := Accident{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		logrus.Infof(">>>>> %s", body)
		err = json.Unmarshal(body, &accident)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logrus.Infof("Creating accident %s", accident)
		// err := json.NewDecoder(r.Body).Decode(&Accident)
		w.Header().Set("Content-Type", "application/json")

		_, err := strconv.ParseFloat(accident.Value, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if accident.Type == "" {
			http.Error(w, "type is required", http.StatusInternalServerError)
			return
		} else if accident.Type != "none" && accident.Type != "latency" && accident.Type != "calls" && accident.Type != "errors" {
			http.Error(w, "type must be one of none|latency|calls|errors", http.StatusInternalServerError)
			return
		}

		if accident.ResourceName == "" {
			http.Error(w, "resource_name is required", http.StatusInternalServerError)
			return
		}

		newAccidents := []Accident{}
		for _, v := range activeAccidents {
			if v.ResourceName != accident.ResourceName {
				newAccidents = append(newAccidents, v)
			} else {
				logrus.Infof("Removing accident %s", v)
			}
		}
		logrus.Infof("Adding accident %s", accident)
		newAccidents = append(newAccidents, accident)
		activeAccidents = newAccidents
		w.Write([]byte("{\"message\":\"AARRRGGG! Crisis planted. Run for your life! 🧟\"}"))
		return
	}
}

func generateSample(median float64, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + median
}

func getRandomElem(items []string) string {
	return items[rand.Intn(len(items))]
}

func getRandomElemNormal(items []string) string {
	return items[randomRangeNormal(0, len(items)-1)]
}

func generateItems(prefix string, qtty int) []string {
	result := []string{}
	for i := 1; i <= qtty; i++ {
		result = append(result, fmt.Sprintf("%s%04d", prefix, i))
	}
	return result
}

func randomRangeNormal(min int, max int) int {
	median := float64((max - min) / 2)
	v := generateSample(median+0.5, median/3)
	mi := float64(min)
	ma := float64(max)
	return min + int(math.Max(float64(math.Min(float64(v), ma)), mi))
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func randomInt(seed int64, max int) int {
	s1 := rand.NewSource(seed)
	return rand.New(s1).Intn(max)
}
