package oasaSyncWeb

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	oasaApplication = "http://telematics.oasa.gr"
)

func getRequest(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Errorf("got error %s", err.Error())
		return nil, err
	}

	if headers != nil && len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%s %d %s %s %s \n", time.Now().Format("2006-01-02 15:04:05"), response.StatusCode, strings.Split(url, "?")[1], req.Method, req.Host)

	// fmt.Printf("client: got response!\n")
	// fmt.Printf("client: status code: %d\n", response.StatusCode)
	return response, nil
}

func MakeRequest(action string) (string, error) {
	response, err := getRequest(oasaApplication+"/api/?act="+action,
		map[string]string{
			"Accept-Encoding": "gzip, deflate"})
	if err != nil {
		return "", err
	}

	reader, err := gzip.NewReader(response.Body)

	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	} else {
		defer reader.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		responseStr := buf.String()
		if response.StatusCode == http.StatusInternalServerError {
			fmt.Println("Response Body ", responseStr)
		}
		return responseStr, nil
	}
}

func OasaRequestApi(action string, extraParams map[string]interface{}) *OasaResponse {
	var oasaResult OasaResponse = OasaResponse{}
	var extraparamUrl string = ""
	// keys := make([]int, len(extraParams))
	for k := range extraParams {
		extraparamUrl = extraparamUrl + "&" + k + "=" + strconv.FormatInt(int64(extraParams[k].(int32)), 10)
	}
	//Error Code for error occured in Request Creation
	response, err := getRequest(oasaApplication+"/api/?act="+action+extraparamUrl, map[string]string{
		"Accept-Encoding": "gzip, deflate"})
	if err != nil {
		oasaResult.Error = &OasaError{
			Error_Code:  123,
			Error_Descr: err.Error(),
			Retry:       false,
		}
		return &oasaResult
	}

	result, error := ioutil.ReadAll(response.Body)
	if error != nil {
		oasaResult.Error = &OasaError{
			Error_Code:  124, //on Reading Response Body
			Error_Descr: err.Error(),
			Retry:       false,
		}
		return &oasaResult
	}

	var responseNonOk = response.StatusCode >= http.StatusInternalServerError &&
		response.StatusCode <= http.StatusNetworkAuthenticationRequired
	if responseNonOk {
		oasaResult.Error = &OasaError{
			Error_Code:  125, //on 5xx http Error Response
			Error_Descr: fmt.Sprintf("Internal Error Oasa Server (%s)", bytes.NewBuffer(result)),
			Retry:       true,
		}
		return &oasaResult
	}

	var tmpResult interface{}
	erro := json.Unmarshal(result, &tmpResult)
	if erro != nil {
		oasaResult.Error = &OasaError{
			Error_Code:  126, //on JSON Body Parsing
			Error_Descr: err.Error(),
		}
		return &oasaResult
	}
	var something = reflect.ValueOf(tmpResult)
	println(something.Bytes())
	// if reflect.ValueOf(tmpResult).FieldByName("error") {
	// 	println("No error Occured!!!")
	// } else {
	// 	println("An error occured!!!")
	// }
	oasaResult.Data = tmpResult

	return &oasaResult
}
