package docker

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBundleDir = "bundle"
)

type TarHandler struct {
	TmpDirName string
}

func NewTarHandler() (*TarHandler, error) {
	dirName, err := os.MkdirTemp(os.TempDir(), "sim-cli")
	if err != nil {
		return nil, fmt.Errorf("error creating temp dir: %v", err)
	}
	return &TarHandler{
		TmpDirName: dirName,
	}, nil
}

func (t *TarHandler) Cleanup() error {
	return os.RemoveAll(t.TmpDirName)
}

// UnzipSupportBundle will unzip bundle into memory FS
// which can then be used to generate a tar ball for providing a context
// to build image with bundle packaged into support-bundle-kit base image
func (t *TarHandler) UnzipSupportBundle(bundleZipFile string) (err error) {

	// ensure destination exists
	destination := t.TmpDirName

	r, err := zip.OpenReader(bundleZipFile)
	if err != nil {
		return err
	}

	for _, f := range r.File {
		destPath := filepath.Join(destination, f.Name)
		if !strings.HasPrefix(destPath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid dest path %s", destPath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
				return err
			}

			destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_CREATE, f.Mode())
			if err != nil {
				return err
			}

			zFile, err := f.Open()
			if err != nil {
				return err
			}

			if _, err = io.Copy(destFile, zFile); err != nil {
				return err
			}
			zFile.Close()
			destFile.Close()
		}

	}

	// rename support bundle to ensure consistent tar file packaging
	baseFilePath := filepath.Base(bundleZipFile)
	baseDirName := strings.Split(baseFilePath, ".zip")
	return os.Rename(filepath.Join(destination, baseDirName[0]), filepath.Join(destination, defaultBundleDir))
}

// GenerateBundleTar attempts to parse FS/bundle to build a tar which can be passed
// as context for image creation step
func (t *TarHandler) GenerateBundleTar() (*bytes.Buffer, error) {
	fs := os.DirFS(t.TmpDirName)
	buf := &bytes.Buffer{}
	contextTar := tar.NewWriter(buf)
	if err := contextTar.AddFS(fs); err != nil {
		return nil, fmt.Errorf("error adding tmp directory %s to tar: %v", t.TmpDirName, err)
	}

	if err := contextTar.Close(); err != nil {
		return nil, fmt.Errorf("error closing tar file %v", err)
	}
	return buf, nil
}

func (t *TarHandler) AddDockerFile(baseImage string) error {
	dockerFile, err := generateTemplate(baseImage)
	if err != nil {
		return fmt.Errorf("error generating dockerfile from embedded template")
	}

	return os.WriteFile(filepath.Join(t.TmpDirName, "Dockerfile"), dockerFile.Bytes(), 0700)
}

// ReadTar is a helper utility to read contents of in memory tar file
// and is mostly used for debugging and testing
func ReadTar(buf *bytes.Buffer) error {
	tr := tar.NewReader(buf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
	}
	return nil
}

// BuildContextTar is a wrapper function tht builds a tar ball with Dockerfile and contents of bundle
// and this can be passed to image builder to ensure support bundle kit image is layered with
// actual support bundle contents to allow for subsequent processing by simulator
func BuildContextTar(bundlePath string, baseImage string) (*bytes.Buffer, error) {
	t, err := NewTarHandler()
	if err != nil {
		return nil, err
	}

	// prepare zip file and extract it into bundle folder
	if err := t.UnzipSupportBundle(bundlePath); err != nil {
		return nil, err
	}

	// add Dockerfile to root of tar image
	if err := t.AddDockerFile(baseImage); err != nil {
		return nil, err
	}

	buf, err := t.GenerateBundleTar()
	if err != nil {
		return nil, err
	}

	if err := t.Cleanup(); err != nil {
		return nil, fmt.Errorf("error cleaning up tmp directory: %v", err)
	}
	return buf, err
}

func generateTemplate(baseImage string) (bytes.Buffer, error) {
	contents := struct {
		BaseImage string
	}{
		BaseImage: baseImage,
	}

	dockerFile := `FROM {{ .BaseImage }}
EXPOSE 6443/tcp
COPY bundle /bundle
`
	dockerTemplate := template.Must(template.New("dockerfile").Parse(dockerFile))
	var output bytes.Buffer
	err := dockerTemplate.Execute(&output, contents)
	return output, err
}
