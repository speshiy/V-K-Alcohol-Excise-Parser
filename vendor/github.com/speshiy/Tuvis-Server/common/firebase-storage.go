package common

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

//FirebaseStorageBucket value
const FirebaseStorageBucket = "tuvisworld.appspot.com"

//ClientToFirebaseStorage connect to Firebase Storage
func ClientToFirebaseStorage() (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

//UploadFileToFirebaseStorageBucket upload file to bucket
func UploadFileToFirebaseStorageBucket(name string, r *io.Reader) (string, error) {
	var err error

	client, err := ClientToFirebaseStorage()
	if err != nil {
		return "", err
	}
	defer client.Close()

	bucket := client.Bucket(FirebaseStorageBucket)
	if err != nil {
		return "", err
	}

	//Get context of application
	ctx := context.Background()

	//Get object from storage bucket
	obj := bucket.Object(name)

	//Create writer for obj
	wc := obj.NewWriter(ctx)

	//Set public access to file
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	wc.CacheControl = "public, max-age=86400"

	if _, err = io.Copy(wc, *r); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	return objectURL(name), nil
}

func objectURL(name string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", FirebaseStorageBucket, name)
}
