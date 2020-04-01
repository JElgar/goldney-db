package models

import (
    "github.com/aws/aws-sdk-go/aws"
    session "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "fmt"
    "mime/multipart"
    "bytes"
    "net/http"
    "path/filepath"
    "time"
)

type AssetsDatastore interface {
  ImageStore(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type DA struct {
  S3 *s3.S3
  Session *session.Session
}

// Used to set up s3 with inital credentials
func InitAssetsDatastore(aws_access_key_id, aws_secret_access_key, aws_session_token string) (*DA, error) {
  creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, aws_session_token)
  creds.Expire()
  // Retrieve the credentials value
  credValue, err := creds.Get()
  if err != nil {
    panic(err)
  }

  fmt.Println("Credentials thing")
  fmt.Println(aws_access_key_id)
  fmt.Println(aws_secret_access_key)
  fmt.Println(aws_session_token)
  fmt.Println("End credentials thing")
  
  fmt.Println("Credentials thing 2")
  fmt.Println(credValue.AccessKeyID)
  fmt.Println(credValue.SecretAccessKey)
  fmt.Println(credValue.SessionToken)
  fmt.Println("End credentials thing")
  
  //sess, err := session.NewSession(&aws.Config{
  //    Region: aws.String("us-east-1"),
  //    Credentials: credentials.NewStaticCredentials(
  //     aws_access_key_id,   // id
  //     aws_secret_access_key, // key
  //     aws_session_token),  // token can be left blank for now
  //})
  //svc := s3.New(sess)
  sess, err := session.NewSession(&aws.Config{
      Region: aws.String("us-east-1"),
      Credentials: creds  })
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
  return &DA{svc, sess}, nil
}

func (da *DA) updateS3Session () {
  creds := da.Session.Config.Credentials
  creds.Expire()
  // Retrieve the credentials value
  _ , err := creds.Get()
  if err != nil {
    panic(err)
  }
  sess, err := session.NewSession(&aws.Config{
      Region: aws.String("us-east-1"),
      Credentials: creds  })
  svc := s3.New(sess)
  da.Session = sess
  da.S3 = svc
}

func (da *DA) ImageStore(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
    // get the file size and read
  // the file content into a buffer
  size := fileHeader.Size
  buffer := make([]byte, size)
  file.Read(buffer)
  now := time.Now().String()

  // create a unique file name for the file
  tempFileName := "pictures/" + fileHeader.Filename  + now + filepath.Ext(fileHeader.Filename)
	
  // config settings: this is where you choose the bucket,
  // filename, content-type and storage class of the file
  // you're uploading
  _, err := da.S3.PutObject(&s3.PutObjectInput{
     Bucket:               aws.String("goldney"),
     Key:                  aws.String(tempFileName),
     ACL:                  aws.String("public-read"),// could be private if you want it to be access by only authorized users
     Body:                 bytes.NewReader(buffer),
     ContentLength:        aws.Int64(int64(size)),
     ContentType:          aws.String(http.DetectContentType(buffer)),
     //ContentDisposition:   aws.String("attachment"),
     //ServerSideEncryption: aws.String("AES256"),
     //StorageClass:         aws.String("INTELLIGENT_TIERING"),
  })
  if err != nil {
     return "", err
  }
  
  fmt.Println("You just uploaded: " + tempFileName)
  return tempFileName, err
}
