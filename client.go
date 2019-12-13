package higgs


import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type (
	Client struct {
		HTTP         *http.Client
		Store        *DB
		Log          *log.Logger
		UserAgent    string
		ESIRateLimit *safeCounter
		RetryLimit   int
		MaxRoutines  int
	}

	safeCounter struct {
		cnt int
		mux sync.Mutex
	}
)

func newClient(config Configuration) (*Client, error) {
	logger := log.New(os.Stdout, "CLIENT:", log.Lshortfile|log.Ldate|log.Ltime)

	// now check we have access to mongo

	store, err := GetDatabaseHandle(config)

	if err != nil {
		return nil, err
	}

	rateLimESI := &safeCounter{}

	go func() {
		for {
			time.Sleep(time.Second)
			if rateLimESI.Value() > 0 {
				rateLimESI.Dec()
			}
		}
	}()

	return &Client{
		HTTP: &http.Client{
			Timeout: time.Second * time.Duration(config.Web.TimeoutSec),
			Transport: &http.Transport{
				MaxConnsPerHost:     10,
				MaxIdleConnsPerHost: 2,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		},
		Store:        store,
		Log:          logger,
		UserAgent:    config.Web.UserAgent,
		ESIRateLimit: rateLimESI,
		RetryLimit:   25,
		MaxRoutines:  config.App.MaxRoutines,
	}, nil

}

func (c *Client) makeRawHTTPGet(url string) ([]byte, int, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to buid http request")
	}

	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTP.Do(req)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to make request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to read response from request")
	}

	return body, res.StatusCode, nil
}

func (c *Client) MakeESIGet(url string) (out []byte, err error) {

	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		for c.ESIRateLimit.Value() > 10 {
			time.Sleep(500 * time.Millisecond)
		}

		body, status, err := c.makeRawHTTPGet(url)
		if err != nil {
			if strings.Contains(err.Error(), "too many open files") {
				// This is not going to hurt to keep retrying
				retriesRemain++
			} else {
				// fmt.Printf("ESI GET ERROR - %v\n", err)
			}
			continue
		}
		if !(status == 200) {
			// Increment the counter :(
			c.ESIRateLimit.Inc()
			// fmt.Printf("ESI GET RESPONSE ERROR - %v - %v - %v\n", status, url, string(body))
			time.Sleep(250 * time.Millisecond)
			continue
		}

		return body, err
	}

	return nil, fmt.Errorf("Max retries exceeded for url: ; err: %v", url, err)
}


func (c *Client) MakeGetRequestWithRetry(url string) ([]byte, error) {
	retriesRemain := c.RetryLimit
	for retriesRemain > 1 {
		retriesRemain--

		body, status, err := c.makeRawHTTPGet(url)
		if err != nil {
			continue
		}
		if !(status >= 200 && status < 300) {
			time.Sleep(250 * time.Millisecond)
			continue
		}

		return body, err
	}

	return nil, fmt.Errorf("Max retries exceeded for url: ", url)
}

// Inc increments the counter.
func (c *safeCounter) Inc() {
	c.mux.Lock()
	c.cnt++
	c.mux.Unlock()
}

// Dec decrements the counter.
func (c *safeCounter) Dec() {
	c.mux.Lock()
	if c.cnt > 0 {
		c.cnt--
	}
	c.mux.Unlock()
}

// Value returns the current value of the counter.
func (c *safeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.cnt
}
