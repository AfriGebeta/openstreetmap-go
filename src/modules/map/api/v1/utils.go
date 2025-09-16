package v1

import (
	"errors"
	"openstreetmap-go/src/config"
	"openstreetmap-go/src/repository/model"
	"strconv"
	"strings"
)

var LON_LIMIT = 180.0
var LAT_LIMIT = 90.0

var SCALED_LON_LIMIT = LON_LIMIT * float64(config.ServerConfigObject.GeoRecordScale)
var SCALED_LAT_LIMIT = LAT_LIMIT * float64(config.ServerConfigObject.GeoRecordScale)

func FromBBoxParams(params string) (model.BBox, error) {
	bbox_array := strings.Split(params, ",")
	if len(bbox_array) < 4 {
		return model.BBox{}, errors.New("The parameter bbox must be of the form min_lon,min_lat,max_lon,max_lat, with all values present")
	}

	minLon, err := strconv.ParseFloat(bbox_array[0], 64)
	if err != nil {
		return model.BBox{}, err
	}
	minLat, err := strconv.ParseFloat(bbox_array[1], 64)
	if err != nil {
		return model.BBox{}, err
	}

	maxLon, err := strconv.ParseFloat(bbox_array[2], 64)
	if err != nil {
		return model.BBox{}, err
	}

	maxLat, err := strconv.ParseFloat(bbox_array[3], 64)
	if err != nil {
		return model.BBox{}, err
	}

	return model.BBox{
		MinLat: minLat,
		MinLon: minLon,
		MaxLat: maxLat,
		MaxLon: maxLon,
	}, nil
}

func Checkboundaries(bbox model.BBox) error {
	// check min_lon and max lon
	if bbox.MinLon > bbox.MaxLon {
		return errors.New("The boundaries are not in range")
	}

	// check the minLat nad maxLat
	if bbox.MinLat > bbox.MaxLat {
		return errors.New("The boundaries are not in range")
	}

	if bbox.MinLat < -LAT_LIMIT || bbox.MinLat > LAT_LIMIT ||
		bbox.MaxLat < -LAT_LIMIT || bbox.MaxLat > LAT_LIMIT ||
		bbox.MinLon < -LON_LIMIT || bbox.MinLon > LON_LIMIT ||
		bbox.MaxLon < -LON_LIMIT || bbox.MaxLon > LON_LIMIT {
		return errors.New("The boundaries are not in range")
	}
	return nil
}

func CheckAreaSize(bbox model.BBox) error {
	if area(bbox) > config.ServerConfigObject.MaxRequestArea {
		return errors.New("Maximum Request Area exceeded")
	}
	return nil
}

func area(bbox model.BBox) float64 {
	return (bbox.MaxLon - bbox.MinLon) * (bbox.MaxLat - bbox.MinLat)
}
