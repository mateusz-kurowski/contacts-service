package bucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Store struct {
	Client *s3.Client
	Bucket string
}

// OpenFromEnv initializes an S3 client pointed at OCI Object Storage S3-compatible endpoint.
func OpenFromEnv(ctx context.Context) (*Store, error) {
	endpoint := os.Getenv("OCI_S3_ENDPOINT")
	region := os.Getenv("OCI_S3_REGION")
	accessKey := os.Getenv("OCI_S3_ACCESS_KEY")
	secretKey := os.Getenv("OCI_S3_SECRET_KEY")
	bucket := os.Getenv("OCI_BUCKET_NAME")

	areEnvsEmpty := endpoint == "" || region == "" || accessKey == "" || secretKey == "" || bucket == ""

	if areEnvsEmpty {
		return nil, errors.New("missing required env vars for S3-compatible OCI: endpoint/region/credentials/bucket")
	}

	// Load a minimal AWS config with static credentials.
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg,
		func(o *s3.Options) {
			o.BaseEndpoint = aws.String("https://fro5sh7tqrh3.compat.objectstorage.eu-frankfurt-1.oraclecloud.com")
			o.UsePathStyle = true
		},
		s3.WithSigV4SigningRegion("eu-frankfurt-1"),
		s3.WithSigV4SigningName("s3"),
	)

	return &Store{Client: client, Bucket: bucket}, nil
}

// Upload writes bytes as an object to the OCI bucket via S3 API.
func (s *Store) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	if s == nil || s.Client == nil {
		return errors.New("nil store/client")
	}
	//nolint: exhaustruct // not necessary
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	_, err := s.Client.PutObject(ctx, input)
	return err
}

// Download retrieves an object and returns its bytes. Helper for future usage.
func (s *Store) Download(ctx context.Context, key string) ([]byte, error) {
	if s == nil || s.Client == nil {
		return nil, errors.New("nil store/client")
	}
	//nolint: exhaustruct // not necessary
	out, err := s.Client.GetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.Bucket), Key: aws.String(key)})
	if err != nil {
		return nil, err
	}
	defer out.Body.Close()
	b, readErr := io.ReadAll(out.Body)
	if readErr != nil {
		return nil, readErr
	}
	return b, nil
}

// GetStream retrieves an object handle from S3 for streaming.
// The caller is responsible for closing the returned object's Body.
func (s *Store) GetStream(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	if s == nil || s.Client == nil {
		return nil, errors.New("nil store/client")
	}

	//nolint: exhaustruct // not necessary
	out, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err // Return the S3 error directly
	}
	return out, nil
}
