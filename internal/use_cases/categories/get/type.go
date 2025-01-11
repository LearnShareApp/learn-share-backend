package get

import "github.com/LearnShareApp/learn-share-backend/internal/jsonutils"

type response struct {
	Categories []category `json:"categories"`
}

type category struct {
	Name   string `json:"name" example:"Programing"`
	MinAge int    `json:"min_age" example:"12"`
}

type errorResponse jsonutils.ErrorStruct
