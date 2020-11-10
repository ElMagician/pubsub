// Package google implement elMagician pubsub interfaces for GCP Pubsub provider.
package google

import (
	"time"
)

// Config for pubsub instance
type Config struct {
	// ProjectID in GCP
	ProjectID string `yaml:"projectId" json:"projectId"`

	// CredentialsPath to your JSON credential provided by GCP.
	CredentialsPath string `yaml:"credentialsPath" json:"credentialsPath"`

	// Concurrency is the default concurrency for listening process.
	Concurrency int `yaml:"concurrency" json:"concurrency"`

	// Timeout is the default timeout for GCP calls
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}
