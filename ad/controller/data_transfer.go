package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MarkLai0317/Advertising/ad"
)

type AdDataTransferer struct{}

func NewAdDataTransferer() *AdDataTransferer {
	return &AdDataTransferer{}
}

// explicitly define struct with json field tag
// because domain object (domain layer) define in ad.go should not know format used in API call
// the domain layer should focus on business logic
type ConditionsJSON struct {
	AgeStart  int               `json:"ageStart"`
	AgeEnd    int               `json:"ageEnd"`
	Genders   []ad.GenderType   `json:"gender"`
	Countries []ad.CountryCode  `json:"country"`
	Platforms []ad.PlatformType `json:"platform"`
}

type AdvertisementJSON struct {
	Title      string         `json:"title"`
	StartAt    time.Time      `json:"startAt"`
	EndAt      time.Time      `json:"endAt"`
	Conditions ConditionsJSON `json:"conditions"`
}

// transfer JSON input to domain object -> advertisement
func (adt *AdDataTransferer) JSONToAdvertisement(req *http.Request) (*ad.Advertisement, error) {
	var adJson AdvertisementJSON
	//err := json.Unmarshal([]byte(jsonData), &adJson)
	err := json.NewDecoder(req.Body).Decode(&adJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	ad := &ad.Advertisement{
		Title:   adJson.Title,
		StartAt: adJson.StartAt,
		EndAt:   adJson.EndAt,
		Conditions: ad.Conditions{
			AgeStart:  adJson.Conditions.AgeStart,
			AgeEnd:    adJson.Conditions.AgeEnd,
			Genders:   adJson.Conditions.Genders,
			Countries: adJson.Conditions.Countries,
			Platforms: adJson.Conditions.Platforms,
		},
	}
	log.Printf("Parsed Advertisement: %+v\n", ad)
	return ad, nil

}

// custom unmaarshalJSON that set default value for empty json field
func (adJson *AdvertisementJSON) UnmarshalJSON(text []byte) error {
	type AdAlias AdvertisementJSON
	ad := AdAlias(*DefaultAdvertisementJSON())
	if err := json.Unmarshal(text, &ad); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return fmt.Errorf("JSON unmarshal fail: %w", err)
	}

	*adJson = AdvertisementJSON(ad)
	return nil
}

// function that return default value setting of advertisement
func DefaultAdvertisementJSON() *AdvertisementJSON {
	return &AdvertisementJSON{
		Conditions: ConditionsJSON{
			AgeStart:  1,
			AgeEnd:    100,
			Genders:   []ad.GenderType{},
			Countries: []ad.CountryCode{},
			Platforms: []ad.PlatformType{},
		},
	}
}

// turn response of slice of Advertisements to JSON

func (adt *AdDataTransferer) QueryToClient(req *http.Request) (client *ad.Client, err error) {

	// default offset and limit
	client = &ad.Client{}
	client.Offset = 0
	client.Limit = 5
	err = nil

	// get query value
	queryValues := req.URL.Query()

	ageStr := queryValues.Get("age")
	genderStr := queryValues.Get("gender")
	countryStr := queryValues.Get("country")
	platformStr := queryValues.Get("platform")
	offsetStr := queryValues.Get("offset")
	limitStr := queryValues.Get("limit")

	//check missing param
	client.AgeMissing = ageStr == ""
	client.GenderMissing = genderStr == ""
	client.CountryMissing = countryStr == ""
	client.PlatformMissing = platformStr == ""

	// age string to int
	if !client.AgeMissing {
		client.Age, err = strconv.Atoi(ageStr)
		if err != nil {
			err = fmt.Errorf("age format incorrect")
			return nil, err
		}
	}

	// convert string to enum type
	client.Gender = ad.GenderType(genderStr)
	client.Country = ad.CountryCode(countryStr)
	client.Platform = ad.PlatformType(platformStr)

	// offset string to int
	if offsetStr != "" {
		tempOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			err = fmt.Errorf("offset format incorrect")
			return nil, err
		}
		client.Offset = tempOffset
	}
	// limit string to int
	if limitStr != "" {
		tempLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			err = fmt.Errorf("limit format incorrect")
			return nil, err
		}
		client.Limit = tempLimit
	}

	return client, nil
}

// =============== advertisement slice to json data transferer ===========

// Simplified structure for the final JSON
type SimplifiedAdvertisement struct {
	Title string       `json:"title"`
	EndAt ResponseTime `json:"endAt"`
}

type ResponseTime time.Time

// Wrapper for the slice of SimplifiedAdvertisement
type AdvertisementResponse struct {
	Items []SimplifiedAdvertisement `json:"items"`
}

func (adt *AdDataTransferer) AdvertisementSliceToJSON(ads []ad.Advertisement) ([]byte, error) {

	simplifiedAds := make([]SimplifiedAdvertisement, len(ads))
	for i, ad := range ads {
		simplifiedAds[i] = SimplifiedAdvertisement{
			Title: ad.Title,
			EndAt: ResponseTime(ad.EndAt),
		}
	}

	// Wrap the transformed slice in the AdvertisementResponse struct
	adResponse := AdvertisementResponse{
		Items: simplifiedAds,
	}

	jsonDataBytes, err := json.Marshal(adResponse)
	if err != nil {
		return nil, fmt.Errorf("eror marshal adResponse")
	}
	// Convert the response struct to JSON
	return jsonDataBytes, err
}

// custom time format for json.NewEncoder(resp).Encode(adResponse)
func (rt ResponseTime) MarshalJSON() ([]byte, error) {
	t := time.Time(rt)
	formattedTime := t.Format("2006-01-02T15:04:05.000Z") // Adjust the layout according to your requirement
	return []byte(fmt.Sprintf(`"%s"`, formattedTime)), nil
}
