package eero

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-resty/resty/v2"
)

type Eero interface {
	Login(identifier string) (LoginData, error)
	LoginVerify(userToken string, verificationCode string) error
	Networks() ([]NetworkData, error)
	Account() (AccountsData, error)
	Devices() ([]DeviceData, error)
}

type eero struct {
	userToken string
	client    *resty.Client
}

func New(userToken string) Eero {
	client := resty.New()
	client.SetBaseURL("https://api-user.e2ro.com/2.2/")

	if userToken != "" {
		client.SetCookie(&http.Cookie{
			Name:  "s",
			Value: userToken,
		})
	}

	return &eero{
		userToken: userToken,
		client:    client,
	}
}

func (e *eero) Login(identifier string) (LoginData, error) {
	data := &LoginResponse{}

	resp, err := e.client.R().
		SetBody(LoginRequest{Identifier: identifier}).
		SetResult(data).
		Post("login")

	// If we error out, return an empty response
	if err != nil {
		return LoginData{}, err
	}

	// If we error out, return an empty response
	if resp.StatusCode() != 200 {
		return LoginData{}, &data.MetaResponse
	}

	return data.Data, nil
}

func (e *eero) LoginVerify(sessionToken string, verificationCode string) error {
	data := &AccountsResponse{}
	resp, err := e.client.R().
		SetBody(LoginVerifyRequest{Code: verificationCode}).
		SetCookie(&http.Cookie{
			Name:  "s",
			Value: sessionToken,
		}).
		SetResult(data).
		SetError(data).
		Post("login/verify")

	// If we error out, return an empty response
	if err != nil {
		return err
	}

	// If we error out, return an empty response
	if resp.StatusCode() != 200 {
		return &data.MetaResponse
	}

	return nil
}

func (e *eero) Account() (AccountsData, error) {
	data := &AccountsResponse{}
	resp, err := e.client.R().
		SetResult(data).
		SetError(data).
		Get("account")

	// If we error out, return an empty response
	if err != nil {
		return AccountsData{}, err
	}
	// If we error out, return an empty response
	if resp.StatusCode() != 200 {
		return AccountsData{}, &data.MetaResponse
	}

	return data.Data, nil
}

func (e *eero) Networks() ([]NetworkData, error) {
	account, err := e.Account()
	if err != nil {
		return nil, err
	}

	if account.Networks.Count == 0 {
		return []NetworkData{}, nil
	}

	return account.Networks.Data, nil
}

func (e *eero) Devices() ([]DeviceData, error) {
	networks, err := e.Networks()

	if err != nil {
		return []DeviceData{}, err
	}

	// Create a new slice DeviceData structs
	devices := []DeviceData{}

	re := regexp.MustCompile(`\d+$`)
	for _, network := range networks {
		// Extract the id number from the URL
		id := re.FindString(network.URL)

		// Error if the match comes back empty
		if id == "" {
			log.Fatalf("Could not find ID in network string '%s'", network.URL)
		}

		temp := &DevicesResponse{}
		resp, err := e.client.R().
			SetResult(temp).
			SetError(temp).
			Get(fmt.Sprintf("networks/%s/devices", id))

		// If we error out, return an empty response
		if err != nil {
			return []DeviceData{}, err
		}

		// If we error out, return an empty response
		if resp.StatusCode() != 200 {
			return []DeviceData{}, &temp.MetaResponse
		}

		// Otherwise add the output into the slice
		devices = append(devices, temp.Data...)
	}

	return devices, nil
}
