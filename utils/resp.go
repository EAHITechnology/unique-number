package utils

import "encoding/json"

type RespData struct {
	ErrorCode int64       `json:"err_code"`
	ErrorMsg  string      `json:"err_msg"`
	Data      interface{} `json:"data"`
}

func GetResponse(errCode int64, errMsg string, data interface{}) ([]byte, error) {
	r := RespData{
		ErrorCode: errCode,
		ErrorMsg:  errMsg,
		Data:      data,
	}
	body, err := json.Marshal(r)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
