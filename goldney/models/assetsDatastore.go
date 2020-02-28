package models

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/credentials"
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

func InitAssetsDatastore(aws_access_key_id, aws_secret_access_key, aws_session_token string) (*DA, *errors.ApiError) {
  fmt.Println("Credentials thing")
  fmt.Println(aws_access_key_id)
  fmt.Println(aws_secret_access_key)
  fmt.Println(aws_session_token)
  fmt.Println("End credentials thing")
  sess, err := session.NewSession(&aws.Config{
      Region: aws.String("us-east-1"),
      Credentials: credentials.NewStaticCredentials(
       aws_access_key_id,   // id
       aws_secret_access_key, // key
       aws_session_token),  // token can be left blank for now
  })
  svc := s3.New(sess)

  result, err := svc.ListBuckets(nil)
  if err != nil {
    panic(err)
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
