// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package testcharms holds a corpus of charms
// for testing.
package testcharms

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/charm.v6"
	"gopkg.in/juju/charmrepo.v3/csclient"
	"gopkg.in/juju/charmrepo.v3/csclient/params"
	"gopkg.in/juju/charmrepo.v3/testing"

	"github.com/juju/juju/charmstore"
	jtesting "github.com/juju/juju/testing"
)

const defaultSeries = "quantal"

// Repo provides access to the test charm repository.
var Repo = testing.NewRepo("charm-repo", defaultSeries)

// RawCharmstoreServer encapsulates HTTP methods
// that are used for interacting with a charmstore
type RawCharmstoreServer interface {
	// Get retrieves data from path.
	Get(path string, data interface{}) error

	// Put sends raw data to the charm store at path. It emulated the HTTP PUT method.
	Put(path string, data []string) error
}

// MinimalCharmstoreClient represents the essential methods of the
// charm store that enable tests to pass.
//
// It has been extracted from a few places, primarily charmrepo.v3 &
// charmrepo.v3/csclient, See also the charmstore/client.go  csWrapper
//  interface.
type MinimalCharmstoreClient interface {
	// charm-related methods
	// TODO(tsm): replace Get with GetCharm

	// GetCharm retrieves a charm. Requesting a missing charm returns an error.
	Get(id *charm.URL) (charm.Charm, error)

	// UploadCharm stores a charm.Charm for later retrieval via Get
	UploadCharm(id *charm.URL, charmData charm.Charm) (*charm.URL, error)
	UploadCharmWithRevision(id *charm.URL, ch charm.Charm, promulgatedRevision int) error

	// UploadCharmMultiSeries
	// UploadCharmWithRevision
	// UploadCharmWithMeta
	Publish(id *charm.URL, channels []params.Channel, resources map[string]int) error

	// GetArchive(id *charm.URL) (r io.ReadCloser, eid *charm.URL, hash string, size int64, err error)

	// GetResource(id *charm.URL, name string, revision int) (result ResourceData, err error)
	ListResources(id *charm.URL) ([]params.Resource, error)
	UploadResource(id *charm.URL, name, path string, file io.ReaderAt, size int64, progress csclient.Progress) (revision int, err error)
	// ResourceMeta(id *charm.URL, name string, revision int) (params.Resource, error)
	// ResumeUploadResource(uploadId string, id *charm.URL, name, path string, file io.ReaderAt, size int64, progress Progress) (revision int, err error)

	// DockerResourceUploadInfo(id *charm.URL, resourceName string) (*params.DockerInfoResponse, error)
	// DockerResourceDownloadInfo(id *charm.URL, resourceName string) (*params.DockerInfoResponse, error)
	AddDockerResource(id *charm.URL, resourceName string, imageName, digest string) (revision int, err error)

	Latest(channel params.Channel, ids []*charm.URL, headers map[string][]string) ([]charmstore.CharmRevision, error)

	/// internal methods

	// WithChannel
	WithChannel(channel params.Channel) MinimalCharmstoreClient

	// StatsUpdate(req params.StatsUpdateRequest) error)

	// http-related
	//Get(path string, result interface{}) error
	//Put(path string, val interface{}) error

	// PutWithResponse // unnecessary?
	// DoWithResponse // unnecessary?

	//PutExtraInfo(id *charm.URL, info map[string]interface{}) error
	//PutCommonInfo(id *charm.URL, info map[string]interface{}) error

	// looks hard
	// Meta(id *charm.URL, result interface{}) (*charm.URL, error)

	UploadBundleWithRevision(id *charm.URL, b charm.Bundle, promulgatedRevision int) error
	Put(path string, data []string) error

	//Get(*charm.URL) (macaroon.Slice, error)
}

type CharmstoreOperations interface {
	RawCharmstoreServer

	Login() error
}

type CharmstoreClientWithBundles interface {
	MinimalCharmstoreClient

	GetBundle(id *charm.URL) (charm.Bundle, error)
	UploadBundle(id *charm.URL, b charm.Bundle) (*charm.URL, error)
}

