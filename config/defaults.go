package config

const (

	// DefaultRetryCount is the default number of request retries.
	DefaultRetryCount = 2

	//DefaultClientID is the default ID used for Hosted API calls (Using UiPath's SaaS platform
	DefaultClientID = "5v7PmPJL6FOGu6RB8I1Y4adLBhIwovQN"
)

// RequestRetryCount returns the number of request retries.
func (*Config) RequestRetryCount() int {
	return DefaultRetryCount
}
