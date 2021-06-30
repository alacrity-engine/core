package system

import (
	"fmt"

	"gopkg.in/go-ini/ini.v1"
)

var (
	configFile *ini.File
)

// LoadConfig loads the contents of the
// .ini file into memory.
func LoadConfig(filename string) error {
	var err error
	configFile, err = ini.Load(filename)

	if err != nil {
		return err
	}

	return nil
}

// ConfigString returns the string value stored
// in the section under the key.
func ConfigString(section, key string) (string, error) {
	sectionObj := configFile.Section(section)

	if sectionObj == nil {
		return "", fmt.Errorf("section %s doesn't exist", section)
	}

	keyObj := sectionObj.Key(key)

	if keyObj == nil {
		return "", fmt.Errorf("key %s doesn't exist", key)
	}

	return keyObj.String(), nil
}

// ConfigBool returns the bool value stored
// in the section under the key.
func ConfigBool(section, key string) (bool, error) {
	sectionObj := configFile.Section(section)

	if sectionObj == nil {
		return false, fmt.Errorf("section %s doesn't exist", section)
	}

	keyObj := sectionObj.Key(key)

	if keyObj == nil {
		return false, fmt.Errorf("key %s doesn't exist", key)
	}

	return keyObj.Bool()
}

// ConfigInt returns the int value stored
// in the section under the key.
func ConfigInt(section, key string) (int, error) {
	sectionObj := configFile.Section(section)

	if sectionObj == nil {
		return -1, fmt.Errorf("section %s doesn't exist", section)
	}

	keyObj := sectionObj.Key(key)

	if keyObj == nil {
		return -1, fmt.Errorf("key %s doesn't exist", key)
	}

	return keyObj.Int()
}
