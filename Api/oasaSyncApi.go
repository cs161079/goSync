package api

import (
	"fmt"

	mapper "github.com/cs161079/goSync/Mappers"
	oasaSyncWeb "github.com/cs161079/goSync/Web"
	models "github.com/cs161079/godbLib/Models"
)

type apierror interface {
	Http_Error()
}

type testType struct {
	Prop1 string
}

func (t testType) Http_Error() {
	fmt.Println(t.Prop1)
}

// func GetData(path string) ([]map[string]interface{}, error) {
// 	var result []map[string]interface{}
// 	responseStr, error := oasaSyncWeb.MakeRequest(path)
// 	if error != nil {
// 		return nil, error
// 	}
// 	file, error := os.Create("/oasa-telematics/" + path + "_data.txt")
// 	if error != nil {
// 		return nil, error
// 	}
// 	defer file.Close()
// 	result = oasaSyncDecode.ReadTextCharByChar(responseStr, nil, file)

// 	if len(result) == 0 {
// 		result = append(result, map[string]interface{}{
// 			"response": responseStr,
// 		})
// 	}
// 	return result, nil
// }

// func GetLines() ([]map[string]interface{}, error) {
// 	var result []map[string]interface{}

// 	responseStr, error := oasaSyncWeb.MakeRequest("getLines")
// 	if error != nil {
// 		return nil, error
// 	}
// 	f, err := os.Create("/oasa-telematics/lines_data.txt")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	result = oasa_sync_decode.ReadTextCharByChar(responseStr, func(dataStr []string) map[string]interface{} {
// 		var resultRec map[string]interface{}
// 		if len(dataStr) == 18 {
// 			resultRec = map[string]interface{}{
// 				"ID":       oasa_sync_utils.StrToInt64(dataStr[0]),
// 				"CODE":     dataStr[1],
// 				"DESCR":    dataStr[2],
// 				"DESCRENG": dataStr[3],
// 				"NUM1":     oasa_sync_utils.StrToInt(dataStr[4]),
// 				"NUM2":     oasa_sync_utils.StrToFloat(dataStr[5]),
// 				"NUM3":     oasa_sync_utils.StrToInt(dataStr[6]),
// 				"NUM4":     oasa_sync_utils.StrToFloat(dataStr[7]),
// 				"NUM5":     oasa_sync_utils.StrToInt(dataStr[8]),
// 				"NUM6":     oasa_sync_utils.StrToFloat(dataStr[9]),
// 				"NUM7":     oasa_sync_utils.StrToInt(dataStr[10]),
// 				"NUM8":     oasa_sync_utils.StrToFloat(dataStr[11]),
// 				"NUM9":     oasa_sync_utils.StrToInt(dataStr[12]),
// 				"NUM10":    oasa_sync_utils.StrToFloat(dataStr[13]),
// 				"NUM11":    oasa_sync_utils.StrToInt(dataStr[14]),
// 				"NUM12":    oasa_sync_utils.StrToFloat(dataStr[15]),
// 				"NUM13":    oasa_sync_utils.StrToInt(dataStr[16]),
// 				"NUM14":    oasa_sync_utils.StrToFloat(dataStr[17]),
// 			}

// 		}
// 		return resultRec
// 	}, f)

// 	//jsonStr, error := json.Marshal(result)
// 	//if error != nil {
// 	//	return "", err
// 	//}
// 	return result, nil
// }

// func GetRoutes() ([]map[string]interface{}, error) {
// 	var result []map[string]interface{}

// 	responseStr, error := oasaSyncWeb.MakeRequest("getRoutes")
// 	if error != nil {
// 		return nil, error
// 	}
// 	result = append(result, map[string]interface{}{
// 		"response": responseStr,
// 	})
// 	f, err := os.Create("/oasa-telematics/routes_data.txt")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	result = oasa_sync_decode.ReadTextCharByChar(responseStr, func(dataStr []string) map[string]interface{} {
// 		var resultRec map[string]interface{}
// 		if len(dataStr) == 6 {
// 			resultRec = map[string]interface{}{
// 				"RouteCode":     oasa_sync_utils.StrToInt64(dataStr[0]),
// 				"LineCode":      oasa_sync_utils.StrToInt64(dataStr[1]),
// 				"RouteDescr":    dataStr[2],
// 				"RouteDescrEng": dataStr[3],
// 				"RouteType":     oasa_sync_utils.StrToInt(dataStr[4]),
// 				"RouteDistance": oasa_sync_utils.StrToFloat(dataStr[5]),
// 			}

