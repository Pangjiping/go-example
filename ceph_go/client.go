package ceph_go

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

var (
	CephConn *s3.S3
)

func init() {
	auth := aws.Auth{
		AccessKey: "",
		SecretKey: "",
	}

	region := aws.Region{
		Name:                 "default",
		EC2Endpoint:          "http://<ceph-rgw ip>:<ceph-rgw port>",
		S3Endpoint:           "http://<ceph-rgw ip>:<ceph-rgw port>",
		S3BucketEndpoint:     "",    // not needed by aws s3
		S3LocationConstraint: false, // true if this region requires a LocationConstraint declaration
		S3LowercaseBucket:    false, // true if the region requires bucket names to be lower case
		Sign:                 aws.SignV2,
	}

	CephConn = s3.New(auth, region)
}
