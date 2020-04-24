package main

import (
	"github.com/minio/minio-go/v6"
	"io"
	"log"
	"strings"
)

type MinioObjectStore struct {
	Bucket string
	Client *minio.Client
}

func (store MinioObjectStore) GetReader(address ObjectAddress) (io.ReadCloser, error) {
	objectName := address.Route + "/" + address.Key
	return store.Client.GetObject(
		store.Bucket,
		objectName,
		minio.GetObjectOptions{},
	)
}

func (store MinioObjectStore) GetWriter(address ObjectAddress) (io.WriteCloser, error) {
	panic("Not implemented")
}


func getMetaMinio(info minio.ObjectInfo) ObjectMeta {
	etag := strings.ReplaceAll(info.ETag, "\"", "")

	return ObjectMeta{
		Size: info.Size,
		ETag: etag,
	}
}

func getAddressMinio(info minio.ObjectInfo) ObjectAddress {
	nonempty := func(s string) bool { return s != "" }
	segments := filter(strings.Split(info.Key, "/"), nonempty)

	return ObjectAddress{
		Key:   segments[len(segments)-1],
		Route: strings.Join(segments[:len(segments)-1], "/"),
	}
}

func getInfoMinio(info minio.ObjectInfo) ObjectInfo {
	return ObjectInfo{
		Meta:    getMetaMinio(info),
		Address: getAddressMinio(info),
	}
}

func (store MinioObjectStore) GetInfos() <-chan ObjectInfo {
	resultsCh := make(chan ObjectInfo)
	doneCh := make(chan struct{})

	go func() {
		for info := range store.Client.ListObjects(store.Bucket, "", true, doneCh) {
			if info.Err != nil {
				log.Println(info.Err)
			}
	
			objInfo := getInfoMinio(info)
	
			resultsCh <- objInfo
		}
		close(resultsCh)
		close(doneCh)
	}() 
	
	return resultsCh
}
