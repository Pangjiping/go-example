package ceph_go

import (
	"gopkg.in/amz.v1/s3"
	"io/ioutil"
	"log"
)

// GetCephBucket 获取一个桶
func GetCephBucket(bucket string) *s3.Bucket {
	return CephConn.Bucket(bucket)
}

// Put2Bucket 将本地文件上传到ceph的一个bucket中
func Put2Bucket(bucket *s3.Bucket, localPath, cephPath string) (*s3.Bucket, error) {
	err := bucket.PutBucket(s3.PublicRead)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	bytes, err := ioutil.ReadFile(localPath)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	err = bucket.Put(cephPath, bytes, "octet-stream", s3.PublicRead)
	return bucket, err
}

// DownloadFromCeph 从ceph下载文件
func DownloadFromCeph(bucket *s3.Bucket, localPath, cephPath string) error {
	data, err := bucket.Get(cephPath)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return ioutil.WriteFile(localPath, data, 0666)
}

// DelCephData 删除指定的文件
func DelCephData(bucket *s3.Bucket, cephPath string) error {
	err := bucket.Del(cephPath)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

// DelBucket 删除桶
// 删除桶时要保证桶内文件已经被删除
func DelBucket(bucket *s3.Bucket) error {
	err := bucket.DelBucket()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func GetBatchFromCeph(bucket *s3.Bucket, prefixCephPath string) []string {
	maxBatch := 100

	// bucket.List() 返回桶内objects的信息，默认1000条
	resultListResp, err := bucket.List(prefixCephPath, "", "", maxBatch)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	keyList := make([]string, 0)
	for _, key := range resultListResp.Contents {
		keyList = append(keyList, key.Key)
	}

	return keyList

}
