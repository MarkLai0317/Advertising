package ad

import (
	"time"
)

// enum of Gender
type GenderType string

const (
	Male   GenderType = "M"
	Female GenderType = "F"
)

var validGender = map[GenderType]bool{
	Male:   true,
	Female: true,
}

func ValidGender(gender GenderType) bool {
	return validGender[gender]
}

// var StringToGenderType = map[string]GenderType{
// 	"M": Male,
// 	"F": Female,
// }

func AllGenders() []GenderType {
	return []GenderType{
		Male, Female,
	}
}

// enum of platform
type PlatformType string

const (
	Android PlatformType = "android"
	Ios     PlatformType = "ios"
	Web     PlatformType = "web"
)

var validPlatform = map[PlatformType]bool{
	Android: true,
	Ios:     true,
	Web:     true,
}

func ValidPlatform(platform PlatformType) bool {
	return validPlatform[platform]
}

// var StringToPlatformType = map[string]PlatformType{
// 	"android": Android,
// 	"ios":     Ios,
// 	"web":     Web,
// }

// func AllGenders() []GenderType {
// 	return []GenderType{
// 		Male, Female,
// 	}
// }

func AllPlatforms() []PlatformType {
	return []PlatformType{
		Ios, Android, Web,
	}
}

// enum of country is at country_code.go since the country enum is very long

// conditions eneity
type Conditions struct {
	AgeStart  int
	AgeEnd    int
	Genders   []GenderType
	Countries []CountryCode
	Platforms []PlatformType
}

// advertisement entity
type Advertisement struct {
	Title      string
	StartAt    time.Time
	EndAt      time.Time
	Conditions Conditions
}

type Client struct {
	Age      int
	Gender   GenderType
	Country  CountryCode
	Platform PlatformType
	Offset   int
	Limit    int

	AgeMissing      bool
	GenderMissing   bool
	CountryMissing  bool
	PlatformMissing bool
}

type UseCase interface {
	CreateAd(ad *Advertisement) error
	Advertise(client *Client) ([]Advertisement, error)
}
