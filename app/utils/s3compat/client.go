package s3compat

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	endpointURL string
	client      *minio.Client
}

func NewClient(endpointURL string, accessKeyID string, secretAccessKey string, region string, _ time.Duration) *Client {
	endpointURL = strings.TrimSpace(endpointURL)
	accessKeyID = strings.TrimSpace(accessKeyID)
	secretAccessKey = strings.TrimSpace(secretAccessKey)
	region = strings.TrimSpace(region)
	if endpointURL == "" || accessKeyID == "" || secretAccessKey == "" {
		return &Client{endpointURL: endpointURL}
	}

	parsed, err := url.Parse(endpointURL)
	if err != nil || parsed.Host == "" {
		return &Client{endpointURL: endpointURL}
	}

	secure := parsed.Scheme == "https"
	mc, err := minio.New(parsed.Host, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:       secure,
		Region:       region,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return &Client{endpointURL: endpointURL}
	}

	return &Client{endpointURL: strings.TrimRight(endpointURL, "/"), client: mc}
}

func (c *Client) Enabled() bool {
	return c != nil && c.client != nil
}

func (c *Client) EndpointURL() string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.endpointURL)
}

func (c *Client) PresignGetObject(bucket string, objectPath string, expires time.Duration) (string, error) {
	if c == nil || c.client == nil {
		return "", fmt.Errorf("s3compat-disabled")
	}
	bucket = strings.TrimSpace(bucket)
	objectPath = strings.Trim(strings.TrimSpace(objectPath), "/")
	if bucket == "" || objectPath == "" {
		return "", fmt.Errorf("s3compat-invalid-params")
	}
	if expires <= 0 {
		expires = 15 * time.Minute
	}

	u, err := c.client.PresignedGetObject(context.Background(), bucket, objectPath, expires, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (c *Client) PublicObjectURL(bucket string, objectPath string) string {
	if c == nil {
		return ""
	}
	endpoint := strings.TrimRight(strings.TrimSpace(c.endpointURL), "/")
	bucket = strings.TrimSpace(bucket)
	objectPath = strings.Trim(strings.TrimSpace(objectPath), "/")
	if endpoint == "" || bucket == "" || objectPath == "" {
		return ""
	}
	return endpoint + "/" + bucket + "/" + objectPath
}
