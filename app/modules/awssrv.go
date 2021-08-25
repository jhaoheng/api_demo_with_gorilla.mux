package modules

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var AWSSrv *AWSSrvObj

type AWSSrvObj struct {
	S3  *s3.S3
	SSM *ssm.SSM
}

func NewAWSSrv() *AWSSrvObj {
	sess := session.Must(session.NewSession())
	AWSSrv = &AWSSrvObj{
		S3:  s3.New(sess),
		SSM: ssm.New(sess),
	}
	return AWSSrv
}

func (srv *AWSSrvObj) S3_ListBuckets() (*s3.ListBucketsOutput, error) {
	input := &s3.ListBucketsInput{}
	return srv.S3.ListBuckets(input)
}

func (srv *AWSSrvObj) SSM_GetParameter(name string) (*ssm.GetParameterOutput, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String(name),
	}
	return srv.SSM.GetParameter(input)
}

func (srv *AWSSrvObj) SSM_GetParametersByPath(path string) (*ssm.GetParametersByPathOutput, error) {
	input := &ssm.GetParametersByPathInput{
		Path:      aws.String(path),
		Recursive: aws.Bool(true),
	}
	return srv.SSM.GetParametersByPath(input)
}
