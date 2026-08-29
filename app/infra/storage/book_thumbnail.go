package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type bookThumbnailRepo struct {
	client *s3.S3
	bucket string
	prefix string
}

func NewBookThumbnailRepo(ctx context.Context, bucket string, prefix string, region string) (*bookThumbnailRepo, error) {
	_ = ctx
	if strings.TrimSpace(bucket) == "" {
		return nil, errors.New("S3 bucket is required")
	}
	cfg := aws.NewConfig()
	if strings.TrimSpace(region) != "" {
		cfg.Region = aws.String(region)
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}
	return &bookThumbnailRepo{
		client: s3.New(sess),
		bucket: bucket,
		prefix: strings.Trim(strings.TrimSpace(prefix), "/"),
	}, nil
}

func (r *bookThumbnailRepo) key(name string) string {
	if r.prefix == "" {
		return name
	}
	return path.Join(r.prefix, name)
}

func (r *bookThumbnailRepo) Save(ctx context.Context, name string, data []byte) error {
	_, err := r.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      &r.bucket,
		Key:         aws.String(r.key(name)),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentTypeByExt(name)),
	})
	return err
}

func (r *bookThumbnailRepo) Load(ctx context.Context, name string) ([]byte, error) {
	out, err := r.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &r.bucket,
		Key:    aws.String(r.key(name)),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey || awsErr.Code() == "NotFound" {
				return nil, errors.New("not found")
			}
		}
		return nil, err
	}
	defer out.Body.Close()
	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func contentTypeByExt(name string) string {
	switch strings.ToLower(path.Ext(name)) {
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
