package szkopul_client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"sio-tool/cookiejar"

	"github.com/fatih/color"
)

// SzkopulClient szkopul client
type SzkopulClient struct {
	Jar            *cookiejar.Jar `json:"cookies"`
	Username       string         `json:"handle"`
	Token          string         `json:"token"`
	LastSubmission *Info          `json:"last_submission"`
	host           string
	path           string
	client         *http.Client
}

// Instance global client
var Instance *SzkopulClient

// Init initialize
func Init(path, host, proxy string) {
	jar, _ := cookiejar.New(nil)
	c := &SzkopulClient{Jar: jar, LastSubmission: nil, path: path, host: host, client: nil}
	if err := c.load(); err != nil {
		color.Red(err.Error())
		color.Green("Create a new session in %v", path)
	}
	Proxy := http.ProxyFromEnvironment
	if len(proxy) > 0 {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			color.Red(err.Error())
			color.Green("Use default proxy from environment")
		} else {
			Proxy = http.ProxyURL(proxyURL)
		}
	}
	c.client = &http.Client{Jar: c.Jar, Transport: &http.Transport{Proxy: Proxy}}
	if err := c.save(); err != nil {
		color.Red(err.Error())
	}
	Instance = c
}

// load from path
func (c *SzkopulClient) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, c)
}

// save file to path
func (c *SzkopulClient) save() (err error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		err = os.WriteFile(c.path, data, 0644)
	}
	if err != nil {
		color.Red("Cannot save session to %v\n%v", c.path, err.Error())
	}
	return
}