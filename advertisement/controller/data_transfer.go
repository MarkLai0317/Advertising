package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MarkLai0317/Advertising/advertisement"
)

type DataTransfer struct{}

// explicitly define struct with json field tag
// because domain object (domain layer) define in advertisement.go should not know format used in API call
// the domain layer should focus on business logic
type ConditionsJSON struct {
	AgeStart  int                          `json:"ageStart"`
	AgeEnd    int                          `json:"ageEnd"`
	Genders   []advertisement.GenderType   `json:"gender"`
	Countries []advertisement.CountryCode  `json:"country"`
	Platforms []advertisement.PlatformType `json:"platform"`
}

type AdvertisementJSON struct {
	Title      string         `json:"title"`
	StartAt    time.Time      `json:"startAt"`
	EndAt      time.Time      `json:"endAt"`
	Conditions ConditionsJSON `json:"conditions"`
}

// functiono that return default value setting of advertisement
func defaultAdvertisementJSON() *AdvertisementJSON {
	return &AdvertisementJSON{
		Conditions: ConditionsJSON{
			AgeStart:  1,
			AgeEnd:    100,
			Genders:   []advertisement.GenderType{},
			Countries: []advertisement.CountryCode{},
			Platforms: []advertisement.PlatformType{},
		},
	}
}

// custom unmaarshalJSON that set default value for empty json field
func (adJson *AdvertisementJSON) UnmarshalJSON(text []byte) error {
	type AdAlias AdvertisementJSON
	ad := AdAlias(*defaultAdvertisementJSON())
	if err := json.Unmarshal(text, &ad); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return fmt.Errorf("JSON unmarshal fail")
	}

	*adJson = AdvertisementJSON(ad)
	return nil
}

// transfer JSON input to domain object -> advertisement
func (dt *DataTransfer) JSONToAdvertisement(req *http.Request) (*advertisement.Advertisement, error) {
	var adJson AdvertisementJSON
	//err := json.Unmarshal([]byte(jsonData), &adJson)
	err := json.NewDecoder(req.Body).Decode(&adJson)
	if err != nil {

		log.Printf("Error parsing JSON: %v", err)
		return nil, err
	}

	ad := &advertisement.Advertisement{
		Title:   adJson.Title,
		StartAt: adJson.StartAt,
		EndAt:   adJson.EndAt,
		Conditions: advertisement.Conditions{
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

// ===== implement method: turn response of slice of Advertisements to JSON

// Simplified structure for the final JSON
type SimplifiedAdvertisement struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}

// Wrapper for the slice of SimplifiedAdvertisement
type AdvertisementResponse struct {
	Items []SimplifiedAdvertisement `json:"items"`
}

func (dt *DataTransfer) AdvertisementSliceToJSON(ads []advertisement.Advertisement) (*AdvertisementResponse, error) {

	simplifiedAds := make([]SimplifiedAdvertisement, len(ads))
	for i, ad := range ads {
		simplifiedAds[i] = SimplifiedAdvertisement{
			Title: ad.Title,
			EndAt: ad.EndAt,
		}
	}

	// Wrap the transformed slice in the AdvertisementResponse struct
	adResponse := AdvertisementResponse{
		Items: simplifiedAds,
	}

	// Convert the response struct to JSON
	return &adResponse, nil

}

func (dt *DataTransfer) QueryToClient(req *http.Request) (client *advertisement.Client, offset int, limit int, err error) {

	// init return value
	client = nil
	offset = 0
	limit = 0
	err = nil

	queryValues := req.URL.Query()

	ageStr := queryValues.Get("age")
	genderStr := queryValues.Get("gender")
	countryStr := queryValues.Get("country")
	platformStr := queryValues.Get("platform")
	offsetStr := queryValues.Get("offset")
	limitStr := queryValues.Get("limit")

	//check missing param
	if ageStr == "" || genderStr == "" || countryStr == "" || platformStr == "" || offsetStr == "" || limitStr == "" {
		err = fmt.Errorf("missing parameter. Should coontain age, gender, country, platform, offset and limit")
		return client, offset, limit, err
	}

	// // client properties checks, the query string need to match enum defined in advertisement packaage
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		err = fmt.Errorf("age format incorrect")
		return client, offset, limit, err
	}

	gender := advertisement.GenderType(genderStr)

	country := advertisement.CountryCode(countryStr)

	platform := advertisement.PlatformType(platformStr)

	// // construct  client prooperty, offset and limit to be used as param of service call
	offset, err = strconv.Atoi(offsetStr)
	if err != nil {
		err = fmt.Errorf("offset format incorrect")
		return client, offset, limit, err
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		err = fmt.Errorf("limit format incorrect")
		return client, offset, limit, err
	}

	client = &advertisement.Client{
		Age:      age,
		Gender:   gender,
		Country:  country,
		Platform: platform,
	}

	return client, offset, limit, err
}
