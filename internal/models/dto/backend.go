package dto

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"strings"
)

type AddBackendRequest struct {
	ServerURL string `json:"server_url" valid:"required,url"`
}

type DeleteBackendRequest struct {
	ServerURL string `json:"server_url" valid:"required,url"`
}

type SetAlgorithmRequest struct {
	Algorithm string `json:"algorithm"`
}

func (r *SetAlgorithmRequest) Validate() error {
	if r.Algorithm != "round-robin" && r.Algorithm != "least-connections" {
		return errors.New("invalid algorithm, must be 'round-robin' or 'least-connections'")
	}
	return nil
}

func (a *AddBackendRequest) Validate() error {
	a.ServerURL = strings.TrimRight(a.ServerURL, "/")
	a.ServerURL = strings.Replace(a.ServerURL, "localhost", "127.0.0.1", 1)
	_, err := govalidator.ValidateStruct(a)
	return err
}

func (d *DeleteBackendRequest) Validate() error {
	d.ServerURL = strings.TrimRight(d.ServerURL, "/")
	d.ServerURL = strings.Replace(d.ServerURL, "localhost", "127.0.0.1", 1)
	_, err := govalidator.ValidateStruct(d)
	return err
}
