package twiliogo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

const ROOT = "https://api.twilio.com"
const VERSION = "2010-04-01"

var HttpClientFactory func(context.Context) *http.Client

func init() {
	HttpClientFactory = func(c context.Context) *http.Client {
		transprt := http.Transport{}
		return &http.Client{
			Transport: &transprt,
		}
	}
}

type Client interface {
	AccountSid() string
	AuthToken() string
	RootUrl() string
	get(url.Values, string) ([]byte, error)
	post(url.Values, string) ([]byte, error)
}

type TwilioClient struct {
	accountSid string
	authToken  string
	rootUrl    string
	context    context.Context
}

func NewClient(accountSid, authToken string, context context.Context) *TwilioClient {
	rootUrl := "/" + VERSION + "/Accounts/" + accountSid
	return &TwilioClient{accountSid, authToken, rootUrl, context}
}

func (client *TwilioClient) post(values url.Values, uri string) ([]byte, error) {
	req, err := http.NewRequest("POST", ROOT+uri, strings.NewReader(values.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.AccountSid(), client.AuthToken())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient := HttpClientFactory(client.context)

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return body, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		if res.StatusCode == 500 {
			return body, Error{"Server Error"}
		} else {
			twilioError := new(TwilioError)
			json.Unmarshal(body, twilioError)
			return body, twilioError
		}
	}

	return body, err
}

func (client *TwilioClient) get(queryParams url.Values, uri string) ([]byte, error) {
	var params *strings.Reader

	if queryParams == nil {
		queryParams = url.Values{}
	}

	params = strings.NewReader(queryParams.Encode())
	req, err := http.NewRequest("GET", ROOT+uri, params)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.AccountSid(), client.AuthToken())
	httpClient := HttpClientFactory(client.context)

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return body, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		if res.StatusCode == 500 {
			return body, Error{"Server Error"}
		} else {
			twilioError := new(TwilioError)
			json.Unmarshal(body, twilioError)
			return body, twilioError
		}
	}

	return body, err
}

func (client *TwilioClient) AccountSid() string {
	return client.accountSid
}

func (client *TwilioClient) AuthToken() string {
	return client.authToken
}

func (client *TwilioClient) RootUrl() string {
	return client.rootUrl
}
