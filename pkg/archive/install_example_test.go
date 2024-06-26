package archive_test

import (
	"log"
	"testing"

	"github.com/uwu-tools/magex/pkg/archive"
	"github.com/uwu-tools/magex/pkg/downloads"
	"github.com/uwu-tools/magex/pkg/gopath"
)

func TestExampleDownloadToGopathBin(t *testing.T) {
	ExampleDownloadToGopathBin()
}

func ExampleDownloadToGopathBin() {
	opts := archive.DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://get.helm.sh/helm-{{.VERSION}}-{{.GOOS}}-{{.GOARCH}}{{.EXT}}",
			Name:        "helm",
			Version:     "v3.15.1",
		},
		ArchiveExtensions: map[string]string{
			"darwin":  ".tar.gz",
			"linux":   ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: "{{.GOOS}}-{{.GOARCH}}/helm{{.EXT}}",
	}
	err := archive.DownloadToGopathBin(opts)
	if err != nil {
		log.Fatal("could not download helm")
	}

	// Add GOPATH/bin to PATH if necessary so that we can immediately
	// use the installed tool
	gopath.EnsureGopathBin()
}
