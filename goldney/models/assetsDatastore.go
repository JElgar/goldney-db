package models

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "fmt"
    errors "jameselgar.com/goldney/errors"
    "mime/multipart"
)

type AssetsDatastore interface {
  ImageStore(f multipart.File) (*errors.ApiError)
}

type DA struct {
  *s3.S3
}

func InitAssetsDatastore() (*DA, *errors.ApiError) {
  sess, err := session.NewSession(&aws.Config{
      Region: aws.String("us-west-2"),
    })
  svc := s3.New(sess)

  result, err := svc.ListBuckets(nil)
  if err != nil {
    return nil, &errors.ApiError{err, "Error listing buckets", 418}
  }
  
  fmt.Println("Buckets:")
  
  for _, b := range result.Buckets {
    fmt.Printf("* %s created on :%s\n",
    aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
  }

  return &DA{svc}, nil
}

func (da *DA) ImageStore(f multipart.File) (*errors.ApiError) {
  return nil
}
