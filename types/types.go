package types

// CreateSecretPayload is the payload for creating a secret.
type CreateSecretPayload struct {
	PlainText string `json:"plain_text"`
}

// CreateSecretResponse is the response for creating a secret.
type CreateSecretResponse struct {
	Id string `json:"id"`
}

// GetSecretPayload is the payload for getting a secret.
type GetSecretResponse struct {
	Data string `json:"data"`
}

// SecretData is the data for a secret.
type SecretData struct {
	Id     string
	Secret string
}
