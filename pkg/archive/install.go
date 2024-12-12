package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/uwu-tools/magex/pkg/downloads"
	"github.com/uwu-tools/magex/xplat"
)

// DownloadArchiveOptions are the set of options available for DownloadToGopathBin.
type DownloadArchiveOptions struct {
	downloads.DownloadOptions

	// ArchiveExtensions maps from the GOOS to the expected extension. Required.
	// For example, windows may use .zip while darwin/linux uses .tgz.
	ArchiveExtensions map[string]string

	// TargetFileTemplate specifies the path to the target binary in the archive. Required.
	// Supports the same templating as downloads.DownloadOptions.UrlTemplate.
	TargetFileTemplate string
}

// DownloadToGopathBin downloads an archived file to GOPATH/bin.
func DownloadToGopathBin(opts DownloadArchiveOptions) error {
	// determine the appropriate file extension based on the OS, e.g. windows gets .zip, otherwise .tgz
	opts.Ext = opts.ArchiveExtensions[runtime.GOOS]
	if opts.Ext == "" {
		return fmt.Errorf("no archive file extension was specified for the current GOOS (%s)", runtime.GOOS)
	}

	if opts.Hook == nil {
		opts.Hook = ExtractBinaryFromArchiveHook(opts)
	}

	return downloads.DownloadToGopathBin(opts.DownloadOptions)
}

// ExtractBinaryFromArchiveHook is the default hook for DownloadToGopathBin.
func ExtractBinaryFromArchiveHook(opts DownloadArchiveOptions) downloads.PostDownloadHook {
	return func(archiveFile string) (binPath string, err error) {
		// Save the binary next to the archive file in the temp directory
		outDir := filepath.Dir(archiveFile)

		// Render the name of the file in the archive
		opts.Ext = xplat.FileExt()
		targetFile, err := downloads.RenderTemplate(opts.TargetFileTemplate, opts.DownloadOptions)
		if err != nil {
			return "", fmt.Errorf("error rendering TargetFileTemplate %q with data %#v: %w", opts.TargetFileTemplate, opts.DownloadOptions, err)
		}

		log.Printf("extracting %s from %s...\n", targetFile, archiveFile)
		destinationFilePath := filepath.Join(outDir, targetFile)
		if err = os.MkdirAll(filepath.Dir(destinationFilePath), 0700); err != nil {
			return "", err
		}

		if strings.HasSuffix(archiveFile, ".zip") {
			if err = extractZip(archiveFile, targetFile, destinationFilePath); err != nil {
				return "", err
			}
		} else if strings.HasSuffix(archiveFile, ".tar.gz") || strings.HasSuffix(archiveFile, ".tgz") {
			if err = extractTar(archiveFile, targetFile, destinationFilePath); err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("unable to determine archive type of file %s", archiveFile)
		}

		// Extract the binary
		// The extracted file may be nested depending on its position in the archive
		binFile := filepath.Join(outDir, targetFile)

		// Check that file was extracted, Extract doesn't error out if you give it a missing targetFile
		if _, err := os.Stat(binFile); os.IsNotExist(err) {
			return "", fmt.Errorf("could not find %s in the archive", targetFile)
		}

		return binFile, nil
	}
}

func extractZip(archiveFile string, targetFile string, destinationFilePath string) error {
	archive, err := zip.OpenReader(archiveFile)
	if err != nil {
		return fmt.Errorf("unable to open %s for reading: %w", archiveFile, err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		if f.Name == targetFile {

			dstFile, err := os.OpenFile(destinationFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("unable to open %s for writing: %w", destinationFilePath, err)
			}

			fileInArchive, err := f.Open()
			if err != nil {
				return err
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				return fmt.Errorf("unable to unpack %s: %w", archiveFile, err)
			}

			dstFile.Close()
			fileInArchive.Close()
		}
	}

	return nil
}

func extractTar(archiveFile string, targetFile string, destinationFilePath string) error {
	gzipStream, err := os.Open(archiveFile)
	if err != nil {
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Name == targetFile {
			dstFile, err := os.Create(destinationFilePath)
			if err != nil {
				return fmt.Errorf("unable to open %s for writing: %w", destinationFilePath, err)
			}

			_, err = io.Copy(dstFile, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
