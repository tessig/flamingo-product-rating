package infrastructure

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/tessig/flamingo-product-rating/src/app/infrastructure"
)

type (
	// Client represents an external data source which provides an REST endpoint delivering JSON bodies to consume
	Client struct {
		baseURL        string
		listEndpoint   string
		detailEndpoint string
		logger         flamingo.Logger
	}
)

var _ infrastructure.Source = new(Client)

// Inject dependencies
func (c *Client) Inject(
	logger flamingo.Logger,
	conf *struct {
		BaseURL        string `inject:"config:productservice.baseurl"`
		ListEndpoint   string `inject:"config:productservice.endpoints.list"`
		DetailEndpoint string `inject:"config:productservice.endpoints.detail"`
	},
) {
	c.logger = logger
	c.baseURL = strings.TrimRight(conf.BaseURL, "/") + "/"
	c.listEndpoint = conf.ListEndpoint
	c.detailEndpoint = conf.DetailEndpoint
}

// Detail returns the raw data for a given  product ID
func (c *Client) Detail(ctx context.Context, pid int) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "products/client/Detail")
	defer span.End()

	return c.Get(ctx, c.baseURL+c.detailEndpoint, nil, []string{":pid", strconv.Itoa(pid)})
}

// All returns the raw data for all products
func (c *Client) All(ctx context.Context) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "products/client/All")
	defer span.End()

	return c.Get(ctx, c.baseURL+c.listEndpoint, nil, nil)
}

// Get does a GET-call with the given parameters
func (c *Client) Get(ctx context.Context, url string, params map[string]string, urlParams []string) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "products/client/Get")
	defer span.End()

	replacer := strings.NewReplacer(urlParams...)
	requestURL := replacer.Replace(url)

	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error on creating request")
	}
	request.WithContext(ctx)

	query := request.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	request.URL.RawQuery = query.Encode()

	start := time.Now()
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Error on executing request")
	}
	duration := time.Since(start) / time.Millisecond
	c.logger.WithFields(map[flamingo.LogKey]interface{}{
		flamingo.LogKeyApicall:           1,
		flamingo.LogKeyRequest:           "",
		flamingo.LogKeyRequestedURL:      request.URL,
		flamingo.LogKeyRequestedEndpoint: requestURL,
		flamingo.LogKeyMethod:            request.Method,
		flamingo.LogKeyRequestTime:       fmt.Sprintf("%d", duration),
		flamingo.LogKeyResponseCode:      resp.StatusCode,
	}).Info("collected data from external source")
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	if statusCode >= 400 {
		return nil, errors.Errorf("Call for %s failed with status code %d", request.URL, statusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
