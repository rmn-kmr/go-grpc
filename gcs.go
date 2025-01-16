package lsp

import (
	gstorage "cloud.google.com/go/storage"
	"context"
)

func GetGCSObject(
	ctx context.Context,
	g gstorage.Client,
	fileName, bucketName string) (*gstorage.Reader, error) {
	bucket := g.Bucket(bucketName)
	object := bucket.Object(fileName)
	return object.NewReader(ctx)
}
