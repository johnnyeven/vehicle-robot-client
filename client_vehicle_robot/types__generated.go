package client_vehicle_robot

import (
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
	github_com_johnnyeven_libtools_courier_swagger "github.com/johnnyeven/libtools/courier/swagger"
)

type DetectivedObject struct {
	//
	Box []float32 `json:"box"`
	//
	Class float32 `json:"class"`
	//
	Probability float32 `json:"probability"`
}

type ErrorField = github_com_johnnyeven_libtools_courier_status_error.ErrorField

type ErrorFields = github_com_johnnyeven_libtools_courier_status_error.ErrorFields

type JSONBytes = github_com_johnnyeven_libtools_courier_swagger.JSONBytes

type ObjectDetectionBody struct {
	//
	Image []uint8 `json:"image"`
}

type StatusError = github_com_johnnyeven_libtools_courier_status_error.StatusError
