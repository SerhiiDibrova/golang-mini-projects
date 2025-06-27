package main

import (
	"errors"
	"strings"
)

type CommonFields struct {
	LayerName    string
	ShapeFileURL string
	OtherDetails string
}

func _loadCommonFields(params map[string]string) (map[string]string, error) {
	if params == nil {
		return nil, errors.New("input map is nil")
	}

	commonFields := make(map[string]string)

	layerName, ok := params["layer_name"]
	if !ok {
		return nil, errors.New("layer_name is required")
	}
	commonFields["layer_name"] = layerName

	shapeFileURL, ok := params["shape_file_url"]
	if !ok {
		return nil, errors.New("shape_file_url is required")
	}
	commonFields["shape_file_url"] = shapeFileURL

	otherDetails, ok := params["other_details"]
	if !ok {
		return nil, errors.New("other_details is required")
	}
	commonFields["other_details"] = otherDetails

	return commonFields, nil
}

func main() {
	params := map[string]string{
		"layer_name":    "test_layer",
		"shape_file_url": "https://example.com/shapefile.shp",
		"other_details":  "test_details",
	}

	commonFields, err := _loadCommonFields(params)
	if err != nil {
		panic(err)
	}

	for key, value := range commonFields {
		println(key + ": " + value)
	}
}