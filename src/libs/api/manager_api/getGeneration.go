package manager_api

const GetGenerationUrl = "/get_generation"

type Generation = uint16

type GetGenerationResponse struct {
	Generation Generation `json:"generation"`
}
