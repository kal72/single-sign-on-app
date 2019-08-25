package builder

import "github.com/sirupsen/logrus"

//response message
const (
	MessageFetchTrxFailed   = "Failed to get transactions data"
	MessageFetchTrxSuccess  = "Transactions data retrieved successfully"
	MessageAuthFailed       = "Authentication failed"
	MessageAuthSuccess      = "Authentication success"
	MessageErrorCreateToken = "Token not created"
	MessageUpdateSuccess    = "Update success"
	MessageSaveSuccess      = "Save success"
	MessageLogoutSuccess    = "Logout success"
	MessageLogoutFailed     = "Logout failed"
	MessageDeleteSuccess    = "Delete success"
)

/** helper for generate doc api **/
type DocResponse struct {
	Message string   `json:"status"`
	Data    struct{} `json:"data"`
}

type DocTransactionHelper struct {
	Merchant             string
	Device               string
	Emoney               string
	TransactionDateStart string
	TransactionDateEnd   string
}

type DocDeviceConfigResponse struct {
	Data     []interface{} `json:"data"`
	ProcType string        `json:"procType"`
	Prop     struct {
		Muid string `json:"muid"`
	} `json:"prop"`
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
	Version  string `json:"version"`
}

type DocLoginDeviceResponse struct {
	Data struct {
		RefreshToken string `json:"refresh_token"`
		Token        string `json:"token"`
		UserData     []struct {
			Fullname     string `json:"fullname"`
			GroupAddress string `json:"group_address"`
			GroupLogo    string `json:"group_logo"`
			Permission   string `json:"permission"`
			Role         string `json:"role"`
		} `json:"userData"`
	} `json:"data"`
	ProcType string `json:"procType"`
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
	Version  string `json:"version"`
}

type LoginDeviceResponse struct {
	ProcType string      `json:"procType"`
	Version  string      `json:"version"`
	RespCode string      `json:"respCode"`
	RespMsg  string      `json:"respMsg"`
	Data     interface{} `json:"data"`
}

/** Deprecated
* please, use Response function
 */
func BaseResponse(success bool, message string, count interface{}, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"status":  success,
		"count":   count,
		"data":    data}
}

func Response(success bool, message string, data interface{}) map[string]interface{} {
	logrus.Debug(message)

	var responseMap = make(map[string]interface{})
	//var status = "success"

	if !success {
		//status = "failed"
		responseMap["status"] = success
		responseMap["message"] = message
	} else {
		responseMap["status"] = success
		responseMap["data"] = data
		//responseMap["message"] = message
	}

	return responseMap
}
