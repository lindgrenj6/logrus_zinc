package logrus_zinc

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// validate that we're implementing the interface at compile time
var _ = (logrus.Hook)(&LocalZincHook{})

// LocalZincHook is the struc that contains the fields + methods for
// implementing a logrus hook that talks to zinc.
//
// In order to use this - just add an instance pointing at your local zinc
// instance with the proper fields set.
type LocalZincHook struct {
	URL   string
	Index string

	Username string
	Password string
}

func (k *LocalZincHook) Fire(entry *logrus.Entry) error {
	if k.URL == "" {
		k.URL = "http://localhost:4080"
	}
	if k.Index == "" {
		k.Index = "default"
	}

	go func() {
		data, err := entry.String()
		if err != nil {
			log.Print(err)
			return
		}
		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%v/api/%v/_doc", k.URL, k.Index),
			strings.NewReader(data),
		)
		if err != nil {
			log.Print(err)
			return
		}
		req.SetBasicAuth(k.Username, k.Password)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Print(err)
			return
		}
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		resp.Body.Close()
	}()

	return nil
}

func (k *LocalZincHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func FromEnv() (*LocalZincHook, error) {
	if os.Getenv("ZINC_SEARCH_USERNAME") == "" {
		return nil, fmt.Errorf("ZINC_SEARCH_USERNAME is required")
	}
	if os.Getenv("ZINC_SEARCH_PASSWORD") == "" {
		return nil, fmt.Errorf("ZINC_SEARCH_PASSWORD is required")
	}

	return &LocalZincHook{
		URL:      os.Getenv("ZINC_SEARCH_URL"),
		Index:    os.Getenv("ZINC_SEARCH_INDEX"),
		Username: os.Getenv("ZINC_SEARCH_USERNAME"),
		Password: os.Getenv("ZINC_SEARCH_PASSWORD"),
	}, nil
}
