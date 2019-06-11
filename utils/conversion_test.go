package utils

import (
	"fmt"
	"testing"
	"time"
)

type Par struct {
	Int     int                    `json:"int"`
	Bool    bool                   `json:"bool"`
	String  string                 `json:"string"`
	Float64 float64                `json:"float64"`
	Slice   []string               `json:"slice"`
	Time    interface{}            `json:"time"`
	Map     map[string]interface{} `json:"map"`
	Man     Man                    `json:"man"`
}

type Man struct {
	Gender string
	Height float64
}

func TestEncode_1(t *testing.T) {
	mapParameter := map[string]interface{}{
		"Int":     18,
		"Bool":    true,
		"String":  123.3,
		"Float64": 180,
		"Time":    time.Now(),
		"Slice":   []string{"slice_1", "slice_2"},
		// "Man": Man{
		// 	Gender: "男",
		// 	Height: 177.7,
		// },
		// "Map": map[string]interface{}{"map1": "map1", "map2": "map2"},
		// "Man": map[string]interface{}{"Gender": "男", "Height": 17.7},
	}

	structPointer := new(Par)
	structPointer.Slice = []string{"1", "2"}
	Encode(mapParameter, structPointer)
	fmt.Println(*structPointer)
	fmt.Println(structPointer.String)
	// mapstructure.Decode(nil, nil)
}

func TestEncode_2(t *testing.T) {
	mapParameter := []map[string]interface{}{
		{
			"Int":     18,
			"Bool":    true,
			"String":  123.3,
			"Float64": 180,
		},
	}
	structPointer := new([]Par)
	Encode_2(mapParameter, structPointer)
}

func TestDecode(t *testing.T) {
	structParameter := Par{
		Int:     18,
		Bool:    true,
		String:  "string",
		Float64: 180.8,
		Time:    time.Now(),
		Slice:   []string{"slice_1", "slice_2"},
		Man: Man{
			Gender: "男",
			Height: 177.77,
		},
	}
	mapParameter := Decode(&structParameter)
	fmt.Println(mapParameter)
}
