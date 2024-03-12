package routes

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	config "github.com/CodeNamor/Config"
	"github.com/stretchr/testify/assert"
)

func TestVersionHandler_NoServicesUsed(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	v := "0.1.test"
	commit := "321654"
	build := "12"
	env := "Local"
	configHash := "0399c455d71ab329b515f55b67f37ab9"

	config.NewHashCode(configHash)

	handler := http.HandlerFunc(VersionHandler(v, commit, build, env))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Handler returned status %v", recorder.Code)
	}

	want := `{"version":"` + v + `","commit":"` + commit + `","buildNumber":"` + build + `","env":"` + env + `","configHash":"` + configHash + `"}`
	assert.Equal(t, want, strings.TrimSpace(recorder.Body.String()))
}

func TestVersionHandler_WithServicesUsed(t *testing.T) {
	t.Run("It should handle a pointer to a config", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/version", nil)
		if err != nil {
			t.Fatal(err)
		}

		v := "0.1.test"
		commit := "321654"
		build := "12"
		env := "Local"
		configHash := "0399c455d71ab329b515f55b67f37ab9"
		servicesUsed := getServicesUsedJsonString()

		space := regexp.MustCompile(`\s+`)
		servicesUsedSpacesRemoved := space.ReplaceAllString(servicesUsed, "")

		config.NewHashCode(configHash)

		serviceConfig := getServiceConfig()
		bigOlConfig := getConfig(serviceConfig)

		handler := VersionHandler(v, commit, build, env, &bigOlConfig)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Handler returned status %v", recorder.Code)
		}

		want := `{"version":"` + v + `","commit":"` + commit + `","buildNumber":"` + build + `","env":"` + env + `","configHash":"` + configHash + `","servicesUsed":` + servicesUsedSpacesRemoved + `}`
		assert.Equal(t, want, strings.TrimSpace(recorder.Body.String()))
	})

	t.Run("It should fallback to older behavior if we don't get a pointer to config", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/version", nil)
		if err != nil {
			t.Fatal(err)
		}

		v := "0.1.test"
		commit := "321654"
		build := "12"
		env := "Local"
		configHash := "0399c455d71ab329b515f55b67f37ab9"

		config.NewHashCode(configHash)

		serviceConfig := getServiceConfig()
		bigOlConfig := getConfig(serviceConfig)

		handler := VersionHandler(v, commit, build, env, bigOlConfig)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Handler returned status %v", recorder.Code)
		}

		want := `{"version":"` + v + `","commit":"` + commit + `","buildNumber":"` + build + `","env":"` + env + `","configHash":"` + configHash + `"}`
		assert.Equal(t, want, strings.TrimSpace(recorder.Body.String()))
	})
}

func getConfig(serviceConfig config.ServiceConfig) config.Config {
	return config.Config{
		ServiceConfigs: map[string]*config.ServiceConfig{
			"AbsMemberComposite": &serviceConfig,
		},
	}
}

func getServiceConfig() config.ServiceConfig {
	return config.ServiceConfig{
		Name:                    "AbsMemberComposite",
		URL:                     "https://test-int-api-gw.centene.com/absmemberinquiry",
		AuthRequired:            false,
		AuthEnvironmentVariable: "ABS_MBR_KEY",
		AuthCredentials: config.AuthCredentials{
			KeyComponent1: "aaa",
			KeyComponent2: "bb",
			Euuid:         "dd",
		},
		AuthKey: "",
		EndPoints: map[string]*config.EndpointConfig{
			"MemberSearch": {
				Name: "MemberSearch",
				Path: "/inq/abs/V1",
			},
		},
		ComponentConfigOverrides: config.ComponentConfigs{
			ServiceLogging: config.ServiceLoggingConfig{
				LogCallDuration: 0,
			},
			Client: config.ClientConfig{
				Timeout:             0,
				IdleConnTimeout:     0,
				MaxIdleConnsPerHost: 0,
				MaxConnsPerHost:     0,
				MaxRetries:          0,
				DisableCompression:  0,
				InsecureSkipVerify:  0,
				CABundlePath:        "",
			},
		},
		HTTPClient: nil,
	}
}

func getServicesUsedJsonString() string {
	return `[
		{
			"Name": "AbsMemberComposite",
            "Url": "https://test-int-api-gw.centene.com/absmemberinquiry",
			"AuthEnvironmentVariable": "ABS_MBR_KEY",
			"EndPoints": {
				"MemberSearch": {
					"Name": "MemberSearch",
                    "Path": "/inq/abs/V1"
				}
			},
			"ComponentConfigOverrides":{
				"ServiceLogging":{
					"LogCallDuration":0
				},
				"Client":{
					"Timeout":0,
					"IdleConnTimeout":0,
					"MaxIdleConnsPerHost":0,
					"MaxConnsPerHost":0,
					"MaxRetries":0,
					"DisableCompression":0,
					"InsecureSkipVerify":0,
					"CABundlePath":""
				}
			}
        }
	]`
}
