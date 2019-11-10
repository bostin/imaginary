package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
)

const imageSourceProtocol string = "s3"
const ImageSourceAWSS3 ImageSourceType = "s3"

type AWSS3ImageSource struct {
	Config *SourceConfig
	Body []byte
}

func NewAWSS3ImageSource(config *SourceConfig) ImageSource {
	return &AWSS3ImageSource{config, nil}
}

func (s *AWSS3ImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodGet && s.getAWSS3Param(r) != ""
}

func (s *AWSS3ImageSource) GetImage(r *http.Request) ([]byte, error) {
	path := s.getAWSS3Param(r)
	if path == "" {
		return nil, ErrMissingParamS3
	}
	return s.fetchImage(r)
}

func (s *AWSS3ImageSource) getAWSS3Param(r *http.Request) string  {
	return r.URL.Query().Get(imageSourceProtocol)
}

func (s *AWSS3ImageSource) fetchImage(r *http.Request) ([]byte, error) {
	downloader := newAWSS3Downloader(s.Config, r)
	_, err := downloader.Download(s, &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key: aws.String("key"),
	})
	if err != nil {
		return nil, err
	}
	return s.Body, nil
}

func (s *AWSS3ImageSource) WriteAt(p []byte, off int64) (n int, err error) {
	s.Body = p
	return len(p), nil
}

func newAWSS3Downloader(config *SourceConfig, r *http.Request) *s3manager.Downloader {
	sess := session.Must(session.NewSession())
	return s3manager.NewDownloader(sess)
}

func init() {
	RegisterSource(ImageSourceAWSS3, NewAWSS3ImageSource)
}
