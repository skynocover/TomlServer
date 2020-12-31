package resp

import "encoding/json"

// Response ...
type Response struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Data         string `json:"data"`
}

// ToBytes ...
func (r *Response) ToBytes() []byte {

	jsondata, _ := json.Marshal(r)
	return jsondata
}
