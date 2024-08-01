package domain

// PostData represents the structure of the data to be processed and stored.
// It includes information about a user's interaction and associated metadata.
//
// Fields:
//   - IPAddress: The IP address of the user, must be a valid IPv4 address.
//   - UserAgent: The user agent string from the HTTP request, indicating the user's browser or client.
//   - ReferringURL: The URL from which the user was referred, must be a valid URL.
//   - AdvertiserID: A unique identifier for the advertiser related to the post data.
//   - Metadata: Additional metadata associated with the post data, stored as a map of string keys to values of any type.
type PostData struct {
	IPAddress    string                 `json:"ip_address" validate:"required,ipv4"`
	UserAgent    string                 `json:"user_agent" validate:"required"`
	ReferringURL string                 `json:"referring_url" validate:"required,url"`
	AdvertiserID string                 `json:"advertiser_id" validate:"required"`
	Metadata     map[string]interface{} `json:"metadata" validate:"required"`
}
