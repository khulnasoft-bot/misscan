package iam

import (
	misscanTypes "github.com/khulnasoft-lab/misscan/pkg/types"
	"github.com/liamg/iamgo"
)

type IAM struct {
	PasswordPolicy     PasswordPolicy
	Policies           []Policy
	Groups             []Group
	Users              []User
	Roles              []Role
	ServerCertificates []ServerCertificate
}

type ServerCertificate struct {
	Metadata   misscanTypes.Metadata
	Expiration misscanTypes.TimeValue
}

type Policy struct {
	Metadata misscanTypes.Metadata
	Name     misscanTypes.StringValue
	Document Document
	Builtin  misscanTypes.BoolValue
}

type Document struct {
	Metadata misscanTypes.Metadata
	Parsed   iamgo.Document
	IsOffset bool
	HasRefs  bool
}

func (d Document) ToRego() interface{} {
	m := d.Metadata
	doc, _ := d.Parsed.MarshalJSON()
	return map[string]interface{}{
		"filepath":  m.Range().GetFilename(),
		"startline": m.Range().GetStartLine(),
		"endline":   m.Range().GetEndLine(),
		"managed":   m.IsManaged(),
		"explicit":  m.IsExplicit(),
		"value":     string(doc),
		"fskey":     misscanTypes.CreateFSKey(m.Range().GetFS()),
	}
}

type Group struct {
	Metadata misscanTypes.Metadata
	Name     misscanTypes.StringValue
	Users    []User
	Policies []Policy
}

type User struct {
	Metadata   misscanTypes.Metadata
	Name       misscanTypes.StringValue
	Groups     []Group
	Policies   []Policy
	AccessKeys []AccessKey
	MFADevices []MFADevice
	LastAccess misscanTypes.TimeValue
}

func (u *User) HasLoggedIn() bool {
	return u.LastAccess.GetMetadata().IsResolvable() && !u.LastAccess.IsNever()
}

type MFADevice struct {
	Metadata  misscanTypes.Metadata
	IsVirtual misscanTypes.BoolValue
}

type AccessKey struct {
	Metadata     misscanTypes.Metadata
	AccessKeyId  misscanTypes.StringValue
	Active       misscanTypes.BoolValue
	CreationDate misscanTypes.TimeValue
	LastAccess   misscanTypes.TimeValue
}

type Role struct {
	Metadata misscanTypes.Metadata
	Name     misscanTypes.StringValue
	Policies []Policy
}

func (d Document) MetadataFromIamGo(r ...iamgo.Range) misscanTypes.Metadata {
	m := d.Metadata
	if d.HasRefs {
		return m
	}
	newRange := m.Range()
	var start int
	if !d.IsOffset {
		start = newRange.GetStartLine()
	}
	for _, rng := range r {
		newRange := misscanTypes.NewRange(
			newRange.GetLocalFilename(),
			start+rng.StartLine,
			start+rng.EndLine,
			newRange.GetSourcePrefix(),
			newRange.GetFS(),
		)
		m = misscanTypes.NewMetadata(newRange, m.Reference()).WithParent(m)
	}
	return m
}