// RepoForSeries returns a new charm repository for the specified series.
func RepoForSeries(series string) *testing.Repo {
	// TODO(ycliuhw): workaround - currently `quantal` is not exact series
	// (for example, here makes deploy charm at charm-repo/quantal/mysql --series precise possible )!
	if series != "kubernetes" {
		series = defaultSeries
	}
	return testing.NewRepo("charm-repo", series)
}

// UploadCharmWithMeta pushes a new charm to the charmstore.
// The uploaded charm takes the supplied charmURL with metadata.yaml and metrics.yaml
// to define the charm, rather than relying on the charm to exist on disk.
// This allows you to create charm definitions directly in yaml and have them uploaded
// here for us in tests.
//
// For convenience the charm is also made public
func UploadCharmWithMeta(c *gc.C, client MinimalCharmstoreClient, charmURL, meta, metrics string, revision int) (*charm.URL, charm.Charm) {
	ch := testing.NewCharm(c, testing.CharmSpec{
		Meta:     meta,
		Metrics:  metrics,
		Revision: revision,
	})
	chURL, err := client.UploadCharm(charm.MustParseURL(charmURL), ch)
	c.Assert(err, jc.ErrorIsNil)
	SetPublic(c, client, chURL)
	return chURL, ch
}

// UploadCharm sets default series to quantal
func UploadCharm(c *gc.C, client MinimalCharmstoreClient, url, name string) (*charm.URL, charm.Charm) {
	return UploadCharmWithSeries(c, client, url, name, defaultSeries)
}

// UploadCharmWithSeries uploads a charm using the given charm store client, and returns
// the resulting charm URL and charm.
//
// It also adds any required resources that haven't already been uploaded
// with the content "<resourcename> content".
func UploadCharmWithSeries(c *gc.C, client MinimalCharmstoreClient, url, name, series string) (*charm.URL, charm.Charm) {
	id := charm.MustParseURL(url)
	promulgatedRevision := -1
	if id.User == "" {
		// We still need a user even if we are uploading a promulgated charm.
		id.User = "who"
		promulgatedRevision = id.Revision
	}
	ch := RepoForSeries(series).CharmArchive(c.MkDir(), name)

	// Upload the charm.
	err := client.UploadCharmWithRevision(id, ch, promulgatedRevision)
	c.Assert(err, jc.ErrorIsNil)

	// Upload any resources required for publishing.
	var resources map[string]int
	if len(ch.Meta().Resources) > 0 {
		// The charm has resources.
		// Ensure that all the required resources are uploaded
		// before we publish.
		resources = make(map[string]int)
		current, err := client.WithChannel(params.UnpublishedChannel).ListResources(id)
		c.Assert(err, gc.IsNil)
		for _, r := range current {
			if r.Revision == -1 {
				// The resource doesn't exist so upload one.
				if r.Type == "oci-image" {
					_, err = client.AddDockerResource(id, r.Name, "Image", "sha")
					c.Assert(err, jc.ErrorIsNil)
				} else {
					content := r.Name + " content"
					_, err := client.UploadResource(id, r.Name, "", strings.NewReader(content), int64(len(content)), nil)
					c.Assert(err, jc.ErrorIsNil)
				}
				r.Revision = 0
			}
			resources[r.Name] = r.Revision
		}
	}

	SetPublicWithResources(c, client, id, resources)

	return id, ch
}

// UploadCharmMultiSeries uploads a charm with revision using the given charm store client,
// and returns the resulting charm URL and charm. This API caters for new multi-series charms
// which do not specify a series in the URL.
func UploadCharmMultiSeries(c *gc.C, client MinimalCharmstoreClient, url, name string) (*charm.URL, charm.Charm) {
	id := charm.MustParseURL(url)
	if id.User == "" {
		// We still need a user even if we are uploading a promulgated charm.
		id.User = "who"
	}
	ch := Repo.CharmArchive(c.MkDir(), name)

	// Upload the charm.
	curl, err := client.UploadCharm(id, ch)
	c.Assert(err, jc.ErrorIsNil)

	SetPublic(c, client, curl)

	// Return the charm and its URL.
	return curl, ch
}

