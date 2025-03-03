// 1. init download
// 2. get assets *.so

package pluginmanager

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
)

const (
	defaultRetryCount = 3
	defaultReleaseUrl = "https://github.com/merico-dev/stream/releases/download"
)

type DownloadClient struct {
	*resty.Client
}

func NewDownloadClient() *DownloadClient {
	dClient := DownloadClient{}
	dClient.Client = resty.New()
	dClient.SetRetryCount(defaultRetryCount)
	return &dClient
}

func (dc *DownloadClient) download(pluginsDir, pluginFilename, version string) error {
	dc.SetOutputDirectory(pluginsDir)

	downloadURL := fmt.Sprintf("%s/v%s/%s", defaultReleaseUrl, version, pluginFilename)
	tmpName := pluginFilename + ".tmp"

	response, err := dc.R().
		SetOutput(tmpName).
		SetHeader("Accept", "application/octet-stream").
		Get(downloadURL)
	if err != nil {
		log.Print(err)
		return err
	}
	if response.StatusCode() != http.StatusOK {
		if err = os.Remove(filepath.Join(pluginsDir, tmpName)); err != nil {
			return err
		}
		err = fmt.Errorf("downloading plugin %s from %s status code %d", pluginFilename, downloadURL, response.StatusCode())
		log.Print(err)
		return err
	}

	// rename, tmp file to real file
	err = os.Rename(
		filepath.Join(pluginsDir, tmpName),
		filepath.Join(pluginsDir, pluginFilename))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
