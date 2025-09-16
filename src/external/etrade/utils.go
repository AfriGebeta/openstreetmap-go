package etrade

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"org-fayda/src/config"
	"org-fayda/src/utils"
	"time"
)

func GetETradeRegistrationInfo(tin string) (*BusinessRegistration, error) {
	var resp, err = MakeETradeRequest(ClientForValue.RegistrationInfo, tin, nil)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return &BusinessRegistration{}, err
	}

	var businessRegistration BusinessRegistration

	err = json.Unmarshal(body, &businessRegistration)

	if err != nil {
		return &BusinessRegistration{}, err
	}

	return &businessRegistration, nil
}

func GetETradeBusinessLicenseNumber(tin string, businessLicenseNumber string) (*BusinessInfo, error) {
	var resp, err = MakeETradeRequest(ClientForValue.BusinessLicenseNumber, tin, &businessLicenseNumber)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return &BusinessInfo{}, err
	}

	var businessInfo BusinessInfo

	err = json.Unmarshal(body, &businessInfo)

	if err != nil {
		return &BusinessInfo{}, err
	}

	return &businessInfo, nil
}

func MakeETradeRequest(_for ClientFor, tin string, businessLicenseNumber *string) (*http.Response, error) {
	req, err := GetETradeRequest(_for, tin, businessLicenseNumber)

	if err != nil {
		return &http.Response{}, err
	}

	return (&http.Client{Timeout: 30 * time.Second}).Do(req)
}

func GetETradeRequest(_for ClientFor, tin string, businessLicenseNumber *string) (*http.Request, error) {
	var (
		method = "GET"
		uri    = config.ServerConfigObject.ETradeBaseUrl
	)

	if tin == "" && (businessLicenseNumber == nil || *businessLicenseNumber == "") {
		return &http.Request{}, errors.New("no identifier provided")
	}

	switch _for {
	case ClientForValue.RegistrationInfo:
		if tin == "" {
			return &http.Request{}, errors.New("no tin provided")
		}

		uri = uri + "/Registration/GetRegistrationInfoByTin/" + tin + "/en"
	case ClientForValue.BusinessLicenseNumber:
		if businessLicenseNumber == nil || *businessLicenseNumber == "" {
			return &http.Request{}, errors.New("no business license number provided")
		}

		uri = uri + "/BusinessMain/GetBusinessByLicenseNo?LicenseNo=" + *businessLicenseNumber + "&Lang=en"
	}

	var tmp, err = http.NewRequest(method, uri, nil)

	if err != nil {
		utils.LogError("proxy", "error creating upstream request: "+err.Error())
		return &http.Request{}, err
	}

	tmp.Header.Set("Referer", config.ServerConfigObject.ETradeAcceptableRefererUrl)

	return tmp, nil
}