// UploadBundle uploads a bundle using the given charm store client, and
// returns the resulting bundle URL and bundle.
func UploadBundle(c *gc.C, client MinimalCharmstoreClient, url, name string) (*charm.URL, charm.Bundle) {
	id := charm.MustParseURL(url)
	promulgatedRevision := -1
	if id.User == "" {
		// We still need a user even if we are uploading a promulgated bundle.
		id.User = "who"
		promulgatedRevision = id.Revision
	}
	b := Repo.BundleArchive(c.MkDir(), name)

	// Upload the bundle.
	err := client.UploadBundleWithRevision(id, b, promulgatedRevision)
	c.Assert(err, jc.ErrorIsNil)

	SetPublic(c, client, id)

	// Return the bundle and its URL.
	return id, b
}

// SetPublicWithResources sets the charm or bundle with the given id to be
// published with global read permissions to the stable channel.
//
// The named resources with their associated revision
// numbers are also published.
func SetPublicWithResources(c *gc.C, client MinimalCharmstoreClient, id *charm.URL, resources map[string]int) {
	// Publish to the stable channel.
	err := client.Publish(id, []params.Channel{params.StableChannel}, resources)
	c.Assert(err, jc.ErrorIsNil)

	// Allow stable read permissions to everyone.
	err = client.WithChannel(params.StableChannel).Put("/"+id.Path()+"/meta/perm/read", []string{params.Everyone})
	c.Assert(err, jc.ErrorIsNil)
}

// SetPublic sets the charm or bundle with the given id to be
// published with global read permissions to the stable channel.
func SetPublic(c *gc.C, client MinimalCharmstoreClient, id *charm.URL) {
	SetPublicWithResources(c, client, id, nil)
}

// CheckCharmReady ensures that a desired charm archive exists and
// has some content.
func CheckCharmReady(c *gc.C, charmArchive *charm.CharmArchive) {
	fileSize := func() int64 {
		f, err := os.Open(charmArchive.Path)
		c.Assert(err, jc.ErrorIsNil)
		defer f.Close()

		fi, err := f.Stat()
		c.Assert(err, jc.ErrorIsNil)
		return fi.Size()
	}

	var oldSize, currentSize int64
	var charmReady bool
	runs := 1
	timeout := time.After(jtesting.LongWait)
	for !charmReady {
		select {
		case <-time.After(jtesting.ShortWait):
			currentSize = fileSize()
			// Since we do not know when the charm is ready, for as long as the size changes
			// we'll assume that we'd need to wait.
			charmReady = oldSize != 0 && currentSize == oldSize
			c.Logf("%d: new file size %v (old size %v)", runs, currentSize, oldSize)
			oldSize = currentSize
			runs++
		case <-timeout:
			c.Fatalf("timed out waiting for charm @%v to be ready", charmArchive.Path)
		}
	}
}

// InjectFilesToCharmArchive overwrites the contents of pathToArchive with a
// new archive containing the original files plus the ones provided in the
// fileContents map (key: file name, value: file contents).
func InjectFilesToCharmArchive(pathToArchive string, fileContents map[string]string) error {
	zr, err := zip.OpenReader(pathToArchive)
	if err != nil {
		return err
	}
	defer func() { _ = zr.Close() }()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	defer func() { _ = zw.Close() }()

	// Copy existing files
	for _, f := range zr.File {
		w, err := zw.CreateHeader(&f.FileHeader)
		if err != nil {
			return err
		}

		r, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(w, r)
		_ = r.Close()
		if err != nil {
			return err
		}
	}

	// Add new files
	for name, contents := range fileContents {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}

		if _, err = w.Write([]byte(contents)); err != nil {
			return err
		}
	}

	// Overwrite original archive with the patched version
	_, _ = zr.Close(), zw.Close()
	return ioutil.WriteFile(pathToArchive, buf.Bytes(), 0644)
}
