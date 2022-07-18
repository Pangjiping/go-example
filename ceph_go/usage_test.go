package ceph_go

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CephUsage(t *testing.T) {
	bucketName := "bucket_test"
	fileName := "~/go/src/github.com/Pangjiping/go-example/ceph_go/test.text"
	cephPath := "/static/default/bucket_test/V1/" + "test_ceph.text"

	// 获取指定桶
	bucket := GetCephBucket(bucketName)
	assert.NotNil(t, bucket)

	// 上传文件
	bucket, err := Put2Bucket(bucket, fileName, cephPath)
	assert.Nil(t, err)
	assert.NotNil(t, bucket)

	// 下载文件
	localPath := "~/go/src/github.com/Pangjiping/go-example/ceph_go/download.text"
	err = DownloadFromCeph(bucket, localPath, cephPath)
	assert.Nil(t, err)

	// 获取下载url
	url := bucket.SignedURL(cephPath, time.Now().Add(time.Hour))
	assert.NotNil(t, url)
	fmt.Println(url)

	// 批量查找
	prefixCephPath := "static/default/bucket_test/V1"
	lists := GetBatchFromCeph(bucket, prefixCephPath)
	assert.NotNil(t, lists)
	for _, list := range lists {
		fmt.Println(list)
	}

	// 删除数据
	err = DelCephData(bucket, cephPath)
	assert.Nil(t, err)

	// 删除桶
	err = DelBucket(bucket)
	assert.Nil(t, err)
}
