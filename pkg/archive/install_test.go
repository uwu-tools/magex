package archive

import (
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/magefile/mage/mg"
	"github.com/stretchr/testify/require"
	"github.com/uwu-tools/magex/pkg/downloads"
	"github.com/uwu-tools/magex/pkg/gopath"
	"github.com/uwu-tools/magex/xplat"
)

func TestDownloadArchiveToGopathBin_OsReplacement(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skipf("skipping test since gh only has binaries for amd64")
	}

	os.Setenv(mg.VerboseEnv, "true")
	err, cleanup := gopath.UseTempGopath()
	require.NoError(t, err, "Failed to set up a temporary GOPATH")
	defer cleanup()

	// gh cli unfortunately uses a different archive schema depending on the OS
	tmpl := "gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}/bin/gh{{.EXT}}"
	if runtime.GOOS == "windows" {
		tmpl = "bin/gh.exe"
	}

	opts := DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://github.com/cli/cli/releases/download/v{{.VERSION}}/gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}{{.EXT}}",
			Name:        "gh",
			Version:     "1.8.1",
			OsReplacement: map[string]string{
				"darwin": "macOS",
			},
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: tmpl,
	}

	err = DownloadToGopathBin(opts)
	require.NoError(t, err)

	_, err = exec.LookPath("gh" + xplat.FileExt())
	require.NoError(t, err)
}

func TestDownloadArchiveToGopathBin(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	err, cleanup := gopath.UseTempGopath()
	require.NoError(t, err, "Failed to set up a temporary GOPATH")
	defer cleanup()

	opts := DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://get.helm.sh/helm-v{{.VERSION}}-{{.GOOS}}-{{.GOARCH}}{{.EXT}}",
			Name:        "helm",
			Version:     "3.15.1",
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: "{{.GOOS}}-{{.GOARCH}}/helm{{.EXT}}",
	}

	err = DownloadToGopathBin(opts)
	require.NoError(t, err)

	_, err = exec.LookPath("helm" + xplat.FileExt())
	require.NoError(t, err)
}
