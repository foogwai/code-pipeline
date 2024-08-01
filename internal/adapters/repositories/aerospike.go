package repositories

import (
	"github.com/aerospike/aerospike-client-go/v7"

	"github.com/crseat/example-data-pipeline/internal/domain"
)

// AerospikeRepository represents a repository for storing post data in an Aerospike database.
type AerospikeRepository struct {
	client *aerospike.Client
}

// NewAerospikeRepository creates a new instance of AerospikeRepository and connects to the Aerospike database.
func NewAerospikeRepository(host string, port int) (*AerospikeRepository, error) {
	client, err := aerospike.NewClient(host, port)
	if err != nil {
		return nil, err
	}
	return &AerospikeRepository{client: client}, nil
}

// SavePostData saves the given post data to the Aerospike database.
func (r *AerospikeRepository) SavePostData(postData domain.PostData) error {
	key, err := aerospike.NewKey("test", "posts", postData.AdvertiserID)
	if err != nil {
		return err
	}

	bins := aerospike.BinMap{
		"ip_address":    postData.IPAddress,
		"user_agent":    postData.UserAgent,
		"referring_url": postData.ReferringURL,
		"metadata":      postData.Metadata,
	}

	err = r.client.Put(nil, key, bins)
	if err != nil {
		return err
	}

	return nil
}
