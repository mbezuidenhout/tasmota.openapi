/*
 * Sample API
 *
 * Device manager for Tasmota devices via MQTT [Source](https://github.com/mbezuidenhout/tdm).
 *
 * API version: 0.1.0
 * Contact: marius.bezuidenhout@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mbezuidenhout/tasmota.mqtt.device.manager/v2"
)

type MQTTOptions struct {
	Host        string
	Username    string
	Password    string
	CustomTopic string
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type ManagerWithTimes struct {
	tasmota.Manager
	created    time.Time
	lastUpdate time.Time
}

type Routes []Route

const Timeout = 5 * time.Minute

var managers map[string]*ManagerWithTimes
var routeerror interface{}

// Cleanup connections and return the number of connections that has been closed
func CleanupConnections() int {
	i := 0
	for key, element := range managers {
		if time.Now().After(element.lastUpdate.Add(Timeout)) {
			element.MQTTclient.Disconnect(1)
			delete(managers, key)
			i++
		}
	}
	// Check connections that has not recevied messages and close them
	return i
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/swaggerui")))

	managers = make(map[string]*ManagerWithTimes)

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func GetRouteError() interface{} {
	return routeerror
}

func recoverError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		routeerror = r
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

var routes = Routes{
	Route{
		Name:        "Index",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v3/",
		HandlerFunc: Index,
	},

	Route{
		Name:        "MqttConnectPost",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/v3/mqtt/connect",
		HandlerFunc: MQTTConnectPost,
	},

	Route{
		Name:        "MqttDisconnectGet",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v3/mqtt/disconnect",
		HandlerFunc: MQTTDisconnectGet,
	},

	Route{
		Name:        "MqttGet",
		Method:      strings.ToUpper("Get"),
		Pattern:     "/v3/mqtt",
		HandlerFunc: MQTTGet,
	},

	Route{
		Name:        "DevicesPost",
		Method:      strings.ToUpper("Post"),
		Pattern:     "/v3/devices",
		HandlerFunc: DevicesPost,
	},
}
