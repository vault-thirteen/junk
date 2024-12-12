package models

import "encoding/json"

// FuncProcessorSettingsJson is a JSON wrapper of FuncProcessorSettings.
type FuncProcessorSettingsJson struct {
	IgnoredWords     []string               `json:"ignored_words"`
	FuncNamePrefix   string                 `json:"func_name_prefix"`
	CToGoTypeMapping []CToGoTypeMappingJson `json:"c_to_go_type_mapping"`
}

type CToGoTypeMappingJson = []string

// NewFuncProcessorSettingsJson is a constructor of FuncProcessorSettingsJson.
func NewFuncProcessorSettingsJson(jsonData []byte) (fpsj *FuncProcessorSettingsJson, err error) {
	fpsj = new(FuncProcessorSettingsJson)
	err = json.Unmarshal(jsonData, fpsj)
	if err != nil {
		return nil, err
	}

	return fpsj, nil
}
