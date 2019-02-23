package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type B struct {
	C string `example:"this is c"`
}

type D struct {
	Name string `example:"john"`
	*B
	Strings []int `example:"car,paper,this"`
}

type A struct {
	ID string `example:"this is id"`
	*B
	Numbers []int `example:"1,2,3"`
	Strings []int `example:"car,paper,this"`
	D
	IsValid bool `example:"true"`
	IsAge   bool `example:"false"`
	IsNone  bool
	Name    string   `json:"name" example:"hello"`
	IDs     []string `json:"ids" example:"1,2,3"`
	Age     int      `json:"age" example:"1"`
	Count   int64    `json:"count" example:"1000"`
	Pages   []int    `json:"pages" example:"1,2,3,4"`
	Pages64 []int64  `json:"pages64" example:"1,2,3,4"`
	Amount  float64  `json:"amount" example:"5.9"`

	Amounts []float64 `json:"amounts" example:"5.9,123"`
}

func main() {
	// Setting nested struct.
	parseTag := ParseTagFactory("example")
	{

		var a A
		out := parseTag(a)
		fmt.Println(out)
		switch o := out.(type) {
		case A:
			fmt.Printf("%#v\n", o)
			fmt.Printf("%#v\n", o.B)
			fmt.Printf("%#v\n", o.D)
			fmt.Printf("%#v\n", o.D.B)
		}
	}
	{
		var a A
		out := parseTag(a)
		fmt.Println(out)
		switch o := out.(type) {
		case A:
			fmt.Printf("%#v\n", o)
			fmt.Printf("%#v\n", o.B)
			fmt.Printf("%#v\n", o.D)
			fmt.Printf("%#v\n", o.D.B)
		}
		b, _ := json.Marshal(out)
		fmt.Println(string(b))
	}
}

var reflectSliceIntType = reflect.TypeOf([]int{})
var reflectSliceInt8Type = reflect.TypeOf([]int8{})
var reflectSliceInt16Type = reflect.TypeOf([]int16{})
var reflectSliceInt32Type = reflect.TypeOf([]int32{})
var reflectSliceInt64Type = reflect.TypeOf([]int64{})

var reflectSliceFloat32Type = reflect.TypeOf([]float32{})
var reflectSliceFloat64Type = reflect.TypeOf([]float64{})

func ParseTagFactory(tagName string) func(interface{}) interface{} {
	return func(in interface{}) interface{} {
		var parse func(in interface{}) reflect.Value
		parse = func(in interface{}) reflect.Value {
			v := reflect.ValueOf(in)
			switch k := v.Kind(); k {
			case reflect.Struct:

				t := reflect.TypeOf(in)
				newEl := reflect.New(t)
				el := newEl.Elem()

				for i := 0; i < el.NumField(); i++ {
					r := el.FieldByIndex([]int{i})
					tagValue := t.Field(i).Tag.Get(tagName)

					switch r.Kind() {
					case reflect.Struct:
						newStruct := reflect.New(reflect.TypeOf(r.Interface()))
						out := parse(newStruct.Elem().Interface())
						r.Set(reflect.ValueOf(out.Interface()))
					case reflect.Ptr:
						newStruct := reflect.New(reflect.TypeOf(r.Interface()).Elem())
						out := parse(newStruct.Elem().Interface())
						r.Set(out.Addr())
					case reflect.Bool:
						r.SetBool(tagValue == "true")
					case reflect.Slice:

						val := strings.Split(tagValue, ",")
						switch r.Type() {
						case reflectSliceIntType:
							out := make([]int, len(val))
							for i, v := range val {
								out[i], _ = strconv.Atoi(v)
							}
							r.Set(reflect.ValueOf(out))
						case reflectSliceInt8Type:
							out := make([]int8, len(val))
							for i, v := range val {
								o, _ := strconv.ParseInt(v, 10, 8)
								out[i] = int8(o)
							}
							r.Set(reflect.ValueOf(out))
						case reflectSliceInt16Type:
							out := make([]int16, len(val))
							for i, v := range val {
								o, _ := strconv.ParseInt(v, 10, 16)
								out[i] = int16(o)
							}
							r.Set(reflect.ValueOf(out))

						case reflectSliceInt32Type:
							out := make([]int32, len(val))
							for i, v := range val {
								o, _ := strconv.ParseInt(v, 10, 32)
								out[i] = int32(o)
							}
							r.Set(reflect.ValueOf(out))
						case reflectSliceInt64Type:
							out := make([]int64, len(val))
							for i, v := range val {
								out[i], _ = strconv.ParseInt(v, 10, 64)
							}
							r.Set(reflect.ValueOf(out))

						case reflectSliceFloat32Type:
							out := make([]float32, len(val))
							for i, v := range val {
								o, _ := strconv.ParseFloat(v, 32)
								out[i] = float32(o)
							}
							r.Set(reflect.ValueOf(out))

						case reflectSliceFloat64Type:
							out := make([]float64, len(val))
							for i, v := range val {
								out[i], _ = strconv.ParseFloat(v, 64)
							}
							r.Set(reflect.ValueOf(out))
						default:
							r.Set(reflect.ValueOf(val))
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						val := parseInt(r.Kind(), tagValue)
						r.Set(reflect.ValueOf(val))
					case reflect.Float32, reflect.Float64:
						val := parseFloat(r.Kind(), tagValue)
						r.Set(reflect.ValueOf(val))
					default:
						r.Set(reflect.ValueOf(tagValue))
					}
				}
				return el
			case reflect.Ptr:
				r := reflect.New(reflect.TypeOf(in).Elem())
				out := parse(r.Elem().Interface())
				return out
			default:
				return v
			}
		}
		out := parse(in)
		return out.Interface()
	}
}

func parseInt(k reflect.Kind, s string) interface{} {
	switch k {
	case reflect.Int:
		val, _ := strconv.Atoi(s)
		return val
	case reflect.Int8:
		val, _ := strconv.ParseInt(s, 10, 8)
		return val
	case reflect.Int16:
		val, _ := strconv.ParseInt(s, 10, 16)
		return val
	case reflect.Int32:
		val, _ := strconv.ParseInt(s, 10, 32)
		return val
	case reflect.Int64:
		val, _ := strconv.ParseInt(s, 10, 64)
		return val
	}
	return 0
}

func parseFloat(k reflect.Kind, s string) interface{} {
	switch k {
	case reflect.Float32:
		val, _ := strconv.ParseFloat(s, 32)
		return val
	case reflect.Float64:
		val, _ := strconv.ParseFloat(s, 64)
		return val
	}
	return 0
}
