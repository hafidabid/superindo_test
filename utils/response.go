package utils

type Response struct {
	Status   int         `json:"status"`
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
	Metadata Metadata    `json:"metadata"`
}

type Metadata struct {
	Page      int     `json:"page"`
	Limit     int     `json:"limit"`
	TotalPage float64 `json:"totalPage"`
	TotalData int64   `json:"totalData"`
}
