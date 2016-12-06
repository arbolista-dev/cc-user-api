package services

import (
	"os"
	"bytes"
	"github.com/arbolista-dev/cc-user-api/app/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	awsBucket, awsRegion, awsEndpoint, awsAkid, awsSak, profileDirectory string
)

func init() {
	awsBucket = os.Getenv("AWS_BUCKETNAME")
	awsRegion = os.Getenv("AWS_REGION")
	awsAkid = os.Getenv("AWS_ACCESS_KEY_ID")
	awsSak = os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsEndpoint = os.Getenv("AWS_ENDPOINT")
	profileDirectory = "/profile-photos/"
}

func UploadFile(file []byte, format string) (location string, err error) {

	f := bytes.NewReader(file)
	uuid := utils.RandString(20)

	path := profileDirectory + uuid + "." + format

	creds := credentials.NewStaticCredentials(awsAkid, awsSak, "")
	sess := session.New(&aws.Config{
	    Credentials: creds,
	    Region:      aws.String(awsRegion),
			Endpoint: 	 aws.String(awsEndpoint),
	})

	uploader := s3manager.NewUploader(sess)
	uploadResult, err := uploader.Upload(&s3manager.UploadInput{
	    Bucket: aws.String(awsBucket),
	    Key:    aws.String(path),
	    Body:   f,
	})

	location = uploadResult.Location
	if err != nil {
		return
	}

	return
}
