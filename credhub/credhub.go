package credhub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	chcli "code.cloudfoundry.org/credhub-cli/credhub"
)

type CredHub interface {
	FindAll() ([]*Credential, error)
	ReGenerate(*Credential) error
	UpdateTransitional(*Credential) error
}

func NewCredHub(ch *chcli.CredHub) CredHub {
	return &credhub{ch}
}

type credhub struct {
	client *chcli.CredHub
}

func (ch *credhub) FindAll() ([]*Credential, error) {
	// use struct to filter get uniuqe paths
	paths := make(map[string]struct{})

	// Note: If a certificate credential only has one version and it is
	// marked as transitional the credential name will not be returned by this endpoint.
	creds, err := ch.client.FindByPath("")
	if err != nil {
		return nil, err
	}

	for _, cred := range creds.Credentials {
		paths[cred.Name] = struct{}{}
	}

	certs, err := ch.client.GetAllCertificatesMetadata()
	if err != nil {
		return nil, err
	}

	for _, cert := range certs {
		paths[cert.Name] = struct{}{}
	}

	keys := make([]string, 0)
	for k := range paths {
		keys = append(keys, k)
	}

	return ch.getAllVersionsForAllPaths(keys)
}

func (ch *credhub) ReGenerate(*Credential) error {
	return nil
}

func (ch *credhub) UpdateTransitional(c *Credential) error {
	certMeta, err := ch.client.GetCertificateMetadataByName(c.Name)
	if err != nil {
		return fmt.Errorf("failed to get certificate meta for: %s got: %s", c.Name, err)
	}

	path := fmt.Sprintf("/api/v1/certificates/%s/update_transitional_version", certMeta.Id)
	body := map[string]interface{}{"version": c.ID}
	if c.Transitional {
		body["version"] = nil
	}
	resp, err := ch.client.Request(http.MethodPut, path, nil, body, true)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed request: %s with body: %s got: %s", path, body, err)
	}
	return nil
}

func (ch *credhub) getAllVersions(path string) ([]*Credential, error) {
	resp, err := ch.client.Request(http.MethodGet, "/api/v1/data",
		url.Values{"name": []string{path}}, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed request got: %s", err)
	}
	defer resp.Body.Close()

	result := struct {
		Data []*Credential `json:"data"`
	}{}

	return result.Data, json.NewDecoder(resp.Body).Decode(&result)
}

func (ch *credhub) getAllVersionsForAllPaths(paths []string) ([]*Credential, error) {
	pathChannel := make(chan string)
	errorChannel := make(chan error)
	resultChannel := make(chan *Credential)
	waitGroup := sync.WaitGroup{}

	for t := 0; t < 50; t++ {
		waitGroup.Add(1)
		go func(pc chan string, ec chan error, rc chan *Credential, wg *sync.WaitGroup) {
			for path := range pc {
				res, err := ch.getAllVersions(path)
				if err != nil {
					ec <- err
				}
				for _, cred := range res {
					rc <- cred
				}
			}
			wg.Done()
		}(pathChannel, errorChannel, resultChannel, &waitGroup)
	}

	go func() {
		for _, path := range paths {
			pathChannel <- path
		}
		close(pathChannel)
	}()

	go func() {
		waitGroup.Wait()
		close(resultChannel)
		close(errorChannel)
	}()

	results := make([]*Credential, 0)

	for {
		select {
		case result, ok := <-resultChannel:
			if !ok {
				return results, nil
			}
			results = append(results, result)

		case err, ok := <-errorChannel:
			if !ok {
				return results, nil
			}
			if err != nil {
				return nil, err
			}
		}
	}
}
