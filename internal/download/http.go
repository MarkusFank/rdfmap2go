package download

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DownloadFile(uri string, target io.Writer) error {

	if target == nil {
		return errors.New("target must not be nil")
	}

	response, err := http.Get(uri)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to download file from URI '%s'. Status %d", uri, response.StatusCode)
	}

	_, err = io.Copy(target, response.Body)

	if err != nil {
		return err
	}

	return nil
}
