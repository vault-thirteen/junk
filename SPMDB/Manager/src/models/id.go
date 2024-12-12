package models

import (
	"encoding/json"
	"fmt"
)

type ID = uint

// ParseJsonArrayOfIds parses the text in JSON format into an array of IDs.
func ParseJsonArrayOfIds(s string) (a []ID, err error) {
	err = json.Unmarshal([]byte(s), &a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// EncodeArrayOfIdsAsJson encodes an array of IDs into the text in JSON format.
func EncodeArrayOfIdsAsJson(a []ID) (s string, err error) {
	var tmp []byte
	tmp, err = json.Marshal(a)
	if err != nil {
		return "", err
	}

	return string(tmp), nil
}

// AddIdToList adds an ID to a list of IDs.
func AddIdToList(list []ID, newId ID) (updatedList []ID, err error) {
	for _, id := range list {
		if id == newId {
			return nil, fmt.Errorf("duplicate id: %v", id)
		}
	}

	return append(list, newId), nil
}

// RemoveIdFromList removes an ID from a list of IDs.
func RemoveIdFromList(list []ID, targetId ID) (updatedList []ID, err error) {
	// We are using a map to remove possible duplicate items.
	var cache = make(map[ID]bool)
	for _, id := range list {
		cache[id] = true
	}

	// Search for an existing item to be deleted.
	var itemExists bool
	_, itemExists = cache[targetId]
	if !itemExists {
		return nil, fmt.Errorf("absent id: %v", targetId)
	}

	delete(cache, targetId)

	updatedList = make([]ID, 0, len(list)-1)
	for id, _ := range cache {
		updatedList = append(updatedList, id)
	}

	return updatedList, nil
}
