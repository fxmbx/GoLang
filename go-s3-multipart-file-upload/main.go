package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *s3.S3
)

const BUCKET_NAME = "funbis3bucket1"
const REGION = "eu-west-2"
const AWS_S3_KEY = ""
const AWS_S3_SECRET = ""

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
		Credentials: credentials.NewStaticCredentials(
			AWS_S3_KEY, AWS_S3_SECRET, ""),
	})))
}

//returns a list of all the list
func listBuckets() (resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func createBucket() *s3.CreateBucketOutput {
	res, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(BUCKET_NAME),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("bucket name already in use")
				panic(err)
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("you own this bucket dummy 🦆")
			default:
				fmt.Println("default")
				panic(err)
			}
		}
	}
	return res

}

func uploadObject(filename string) *s3.PutObjectOutput {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Uploading: ", filename)
	resp, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(strings.Split(filename, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
		// ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func listObjects() *s3.ListObjectsV2Output {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKET_NAME),
	})
	if err != nil {
		panic(err)
	}
	return resp
}

func getObject(filename string) {
	fmt.Println("downloading: ", filename)
	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}

}

func deleteObject(filename string) *s3.DeleteObjectOutput {
	fmt.Println("deleting : ", filename)
	res, err := s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}
	return res

}
func main() {
	fmt.Println(listBuckets())
	fmt.Println(createBucket())
	uploadObject("files/gopher_logo")
	fmt.Println(listObjects())
	getObject("gopher_logo")

}
