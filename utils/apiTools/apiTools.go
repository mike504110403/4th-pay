package apitools

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
	"github.com/valyala/fasthttp"
)

// BuildFormValues : 將 GomypayRequest 結構轉換為 form-urlencoded 格式參數
func BuildFormValues(reqData interface{}) url.Values {
	formValues := url.Values{}

	v := reflect.ValueOf(reqData)
	t := reflect.TypeOf(reqData)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("form")

		if tag == "" {
			continue
		}

		// 只處理 string 型別，且非空值才加入
		if field.Kind() == reflect.String && field.String() != "" {
			formValues.Set(tag, field.String())
		}
	}

	return formValues
}

// PostForm 發送 application/x-www-form-urlencoded 的 POST 請求
func PostForm(endpoint string, formValues url.Values) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(endpoint)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetBodyString(formValues.Encode())

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ProcessUrl : 處理 Gomypay 交易 URL
func ProcessUrl(root string, reqData interface{}) string {
	processUrl := ""
	if v, err := query.Values(reqData); err != nil {
		fmt.Println(err)
		return processUrl
	} else {
		processUrl = fmt.Sprintf("%s?%s", root, v.Encode())
	}

	// processUrl := fmt.Sprintf(
	// 	"%s?Send_Type=%s&Amount=%s&Pay_Mode_No=%s&CustomerId=%s&Order_No=%s&StrCheck=%s",
	// 	TEST_URL,
	// 	reqData.SendType,
	// 	reqData.Amount,
	// 	reqData.PayModeNo,
	// 	reqData.CustomerID,
	// 	reqData.OrderNo,
	// 	reqData.StrCheck,
	// )

	return processUrl
}
