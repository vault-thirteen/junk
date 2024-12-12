package media_info

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"spdb_manager/settings"
	"strings"
)

type Response struct {
	Media ResponseMedia `json:"media"`
}

type ResponseMedia struct {
	Ref    string               `json:"@ref"`
	Tracks []ResponseMediaTrack `json:"track"`
}

type ResponseMediaTrack struct {
	AudioCount                uint    `json:"AudioCount,string"`
	BitRate                   uint    `json:"BitRate,string"`
	BitRate_Mode              string  `json:"BitRate_Mode"`
	BitRate_Maximum           string  `json:"BitRate_Maximum"` // Some movies have this field bugged.
	DisplayAspectRatio        float32 `json:"DisplayAspectRatio,string"`
	Duration                  float32 `json:"Duration,string"`
	FileExtension             string  `json:"FileExtension"`
	FileSize                  uint    `json:"FileSize,string"`
	Format                    string  `json:"Format"`
	Format_AdditionalFeatures string  `json:"Format_AdditionalFeatures"`
	Format_Profile            string  `json:"Format_Profile"`
	FrameCount                uint    `json:"FrameCount,string"`
	FrameRate                 float32 `json:"FrameRate,string"`
	Height                    uint    `json:"Height,string"`
	OverallBitRate            float32 `json:"OverallBitRate,string"`
	OverallBitRate_Mode       string  `json:"OverallBitRate_Mode"`
	PixelAspectRatio          float32 `json:"PixelAspectRatio,string"`
	Rotation                  float32 `json:"Rotation,string"`
	Sampled_Height            uint    `json:"Sampled_Height,string"`
	Sampled_Width             uint    `json:"Sampled_Width,string"`
	Stored_Height             uint    `json:"Stored_Height,string"`
	Stored_Width              uint    `json:"Stored_Width,string"`
	Type                      string  `json:"@type"`
	VideoCount                uint    `json:"VideoCount,string"`
	Width                     uint    `json:"Width,string"`
}

// GetExeVersion gets version text of the MediaInfo tool.
func GetExeVersion() (version string, err error) {
	cmd := fmt.Sprintf("%s", settings.ExeMediaInfo)

	var output []byte
	output, err = exec.Command(cmd, "--version").Output()
	if err != nil {
		return "", err
	}

	version = strings.TrimSpace(string(output))

	return version, nil
}

// GetData gets information about the specified media file using the MediaInfo
// tool.
func GetData(filePath string) (resp *Response, err error) {
	cmd := fmt.Sprintf("%s", settings.ExeMediaInfo)

	var output []byte
	output, err = exec.Command(cmd, "--output=JSON", filePath).Output()
	if err != nil {
		return nil, err
	}

	resp = new(Response)
	err = json.Unmarshal(output, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
