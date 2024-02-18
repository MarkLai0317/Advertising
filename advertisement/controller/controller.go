package controller

import (
	"encoding/json"
	"net/http"

	"github.com/MarkLai0317/Advertising/advertisement"
)

type ServiceError struct {
	Message string `json:"message"`
}

type DataTransferer interface {
	JSONToAdvertisement(req *http.Request) (*advertisement.Advertisement, error)
	AdvertisementSliceToJSON(ads []advertisement.Advertisement) (*AdvertisementResponse, error)
	QueryToClient(req *http.Request) (client *advertisement.Client, offset int, limit int, err error)
}

type Controller struct {
	adService      advertisement.UseCase
	dataTransferer DataTransferer
}

func NewAdvertisementController(service advertisement.UseCase, dataDataTransferer DataTransferer) *Controller {

	return &Controller{
		adService:      service,
		dataTransferer: dataDataTransferer,
	}
}

func (c *Controller) CreateAdvertisement(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")
	newAd, err := c.dataTransferer.JSONToAdvertisement(req)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: "Error decode req.Body"})
		return
	}
	err = c.adService.Create(newAd)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(ServiceError{Message: "Error creating advertisement"})
		return
	}
	resp.WriteHeader(http.StatusOK)

	return
}

func (c *Controller) advertise(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")
	client, offset, limit, err := c.dataTransferer.QueryToClient(req)
	if err != nil {
		// Use http.Error to send the error message back to the client
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	adSlice, err := c.adService.Advertise(client, offset, limit)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	adResponse, err := c.dataTransferer.AdvertisementSliceToJSON(adSlice)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(adResponse)

}
