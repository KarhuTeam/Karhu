package models

import (
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
)

type DONodeCreateForm struct {
	Region         string `form:"region" json:"region" valid:"required"`
	Hostname       string `form:"hostname" json:"hostname" valid:"required"`
	Description    string `form:"description" json:"description" valid:"-"`
	InstanceType   string `form:"instance_type" json:"instance_type" valid:"required"`
	Backups        string `form:"backups" json:"backups" valid:"-"`
	IpV6           string `form:"ipv6" json:"ipv6" valid:"-"`
	PrivateNetwork string `form:"private_network" json:"private_network" valid:"-"`
}

// Validator for node creation
func (f DONodeCreateForm) Validate() *errors.Errors {
	return validator.Validate(&f)
}
