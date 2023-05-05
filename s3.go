package files

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sokool/domain"
)

type s3 struct {
	bucket     string
	dir        Location
	context    context.Context
	connection *minio.Client
}

func newS3(u domain.URL) (*s3, error) {
	var s = s3{context: context.TODO(), dir: u.Path.Trim(1)}
	var err error

	if u.Schema != "s3" {
		return nil, Err("s3 schema required in url")
	}
	if u.Host == "" {
		return nil, Err("s3 hostname required")
	}
	if s.bucket = u.Path.Trim(0, 1).Replace("/", ""); u.Path.IsZero() {
		return nil, Err("s3 bucket name required in first element of url path")
	}

	if s.connection, err = minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(u.Username, u.Password, ""),
		Secure: true,
		Region: u.Query.Get("region"),
	}); err != nil {
		return nil, Err("can not establish s3 connection %w", err)
	}
	return &s, nil
}

func (c *s3) Write(l Location, from io.Reader, m ...Meta) error {
	if len(m) == 0 {
		m = append(m, nil)
	}
	p := c.dir.Append(l).Replace("/", "")
	r, err := c.connection.PutObject(c.context, c.bucket, p, from, -1, minio.PutObjectOptions{
		UserMetadata: m[0],
	})
	if err != nil {
		return Err("%w", err)
	}
	if m[0] != nil {
		m[0]["Location"] = r.Location
	}
	return nil
}

func (c *s3) Read(l Location, to io.Writer, m ...Meta) error {
	p := c.dir.Append(l).Replace("/", "")
	o, err := c.connection.GetObject(c.context, c.bucket, p, minio.GetObjectOptions{})
	if err != nil {
		return Err("%w", err)
	}
	if _, err = io.Copy(to, o); err != nil {
		return Err("%w", err)
	}
	if len(m) == 0 || m[0] == nil {
		return nil
	}
	s, err := o.Stat()
	if err != nil {
		return Err("%w", err)
	}
	if err = m[0].Merge(s.UserMetadata); err != nil {
		return Err("merging metadata failed %w", err)
	}
	return nil
}

func (c *s3) Files(l Location, recursive ...bool) ([]string, error) {
	p := l.Replace("/", "")
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
