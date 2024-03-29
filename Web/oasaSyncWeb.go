package oasaSyncWeb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	logger "github.com/cs161079/godbLib/Utils/goLogger"
)

const (
	oasaApplication = "http://telematics.oasa.gr"
)

type OpswHttpRequest struct {
	Method   string
	Headers  map[string]string
	Body     io.Reader
	Endpoint string
}

func getProperty(v interface{}, property string) any {
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		return nil
	} else {
		result := v.(map[string]any)[property]
		return result
	}
}

func checkFields(request *OpswHttpRequest) error {
	if request.Endpoint == "" {
		return fmt.Errorf("REQUEST ENDPOINT IS NOT SET.")
	}
	if request.Method == "" {
		return fmt.Errorf("REQUEST HTTP METHOD IS NOT SET.")
	}
	return nil
}

func httpRequest(request *OpswHttpRequest) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	if &request == nil {
		return nil, fmt.Errorf("REQUEST OBJECT-STRUCT IS NIL OR IS NOT SET CORRECTLY.")
	}

	err := checkFields(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(request.Method, request.Endpoint, request.Body)
	if err != nil {
		return nil, err
	}

	if request.Headers != nil && len(request.Headers) > 0 {
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	logger.INFO(fmt.Sprintf("%s %s %d", response.Request.Method, response.Request.URL.Path, response.StatusCode))

	return response, nil
}

// func MakeRequest(action string) (string, error) {
// 	var req OpswHttpRequest = OpswHttpRequest{
// 		Method: http.MethodGet,
// 	}
// 	response, err := httpRequest(oasaApplication+"/api/?act="+action,
// 		map[string]string{
// 			"Accept-Encoding": "gzip, deflate"})
// 	if err != nil {
// 		return "", err
// 	}

// 	reader, err := gzip.NewReader(response.Body)

// 	if err != nil {
// 		fmt.Printf(err.Error())
// 		return "", err
// 	} else {
// 		defer reader.Close()

// 		buf := new(bytes.Buffer)
// 		buf.ReadFrom(reader)
// 		responseStr := buf.String()
// 		if response.StatusCode == http.StatusInternalServerError {
// 			fmt.Println("Response Body ", responseStr)
// 		}
// 		return responseStr, nil
// 	}
// }

func OasaRequestApi(action string, extraParams map[string]interface{}) *OasaResponse {
	var oasaResult OasaResponse = OasaResponse{}
	var extraparamUrl string = ""
	// keys := make([]int, len(extraParams))
	for k := range extraParams {
		extraparamUrl = extraparamUrl + "&" + k + "=" + strconv.FormatInt(int64(extraParams[k].(int32)), 10)
	}
	var req OpswHttpRequest = OpswHttpRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("%s/api/?act=%s%s", oasaApplication, action, extraparamUrl),
	}
	//Error Code for error occured in Request Creation
	response, err := httpRequest(&req)

	if err != nil {
		oasaResult.Error = err
	} else {
		if response.StatusCode >= http.StatusBadRequest && response.StatusCode <= http.StatusUnavailableForLegalReasons {
			fmt.Println("Client Error Response from Server")
			oasaResult.Error = fmt.Errorf("%s %s", response.Status, "Request contains bad syntax or cannot be fulfilled.")
		} else if response.StatusCode >= http.StatusInternalServerError && response.StatusCode <= http.StatusNetworkAuthenticationRequired {
			oasaResult.Error = fmt.Errorf("%s %s", response.Status, "Internal Server Error.")
		} else {
			responseBody, error := io.ReadAll(response.Body)
			if error != nil {
				oasaResult.Error = fmt.Errorf("AN ERROR OCCURED ANALYZE RESPONSE BODY. %s", error.Error())
			} else {
				var tmpResult interface{}
				err := json.Unmarshal(responseBody, &tmpResult)
				if err != nil {
					oasaResult.Error = fmt.Errorf("AN ERROR OCCURED WHEN CONVERT JSON STRING TO INTERFACE. %s", err.Error())
				} else {
					hasError := getProperty(tmpResult, "error")
					if hasError != nil {
						oasaResult.Error = fmt.Errorf("SERVER RESPONSES ERROR. %s", hasError)
					} else {
						oasaResult.Data = tmpResult
					}
				}
			}
		}
	}

	return &oasaResult
}
