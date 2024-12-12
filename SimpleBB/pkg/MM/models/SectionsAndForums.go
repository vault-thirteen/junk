package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
)

type SectionsAndForums struct {
	Sections []derived2.ISection `json:"sections"`
	Forums   []derived2.IForum   `json:"forums"`
}
