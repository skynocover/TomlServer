package resp

import "unsafe"

// Response ...
type Response struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Data         string `json:"data"`
}

// ToBytes ...
func (r *Response) ToBytes() []byte {

	return *(*[]byte)(unsafe.Pointer(r))
}
