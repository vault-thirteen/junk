package models

import (
	"time"
)

type Movie struct {
	Id                 ID
	Title              string
	Description        string
	FileServerId       ID
	FilePath           string
	FileName           string
	FileExtension      string
	ContainerTypeId    ID
	OverallBitRate     float32
	VideoCodecId       ID
	AudioCodecId       ID
	FrameWidth         uint
	FrameHeight        uint
	Duration           float32
	TagsList           []ID
	ActorsList         []ID
	GenresList         []ID
	TimeOfRegistration time.Time
	TimeOfUpdate       *time.Time
	LastViewTime       *time.Time
	ViewsCount         uint
	CommentsList       []ID
	HardcoreLevel      byte
}

type MovieRaw struct {
	Id                 ID
	Title              string
	Description        string
	FileServerId       ID
	FilePath           string
	FileName           string
	FileExtension      string
	ContainerTypeId    ID
	OverallBitRate     float32
	VideoCodecId       ID
	AudioCodecId       ID
	FrameWidth         uint
	FrameHeight        uint
	Duration           float32
	TagsList           string
	ActorsList         string
	GenresList         string
	TimeOfRegistration int64
	TimeOfUpdate       *int64
	LastViewTime       *int64
	ViewsCount         uint
	CommentsList       string
	HardcoreLevel      byte
}

// Parse parses the raw object into a normal object.
func (mr *MovieRaw) Parse() (m *Movie, err error) {
	m = &Movie{
		Id:              mr.Id,
		Title:           mr.Title,
		Description:     mr.Description,
		FileServerId:    mr.FileServerId,
		FilePath:        mr.FilePath,
		FileName:        mr.FileName,
		FileExtension:   mr.FileExtension,
		ContainerTypeId: mr.ContainerTypeId,
		OverallBitRate:  mr.OverallBitRate,
		VideoCodecId:    mr.VideoCodecId,
		AudioCodecId:    mr.AudioCodecId,
		FrameWidth:      mr.FrameWidth,
		FrameHeight:     mr.FrameHeight,
		Duration:        mr.Duration,
		//TagsList:           mr.TagsList,
		//ActorsList:         mr.ActorsList,
		//GenresList:         mr.GenresList,
		//TimeOfRegistration: mr.TimeOfRegistration,
		//TimeOfUpdate:       mr.TimeOfUpdate,
		//LastViewTime:       mr.LastViewTime,
		ViewsCount: mr.ViewsCount,
		//CommentsList:       mr.CommentsList,
		HardcoreLevel: mr.HardcoreLevel,
	}

	if len(mr.TagsList) == 0 {
		mr.TagsList = EmptyList
	}
	m.TagsList, err = ParseJsonArrayOfIds(mr.TagsList)
	if err != nil {
		return nil, err
	}

	if len(mr.ActorsList) == 0 {
		mr.ActorsList = EmptyList
	}
	m.ActorsList, err = ParseJsonArrayOfIds(mr.ActorsList)
	if err != nil {
		return nil, err
	}

	if len(mr.GenresList) == 0 {
		mr.GenresList = EmptyList
	}
	m.GenresList, err = ParseJsonArrayOfIds(mr.GenresList)
	if err != nil {
		return nil, err
	}

	if mr.TimeOfRegistration != 0 {
		m.TimeOfRegistration = time.Unix(mr.TimeOfRegistration, 0)
	} else {
		m.TimeOfRegistration = time.Time{}
	}

	if mr.TimeOfUpdate != nil {
		tmp := time.Unix(*(mr.TimeOfUpdate), 0)
		m.TimeOfUpdate = &tmp
	}

	if mr.LastViewTime != nil {
		tmp := time.Unix(*(mr.LastViewTime), 0)
		m.LastViewTime = &tmp
	}

	if len(mr.CommentsList) == 0 {
		mr.CommentsList = EmptyList
	}
	m.CommentsList, err = ParseJsonArrayOfIds(mr.CommentsList)
	if err != nil {
		return nil, err
	}

	return m, nil
}
