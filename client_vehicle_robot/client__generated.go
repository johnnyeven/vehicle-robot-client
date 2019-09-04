package client_vehicle_robot

import (
	"fmt"

	github_com_johnnyeven_libtools_courier "github.com/johnnyeven/libtools/courier"
	github_com_johnnyeven_libtools_courier_client "github.com/johnnyeven/libtools/courier/client"
	github_com_johnnyeven_libtools_courier_status_error "github.com/johnnyeven/libtools/courier/status_error"
)

type ClientVehicleRobotInterface interface {
	ObjectDetection(req ObjectDetectionRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *ObjectDetectionResponse, err error)
}

type ClientVehicleRobot struct {
	github_com_johnnyeven_libtools_courier_client.Client
}

func (ClientVehicleRobot) MarshalDefaults(v interface{}) {
	if cl, ok := v.(*ClientVehicleRobot); ok {
		cl.Name = "vehicle-robot"
		cl.Client.MarshalDefaults(&cl.Client)
	}
}

func (c ClientVehicleRobot) Init() {
	c.CheckService()
}

func (c ClientVehicleRobot) CheckService() {
	err := c.Request(c.Name+".Check", "HEAD", "/", nil).
		Do().
		Into(nil)
	statusErr := github_com_johnnyeven_libtools_courier_status_error.FromError(err)
	if statusErr.Code == int64(github_com_johnnyeven_libtools_courier_status_error.RequestTimeout) {
		panic(fmt.Errorf("service %s have some error %s", c.Name, statusErr))
	}
}

type ObjectDetectionRequest struct {
	//
	Body ObjectDetectionBody `fmt:"json" in:"body"`
}

func (c ClientVehicleRobot) ObjectDetection(req ObjectDetectionRequest, metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *ObjectDetectionResponse, err error) {
	resp = &ObjectDetectionResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".ObjectDetection", "POST", "/vehicle-robot/v0/detections/object", req, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type ObjectDetectionResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body []DetectivedObject
}

func (c ClientVehicleRobot) Swagger(metas ...github_com_johnnyeven_libtools_courier.Metadata) (resp *SwaggerResponse, err error) {
	resp = &SwaggerResponse{}
	resp.Meta = github_com_johnnyeven_libtools_courier.Metadata{}

	err = c.Request(c.Name+".Swagger", "GET", "/vehicle-robot", nil, metas...).
		Do().
		BindMeta(resp.Meta).
		Into(&resp.Body)

	return
}

type SwaggerResponse struct {
	Meta github_com_johnnyeven_libtools_courier.Metadata
	Body JSONBytes
}
