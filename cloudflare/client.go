package cloudflare

import (
	//"bufio"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/franela/goreq"
)

/**
	View details of API calls here: https://support.cloudflare.com/hc/en-us/articles/216672448-Enterprise-Log-Share-REST-API
**/

const (
	API_BASE = "https://api.cloudflare.com"
)

type CloudflareClient struct {
	ApiKey         string
	Email          string
	UserServiceKey string
	RequestLogFile *RequestLogFile
	LogfileName    string
	uri            string
	StatePath      string
	debug          bool
}

// NewClient returns a new instance of a CloudflareClient struct
func NewClient(params map[string]interface{}) *CloudflareClient {

	c := &CloudflareClient{
		uri: "/client/v4/zones/%s/logs/requests",
	}

	if _, ok := params["api_key"]; ok {
		c.ApiKey = params["api_key"].(string)
		c.Email = params["email"].(string)
	} else {
		c.UserServiceKey = params["user_service_key"].(string)
	}

	if _, ok := params["debug"]; ok {
		c.debug = params["debug"].(bool)
	}

	if _, ok := params["state_path"]; ok {
		c.StatePath = params["state_path"].(string)
	}
	return c
}

func (c *CloudflareClient) doRequest(params map[string]interface{}) (string, error) {

	qsa := url.Values{}
	apiURL := API_BASE + fmt.Sprintf(c.uri, params["zone_tag"].(string))

	if _, ok := params["time_start"]; ok {
		qsa.Set("start", fmt.Sprintf("%d", params["time_start"].(int)))
	}
	if _, ok := params["time_end"]; ok {
		qsa.Set("end", fmt.Sprintf("%d", params["time_end"].(int)))
	}
	if _, ok := params["count"]; ok {
		qsa.Set("count", fmt.Sprintf("%d", params["count"].(int)))
	}

	// goreq.DefaultTransport = &http.Transport{Dial: goreq.DefaultDialer.Dial,
	// 	Proxy:              http.ProxyFromEnvironment,
	// 	DisableCompression: true}

	req := goreq.Request{
		Uri:         apiURL,
		Timeout:     10 * time.Minute,
		ShowDebug:   c.debug, // true
		QueryString: qsa,
	}

	req.AddHeader("Accept-encoding", "gzip")
	if c.UserServiceKey != "" {
		req.AddHeader("X-User-Service-Key", c.UserServiceKey)
	} else {
		req.AddHeader("X-Auth-Key", c.ApiKey)
		req.AddHeader("X-Auth-Email", c.Email)
	}

	logp.Debug("http", "%d Downloading log file...", params["time_start"])

	for i := 0; i < 4; i++ {
		response, err := req.Do()
		if err != nil {
			logp.Debug("http", "%d Error request  %s", err)
			return "", err
		}

		logp.Debug("http", "%d Try %d Request code %d, %s %s", params["time_start"], i, response.StatusCode, response.Status, qsa)
		logp.Debug("http", "%d Try %d Request Content type %s", params["time_start"], i, response.Header.Get("Content-Type"))

		if (response.StatusCode == 200) && (response.Header.Get("Content-Type") == "application/json") {
			// Now need to save all the resposne content to a file
			logFileName := fmt.Sprintf("cloudflare_logs_%d_to_%d.txt.gz", params["time_start"].(int), params["time_end"].(int))
			rlf := NewRequestLogFile(logFileName, c.StatePath)

			//	logp.Debug("http", "Body"+nBody)
			nBytes, err := rlf.SaveFromHttpResponseBody(response.Body)
			if err != nil {
				logp.Debug("http", "%d Error saving %s", params["time_start"], err)
				return "", err
			}
			logp.Debug("http", "%d Downloaded %d bytes to %s %s", params["time_start"], nBytes, logFileName, c.StatePath)
			if nBytes <= 23 {
				logp.Debug("http", "%d Error body %s", params["time_start"], response.Status)
				return "", errors.New("Request body is empty")
			}
			return logFileName, nil
		} else {
			logp.Debug("http", "%d Wrong status code : %s", params["time_start"], response.Status)
		}
		// sleep a second
		logp.Debug("http", "%d Sleeping 2 seconds", params["time_start"])
		time.Sleep(2 * time.Second)
	}
	logp.Debug("http", "%d Returning empty", params["time_start"])
	return "", nil
}

func (c *CloudflareClient) GetLogRangeFromTimestamp(opts map[string]interface{}) (string, error) {
	filename, err := c.doRequest(opts)
	if err != nil {
		logp.Debug("http", "Error doing request %s", err)
		return "", err
	}
	return filename, nil
}
