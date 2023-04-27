package files

import (
	"context"
	"fmt"
	"io"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sokool/domain"
)

type s3 struct {
	bucket     string
	context    context.Context
	connection *minio.Client
}

func newS3(u domain.URL) (*s3, error) {
	var c *minio.Client
	var err error

	if u.Schema != "s3" {
		return nil, Err("s3 schema required in url")
	}
	if u.Host == "" {
		return nil, Err("s3 hostname required")
	}
	if len(u.Path) == 0 {
		return nil, Err("bucket name required in first element of url path")
	}
	if c, err = minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(u.Username, u.Password, ""),
		Secure: true,
		Region: u.Query.Get("region"),
	}); err != nil {
		return nil, Err("%w", err)
	}

	return &s3{bucket: u.Path[0], connection: c, context: context.TODO()}, nil
}

func (c *s3) Write(l Location, from io.Reader) error {
	if _, err := c.connection.PutObject(c.context, c.bucket, l.Cut(1), from, -1, minio.PutObjectOptions{}); err != nil {
		return Err("%w", err)
	}
	return nil
}

func (c *s3) Read(l Location, to io.Writer) error {
	o, err := c.connection.GetObject(c.context, c.bucket, l.Cut(1), minio.GetObjectOptions{})
	if err != nil {
		return Err("%w", err)
	}
	if _, err = io.Copy(to, o); err != nil {
		return Err("%w", err)
	}
	return nil
}

func (c *s3) Files(l Location, recursive ...bool) ([]string, error) {
	p := l.Cut(1)
	o := c.connection.ListObjects(c.context, c.bucket, minio.ListObjectsOptions{
		WithVersions: true,
		WithMetadata: true,
		Prefix:       p,
		Recursive:    len(recursive) != 0 && recursive[0],
		MaxKeys:      0,
		StartAfter:   p,
		UseV1:        false,
	})

	var s []string
	for n := range o {
		if n.Err != nil {
			return s, Err("%w", n.Err)
		}
		s = append(s, fmt.Sprintf("%s %s %d %s", n.LastModified, n.Key, n.Size, n.ContentType))
	}
	return s, nil
}
