package advertisement

import (
	"fmt"
	"time"
)

// abstraction of database operation that can be imlemented by posgresSQL, MYSQL, MongoDB ...
type Repository interface {
	CreateAdvertisement(ad *Advertisement) error
	GetAdvertisements(client *Client, offset int, limit int) ([]Advertisement, error)
}

type Service struct {
	repo Repository
}

// injection of Repository for Inversion of Control
func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
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

type Validator interface {
	ValidateGenderSlice(slice []GenderType) error
	ValidateCountrySlice(slice []CountryCode) error
	ValidatePlatformSlice(slice []PlatformType) error
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

func ValidateAdvertisement(ad *Advertisement, validator Validator) error {

	// check Title cannot be empty
	if ad.Title == "" {
		return fmt.Errorf("Title cannot be empty")
	}

	// check EndAt need to be after time.Now,  set StartAt to now if it's before now
	now := time.Now()

	if ad.StartAt.Before(now) {
		ad.StartAt = now
	}
	if !ad.EndAt.After(now) {
		return fmt.Errorf("EndAt cannot be smaller than current Time %s", now.Format("2006-01-02T15:04:05.000Z"))
	}

	if ad.StartAt.After(ad.EndAt) {
		return fmt.Errorf("EndAt cannot be smaller than StartAt")
	}

	// check age in range 1 to 100 and AgeEnd cannot be smaller than AgeStart
	if ad.Conditions.AgeStart < 1 || ad.Conditions.AgeStart > 100 || ad.Conditions.AgeEnd < 1 || ad.Conditions.AgeEnd > 100 {
		return fmt.Errorf("Age should be in range 1 to 100")
	}
	if ad.Conditions.AgeStart > ad.Conditions.AgeEnd {
		return fmt.Errorf("AgeStart need to be less than AgeEnd")
	}

	// check gender
	if len(ad.Conditions.Genders) == 0 {
		ad.Conditions.Genders = AllGenders()
	} else if err := validator.ValidateGenderSlice(ad.Conditions.Genders); err != nil {
		return fmt.Errorf("invalid advertisement : %w", err)
	}

	// check Country
	if len(ad.Conditions.Countries) == 0 {
		ad.Conditions.Countries = AllCountries()
	} else if err := validator.ValidateCountrySlice(ad.Conditions.Countries); err != nil {
		return fmt.Errorf("invalid advertisement : %w", err)
	}

	// check platform
	if len(ad.Conditions.Platforms) == 0 {
		ad.Conditions.Platforms = AllPlatforms()
	} else if err := validator.ValidatePlatformSlice(ad.Conditions.Platforms); err != nil {
		return fmt.Errorf("invalid advertisement : %w", err)
	}

	return nil

}

// validate client
func ValidateClient(client *Client) error {
	if client.Age < 1 || client.Age > 100 {
		return fmt.Errorf("age should be in range 1 - 100")
	}

	if !ValidGender(client.Gender) {
		return fmt.Errorf("invalid Gender type")
	}

	if !ValidCountry(client.Country) {
		return fmt.Errorf("invalid Country type")
	}

	if !ValidPlatform(client.Platform) {
		return fmt.Errorf("invalid platform type")
	}

	return nil

}

// implement domain use case -> create advertisement
func (s *Service) Create(ad *Advertisement) error {

	validator := &SliceValidator{}

	if err := ValidateAdvertisement(ad, validator); err != nil {
		return fmt.Errorf("invalid Advertisement: %w", err)
	}
	if err := s.repo.CreateAdvertisement(ad); err != nil {
		return fmt.Errorf("createAdvertisement: %w", err)
	}
	return nil

}

// implement domain use case -> return advertisement with specified conodition
func (s *Service) Advertise(client *Client, offset int, limit int) ([]Advertisement, error) {

	ValidateClient(client)
	adSlice, err := s.repo.GetAdvertisements(client, offset, limit)
	return adSlice, err
}
