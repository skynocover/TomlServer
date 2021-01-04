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

// ToBytesWithObject ...
func (r *Response) ToBytesWithObject(v map[string]string) []byte {
	m := map[string]string{}
	json.Unmarshal(r.ToBytes(), &m)
	for k, v := range v {
		m[k] = v
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}
