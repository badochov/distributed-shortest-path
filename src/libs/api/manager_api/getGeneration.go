package manager_api

const GetGenerationUrl = "/get_generation"

type GetGenerationResponse struct {
	Generation int `json:"generation"`
}