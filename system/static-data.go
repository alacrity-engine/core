package system

import (
	"fmt"
)

// persistentStaticData is a storage to pass data between scenes.
var persistentStaticData map[string]interface{}

// StaticData returns the persistent data stored by
// the given key.
func StaticData(key string) (interface{}, error) {
	data, ok := persistentStaticData[key]

	if !ok {
		return nil, fmt.Errorf("no data stored by the key %s", key)
	}

	return data, nil
}

// SetStaticData saves the given value by the given key
// in the persistent storage.
func SetStaticData(key string, value interface{}) {
	persistentStaticData[key] = value
}

// RemoveStaticData removes the value by the given key
// from the persistent storage.
func RemoveStaticData(key string) {
	delete(persistentStaticData, key)
}
