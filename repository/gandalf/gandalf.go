// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gandalf provides an implementation of the RepositoryManager, that
// uses Gandalf (https://github.com/tsuru/gandalf). This package doesn't expose
// any public types, in order to use it, users need to import the package and
// then configure tsuru to use the "gandalf" repo-manager.
//
//     import _ "github.com/tsuru/tsuru/repository/gandalf"
package gandalf

import (
	"github.com/tsuru/config"
	"github.com/tsuru/go-gandalfclient"
	"github.com/tsuru/tsuru/repository"
)

func init() {
	repository.Register("gandalf", gandalfManager{})
}

type gandalfManager struct{}

func (gandalfManager) client() (*gandalf.Client, error) {
	url, err := config.GetString("git:api-server")
	if err != nil {
		return nil, err
	}
	client := gandalf.Client{Endpoint: url}
	return &client, nil
}

func (m gandalfManager) CreateUser(username string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	_, err = client.NewUser(username, nil)
	return err
}

func (m gandalfManager) RemoveUser(username string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return client.RemoveUser(username)
}

func (m gandalfManager) CreateRepository(name string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	_, err = client.NewRepository(name, nil, true)
	return err
}

func (m gandalfManager) RemoveRepository(name string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return client.RemoveRepository(name)
}

func (m gandalfManager) GetRepository(name string) (repository.Repository, error) {
	client, err := m.client()
	if err != nil {
		return repository.Repository{}, err
	}
	repo, err := client.GetRepository(name)
	if err != nil {
		return repository.Repository{}, err
	}
	return repository.Repository{
		Name:         repo.Name,
		ReadOnlyURL:  repo.GitURL,
		ReadWriteURL: repo.SshURL,
	}, nil
}

func (m gandalfManager) GrantAccess(repository, user string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return client.GrantAccess([]string{repository}, []string{user})
}

func (m gandalfManager) RevokeAccess(repository, user string) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return client.RevokeAccess([]string{repository}, []string{user})
}

func (m gandalfManager) AddKey(username string, key repository.Key) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	keyMap := map[string]string{key.Name: key.Body}
	return client.AddKey(username, keyMap)
}

func (m gandalfManager) RemoveKey(username string, key repository.Key) error {
	client, err := m.client()
	if err != nil {
		return err
	}
	return client.RemoveKey(username, key.Name)
}

func (m gandalfManager) ListKeys(username string) ([]repository.Key, error) {
	client, err := m.client()
	if err != nil {
		return nil, err
	}
	keyMap, err := client.ListKeys(username)
	if err != nil {
		return nil, err
	}
	keys := make([]repository.Key, 0, len(keyMap))
	for name, body := range keyMap {
		keys = append(keys, repository.Key{Name: name, Body: body})
	}
	return keys, nil
}
