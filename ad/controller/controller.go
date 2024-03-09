package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MarkLai0317/Advertising/ad"
)

type ServiceError struct {
	Message string `json:"message"`
}

type DataTransferer interface {
	JSONToAdvertisement(req *http.Request) (*ad.Advertisement, error)
	AdvertisementSliceToJSON(ads []ad.Advertisement) ([]byte, error)
	QueryToClient(req *http.Request) (client *ad.Client, err error)
}

type Controller struct {
	AdService      ad.UseCase
	DataTransferer DataTransferer
}

func NewAdvertisementController(service ad.UseCase, dataDataTransferer DataTransferer) *Controller {
	return &Controller{
		AdService:      service,
		DataTransferer: dataDataTransferer,
	}
}

func (c *Controller) CreateAdvertisement(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")
	newAd, err := c.DataTransferer.JSONToAdvertisement(req)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: fmt.Sprintf("Error JSON to advertisement: %v", err.Error())})
		return
	}

	err = c.AdService.CreateAd(newAd)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: fmt.Sprintf("Error creating advertisement: %v", err.Error())})
		return
	}

	resp.WriteHeader(http.StatusOK)

}

func (c *Controller) Advertise(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")
	client, err := c.DataTransferer.QueryToClient(req)
	if err != nil {
		// Use http.Error to send the error message back to the client
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: fmt.Sprintf("Error decode url query: %v", err.Error())})
		return
	}

	adSlice, err := c.AdService.Advertise(client)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: fmt.Sprintf("Error getting ad: %v", err.Error())})
		return
	}

	adResponse, err := c.DataTransferer.AdvertisementSliceToJSON(adSlice)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: fmt.Sprintf("Error ad slice to JSON: %v", err.Error())})
		return
	}

	_, err = resp.Write(adResponse)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: fmt.Sprintf("Error writing to body %v", err.Error())})
		return
	}
	resp.WriteHeader(http.StatusOK)

}
