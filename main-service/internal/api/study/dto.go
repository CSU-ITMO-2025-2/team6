package study

type CreateStudyResponse struct {
	ID string `json:"id"`
}

type GetStudyResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}
