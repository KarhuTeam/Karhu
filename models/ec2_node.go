package models

import (
	"github.com/gotoolz/errors"
	"github.com/gotoolz/validator"
)

type EC2NodeCreateForm struct {
	AvailabilityZone string `form:"availability_zone" json:"availability_zone" valid:"required"`
	VPC              string `form:"vpc" json:"vpc" valid:"required"`
	SecurityGroup    string `form:"security_group" json:"security_group" valid:"required"`
	Hostname         string `form:"hostname" json:"hostname" valid:"required"`
	Description      string `form:"description" json:"description" valid:"-"`
	InstanceType     string `form:"instance_type" json:"instance_type" valid:"required"`
	// RootSize         string `form:"root_size" json:"root_size" valid:"int,required"`
	Monitoring string `form:"monitoring" json:"monitoring" valid:"-"`
}

// Validator for node creation
func (f EC2NodeCreateForm) Validate() *errors.Errors {
	return validator.Validate(&f)
}
