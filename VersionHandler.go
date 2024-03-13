package routes

import (
	"encoding/json"
	"net/http"

	config "github.com/CodeNamor/Config"
	log "github.com/sirupsen/logrus"
)

//Version global string, it is exposed in order main to load the full API version
//var Version string

type response struct {
	Version      string        `json:"version"`
	Commit       string        `json:"commit"`
	BuildNumber  string        `json:"buildNumber"`
	Env          string        `json:"env"`
	ConfigHash   string        `json:"configHash"`
	ServicesUsed []serviceInfo `json:"servicesUsed,omitempty"`
}

type serviceInfo struct {
	Name                     string
	URL                      string
	AuthEnvironmentVariable  string
	EndPoints                config.EndpointMap
	ComponentConfigOverrides config.ComponentConfigs
}

const warningMessage = `Attempted to create version info with wrong config type (should be *config.Config) handed to the func 'VersionHandler'. Version endpoint will not have service information.`

func VersionHandler(v, c, bn, e string, interfacesForExpansion ...interface{}) http.HandlerFunc {
	if len(interfacesForExpansion) == 1 {
		if configData, ok := interfacesForExpansion[0].(*config.Config); ok {
			return versionHandlerWithServicesUsed(v, c, bn, e, configData.ServiceConfigs)
		} else {
			log.Warning(warningMessage)
		}
	}

	return standardVersionHandler(v, c, bn, e)
}

func standardVersionHandler(v, c, bn, e string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(response{Version: v, Commit: c, BuildNumber: bn, Env: e, ConfigHash: config.HashCode()})
	}
}

func versionHandlerWithServicesUsed(v, c, bn, e string, servicesMap config.ServicesMap) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		serviceInfos := convertServicesMapToArrayOfServiceInfo(servicesMap)
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(response{Version: v, Commit: c, BuildNumber: bn, Env: e, ConfigHash: config.HashCode(), ServicesUsed: serviceInfos})
	}
}

func convertServicesMapToArrayOfServiceInfo(servicesMap config.ServicesMap) []serviceInfo {
	var serviceInfos []serviceInfo
	for _, v := range servicesMap {
		info := serviceInfo{
			Name:                     v.Name,
			URL:                      v.URL,
			AuthEnvironmentVariable:  v.AuthEnvironmentVariable,
			EndPoints:                v.EndPoints,
			ComponentConfigOverrides: v.ComponentConfigOverrides,
		}
		serviceInfos = append(serviceInfos, info)
	}
	return serviceInfos
}
