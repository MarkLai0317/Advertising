package ad

import (
	"fmt"
	"time"
)

// abstraction of database operation that can be imlemented by posgresSQL, MYSQL, MongoDB ...
type Repository interface {
	CreateAdvertisement(advertisement *Advertisement) error
	GetAdvertisements(client *Client, now time.Time) ([]Advertisement, error)
}

type Service struct {
	repo Repository
	// GetRepo    Repository
}

// injection of Repository for Inversion of Control
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		//GetRepo:    getRepository,
	}
}

// service need to implement {
// 	   Create(ad *Advertisement) error
// 	   Advertise(client *Client, offset int, limit int) ([]Advertisement, error)
// }

// implement domain use case -> create advertisement
func (s *Service) CreateAd(advertisement *Advertisement) error {

	sliceValidator := &SliceValidator{}

	if err := ValidateAdvertisement(advertisement, sliceValidator); err != nil {
		return fmt.Errorf("invalid Advertisement: %w", err)
	}
	if err := s.repo.CreateAdvertisement(advertisement); err != nil {
		return fmt.Errorf("error creating Advertisement in DB: %w", err)
	}
	return nil

}

// helper function for CreateAd
type Validator interface {
	ValidateGenderSlice(slice []GenderType) error
	ValidateCountrySlice(slice []CountryCode) error
	ValidatePlatformSlice(slice []PlatformType) error
}

// generic function to check the invalid slice enum string
func ValidateSlice[T comparable](slice []T, valid func(T) bool) error {
	for _, item := range slice {
		if !valid(item) {
			return fmt.Errorf("invalid item in slice of type %T: %v", item, item)
		}
	}
	return nil
}

type SliceValidator struct{}

func (sv SliceValidator) ValidateGenderSlice(slice []GenderType) error {
	return ValidateSlice(slice, ValidGender)
}

func (sv SliceValidator) ValidateCountrySlice(slice []CountryCode) error {
	return ValidateSlice(slice, ValidCountry)
}

func (sv SliceValidator) ValidatePlatformSlice(slice []PlatformType) error {
	return ValidateSlice(slice, ValidPlatform)
}

func ValidateAdvertisement(advertisement *Advertisement, validator Validator) error {

	// check Title cannot be empty
	if advertisement.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// check EndAt need to be after time.Now,  set StartAt to now if it's before now
	now := time.Now()

	if advertisement.StartAt.Before(now) {
		advertisement.StartAt = now
	}
	if !advertisement.EndAt.After(now) {
		return fmt.Errorf("endAt cannot be smaller than current Time %s", now.Format("2006-01-02T15:04:05.000Z"))
	}

	if advertisement.StartAt.After(advertisement.EndAt) {
		return fmt.Errorf("endAt cannot be smaller than StartAt")
	}

	// check age in range 1 to 100 and AgeEnd cannot be smaller than AgeStart
	if advertisement.Conditions.AgeStart < 1 || advertisement.Conditions.AgeStart > 100 || advertisement.Conditions.AgeEnd < 1 || advertisement.Conditions.AgeEnd > 100 {
		return fmt.Errorf("age should be in range 1 to 100")
	}
	if advertisement.Conditions.AgeStart > advertisement.Conditions.AgeEnd {
		return fmt.Errorf("ageStart need to be less than AgeEnd")
	}

	// check gender
	if len(advertisement.Conditions.Genders) == 0 {
		advertisement.Conditions.Genders = AllGenders()
	} else if err := validator.ValidateGenderSlice(advertisement.Conditions.Genders); err != nil {
		return fmt.Errorf("invalid advertisement : %w", err)
	}

	// check Country
	if len(advertisement.Conditions.Countries) == 0 {
		advertisement.Conditions.Countries = AllCountries()
	} else if err := validator.ValidateCountrySlice(advertisement.Conditions.Countries); err != nil {
		return fmt.Errorf("invalid advertisement : %w", err)
	}

	// check platform
	if len(advertisement.Conditions.Platforms) == 0 {
		advertisement.Conditions.Platforms = AllPlatforms()
	} else if err := validator.ValidatePlatformSlice(advertisement.Conditions.Platforms); err != nil {
		return fmt.Errorf("invalid advertisement : %w", err)
	}

	return nil

}

// implement domain use case -> return advertisement with specified conodition
func (s *Service) Advertise(client *Client) ([]Advertisement, error) {

	err := ValidateClient(client)
	if err != nil {
		return nil, fmt.Errorf("invalid client: %w", err)
	}
	// if offset < 0 {
	// 	return nil, fmt.Errorf("offset cannot less than 0")
	// }

	// if limit < 1 {
	// 	return nil, fmt.Errorf("limit cannot less than 1")
	// }
	adSlice, err := s.repo.GetAdvertisements(client, time.Now())

	if err != nil {
		return nil, fmt.Errorf("GetAdvertisements Error: %w", err)
	}

	return adSlice, err
}

// validate client
func ValidateClient(client *Client) error {
	if !client.AgeMissing && (client.Age < 1 || client.Age > 100) {
		return fmt.Errorf("age should be in range 1 - 100")
	}

	if !client.GenderMissing && !ValidGender(client.Gender) {
		return fmt.Errorf("invalid Gender type")
	}

	if !client.CountryMissing && !ValidCountry(client.Country) {
		return fmt.Errorf("invalid Country type")
	}

	if !client.PlatformMissing && !ValidPlatform(client.Platform) {
		return fmt.Errorf("invalid Platform type")
	}

	if client.Offset < 0 {
		return fmt.Errorf("invalid offset")
	}

	if client.Limit < 1 {
		return fmt.Errorf("invalid limit")
	}

	return nil

}

// high level public function defined first
// if interface as param, define interface first
// helper and private function called by high level function define below ordered by called order

// do the same for each public function
