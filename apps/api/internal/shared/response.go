package shared

type OKResponse struct {
	OK bool `json:"ok"`
}

func OK() OKResponse {
	return OKResponse{OK: true}
}
