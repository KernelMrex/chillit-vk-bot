package vkapi

type Response struct {
	Response int        `json:"reponse"`
	Error    *RespError `json:"error"`
}

type RespError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_msg"`
}
