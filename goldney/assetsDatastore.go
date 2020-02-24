package models

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "fmt"
)

type AssetsDatastore interface {
}

func InitAssetsDatastore() () {
  sess, err := session.NewSession(&aws.Config{
      Region: aws.String("us-west-2")},
  )
  svc := s3.New(sess)

  result, err := svc.ListBuckets(nil)
  if err != nil {
    panic(err)
  }
  
  fmt.Println("Buckets:")
  
  for _, b := range result.Buckets {
    fmt.Printf("* %s created on %s\n",
    aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
  }
}


