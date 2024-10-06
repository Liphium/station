package settings

import (
	"errors"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/localization"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type setting[T any] struct {
	Label        localization.Translations // The label for the UI
	Name         string                    // Id of the setting in the database
	DefaultValue T                         // The default value

	// Cached current value (for quicker access)
	currentValueSet bool
	currentValue    T
}

// Set the value of the setting in the database (and in the cache)
func (s *setting[T]) SetValue(value T) error {

	// Encode the value to a string
	str, err := s.encode(value)
	if err != nil {
		return err
	}

	// Insert the setting into the database (or update)
	if err := database.DBConn.Save(&database.Setting{
		Name:  s.Name,
		Value: str,
	}).Error; err != nil {
		return err
	}

	// Cache the value
	s.currentValueSet = true
	s.currentValue = value

	return nil
}

func (s *setting[T]) GetValue() (T, error) {

	// Check if the value is already cached
	if s.currentValueSet {
		return s.currentValue, nil
	}

	// Get the value from the database
	var dbSetting database.Setting
	if err := database.DBConn.Where("name = ?", s.Name).Take(&dbSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			// Cache default value as current value
			s.currentValueSet = true
			s.currentValue = s.DefaultValue

			return s.DefaultValue, nil
		}

		return s.DefaultValue, err
	}

	// Decode the value from the database
	value, err := s.Decode(dbSetting.Value)
	if err != nil {
		return s.DefaultValue, err
	}

	// Cache the value
	s.currentValueSet = true
	s.currentValue = value

	return value, nil
}

// Encode a value to string using sonic (json)
func (s setting[T]) encode(value T) (string, error) {
	return sonic.MarshalString(value)
}

// Decode a string json value using sonic
func (s setting[T]) Decode(value string) (T, error) {
	var result T
	err := sonic.UnmarshalString(value, &result)
	return result, err
}

func (s setting[T]) ToMap(locale string) fiber.Map {
	return fiber.Map{
		"name":  s.Name,
		"label": localization.TranslateLocale(locale, s.Label),
		"value": s.currentValue,
	}
}

// Settings registry for integers (int64)
var SettingRegistryInteger = map[string]*setting[int64]{
	FilesMaxUploadSize.Name:   FilesMaxUploadSize,
	FilesMaxTotalStorage.Name: FilesMaxTotalStorage,
}

// Settings registry for booleans
var SettingRegistryBoolean = map[string]*setting[bool]{
	DecentralizationEnabled.Name:     DecentralizationEnabled,
	DecentralizationAllowUnsafe.Name: DecentralizationAllowUnsafe,
}