// 		}
// 		return resultRec
// 	}, f)
// 	//
// 	////jsonStr, error := json.Marshal(result)
// 	//if error != nil {
// 	//	return "", err
// 	//}
// 	return result, nil
// }

// func GetStops() ([]map[string]interface{}, error) {
// 	var result []map[string]interface{}

// 	responseStr, error := oasaSyncWeb.MakeRequest("getStops")
// 	if error != nil {
// 		return nil, error
// 	}
// 	f, err := os.Create("/oasa-telematics/stops-data.txt")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()
// 	result = append(result, map[string]interface{}{
// 		"response": responseStr,
// 	})
// 	result = oasa_sync_decode.ReadTextCharByChar(responseStr, func(dataStr []string) map[string]interface{} {
// 		var resultRec map[string]interface{}
// 		if len(dataStr) == 13 {
// 			resultRec = map[string]interface{}{
// 				"StopCode":      oasa_sync_utils.StrToInt64(dataStr[0]),
// 				"StopID":        dataStr[1],
// 				"StopDescr":     dataStr[2],
// 				"StopDescrEng":  dataStr[3],
// 				"StopStreet":    dataStr[4],
// 				"StopStreetEng": dataStr[5],
// 				"StopHeading":   dataStr[6],
// 				"StopLat":       oasa_sync_utils.StrToFloat(dataStr[7]),
// 				"StopLng":       oasa_sync_utils.StrToFloat(dataStr[8]),
// 				//"RouteStopOrder": oasa_sync_utils.StrToInt(dataStr[9]),
// 				"StopType":    dataStr[9],
// 				"StopAmea":    dataStr[10],
// 				"DescrStr":    dataStr[11],
// 				"DescrStrEng": dataStr[12],
// 			}

// 		}
// 		return resultRec
// 	}, f)
// 	//
// 	////jsonStr, error := json.Marshal(result)
// 	//if error != nil {
// 	//	return "", err
// 	//}
// 	return result, nil
// }

// This function call OASA Server to get all information about Bus Lines
// func GetBusLines(mappinOut func(interface{}) interface{}) ([]oasa_sync_model.LineDto, error) {
// 	var result []oasa_sync_model.LineDto
// 	_, error := oasaSyncWeb.OasaRequestApi("webGetLinesWithMLInfo", nil)
// 	if error != nil {
// 		return nil, error
// 	}

// 	return mappinOut(result).([]oasa_sync_model.LineDto), nil
// }

// This function call OASA Server to get all information about Bus Lines and store in general Interface{}
func GetBusLinesTest() ([]models.Busline, error) {
	var result []models.Busline
	response := oasaSyncWeb.OasaRequestApi("webGetLinesWithMLInfoo", nil)
	if response.Error != nil {
		return nil, response.Error
	}

	for _, record := range response.Data.([]interface{}) {
		result = append(result, mapper.BussLineMapper(record.(map[string]interface{})))
	}

	return result, nil
}

func GetBusRoutes(lined_id int32) ([]models.BusRoute, error) {
	var result []models.BusRoute
	response := oasaSyncWeb.OasaRequestApi("webGetRoutes", map[string]interface{}{"p1": lined_id})
	if response.Error != nil {
		return nil, response.Error
	}

	for _, record := range response.Data.([]interface{}) {
		// fmt.Println("at Index Route ", index)
		result = append(result, mapper.BusRouteMapper(record.(map[string]interface{})))
	}
	fmt.Printf("Routes results %d \n", len(result))

	return result, nil
}

func GetBusStops(route_code int32) ([]models.BusStopDto, error) {
	var result []models.BusStopDto
	response := oasaSyncWeb.OasaRequestApi("webGetStops", map[string]interface{}{"p1": route_code})
	if response.Error != nil {
		return nil, response.Error
	}

	for _, record := range response.Data.([]interface{}) {
		// fmt.Println("at Index Route ", index)
		result = append(result, mapper.BusStopDtoMapper(record.(map[string]interface{})))
	}
	fmt.Printf("Routes results %d \n", len(result))

	return result, nil
}

func GetBusScheduleMaster(line_code int32) ([]models.BusScheduleMasterLineDto, error) {
	var result []models.BusScheduleMasterLineDto
	response := oasaSyncWeb.OasaRequestApi("getScheduleDaysMasterline", map[string]interface{}{"p1": line_code})

	if response.Error != nil {
		return nil, response.Error
	}

	for _, record := range response.Data.([]interface{}) {
		result = append(result, mapper.SheduleMasterLineDtoMapper(record.(map[string]interface{})))
	}
	fmt.Printf("Schedule Master Line  get successfully %d \n", len(result))

	return result, nil

}
