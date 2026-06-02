package datastore

import (
	"context"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3PresignerInterface interface {
	GetObject(bucket, key string, lifetime time.Duration) (*url.URL, error)
	PostObject(bucket, key string, lifetime time.Duration) (*url.URL, map[string]string, error)
	PutObject(bucket, key string, lifetime time.Duration) (*url.URL, error)
}

type S3Presigner struct {
	s3 *s3.Client
}

func (p *S3Presigner) GetObject(bucket, key string, lifetime time.Duration) (*url.URL, error) {

	presigner := s3.NewPresignClient(p.s3)

	presignedReq, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{Bucket: &bucket, Key: &key}, func(po *s3.PresignOptions) { po.Expires = lifetime })

	if err != nil {
		return nil, err
	}

	presignedURL, err := url.Parse(presignedReq.URL)
	if err != nil {
		return nil, err
	}
	return presignedURL, nil
}

func (p *S3Presigner) PostObject(bucket, key string, lifetime time.Duration) (*url.URL, map[string]string, error) {
	presignClient := s3.NewPresignClient(p.s3)

	presignedReq, err := presignClient.PresignPostObject(context.Background(), &s3.PutObjectInput{Bucket: &bucket, Key: &key}, func(ppo *s3.PresignPostOptions) { ppo.Expires = lifetime })

	if err != nil {
		return nil, nil, err
	}

	presignedURL, err := url.Parse(presignedReq.URL)

	if err != nil {
		return nil, nil, err
	}

	return presignedURL, presignedReq.Values, nil
}

func (p *S3Presigner) PutObject(bucket, key string, lifetime time.Duration) (*url.URL, error) {

	presignClient := s3.NewPresignClient(p.s3)

	presignedReq, err := presignClient.PresignPutObject(context.Background(), &s3.PutObjectInput{Bucket: &bucket, Key: &key}, func(ppo *s3.PresignOptions) { ppo.Expires = lifetime })

	if err != nil {
		return nil, err
	}

	presignedURL, err := url.Parse(presignedReq.URL)

	if err != nil {
		return nil, err
	}
	return presignedURL, nil
}

func NewS3Presigner(s3Client *s3.Client) *S3Presigner {
	return &S3Presigner{
		s3: s3Client,
	}
}
