package vkapi

// Response vk reponse struct
type Response struct {
	Response int        `json:"reponse"`
	Error    *RespError `json:"error"`
}

// RespError vk response error detailed info struct
type RespError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}
