package models

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "fmt"
    "mime/multipart"
    "bytes"
    "net/http"
    "path/filepath"
)

type AssetsDatastore interface {
  ImageStore(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type DA struct {
  *s3.S3
}

func InitAssetsDatastore(aws_access_key_id, aws_secret_access_key, aws_session_token string) (*DA, error) {
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
    return nil, err
  }
  
  fmt.Println("Buckets:")
  
  for _, b := range result.Buckets {
    fmt.Printf("* %s created on :%s\n",
    aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
  }

  fmt.Println("Returning the thing")
  return &DA{svc}, nil
}

func (da *DA) ImageStore(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
    // get the file size and read
  // the file content into a buffer
  size := fileHeader.Size
  buffer := make([]byte, size)
  file.Read(buffer)

  // create a unique file name for the file
  tempFileName := "pictures/" + filepath.Ext(fileHeader.Filename)
	
  // config settings: this is where you choose the bucket,
  // filename, content-type and storage class of the file
  // you're uploading
  _, err := da.S3.PutObject(&s3.PutObjectInput{
     Bucket:               aws.String("goldney"),
     Key:                  aws.String(tempFileName),
     ACL:                  aws.String("public-read"),// could be private if you want it to be access by only authorized users
     Body:                 bytes.NewReader(buffer),
     ContentLength:        aws.Int64(int64(size)),
     ContentType:        aws.String(http.DetectContentType(buffer)),
     //ContentDisposition:   aws.String("attachment"),
     //ServerSideEncryption: aws.String("AES256"),
     //StorageClass:         aws.String("INTELLIGENT_TIERING"),
  })
  if err != nil {
     return "", err
  }

  return tempFileName, err
}
