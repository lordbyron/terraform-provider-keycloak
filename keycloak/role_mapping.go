package keycloak

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
)

// Not a real object in keycloak, just convenient
type RoleMapping struct {
	Realm     string
	RoleId    string
	roleMemo  *Role
	UserName  string
	GroupName string
	UserId    string
	GroupId   string
	ClientId  string
}

const (
	roleMapBaseUri = "%s/auth/admin/realms/%s/%s/%s/role-mappings/%s"
)

func (rm *RoleMapping) Validate(c *KeycloakClient) error {
	if rm.UserName == "" && rm.GroupName == "" && rm.UserId == "" && rm.GroupId == "" {
		return errors.New("Must specify one of user, group, user_id, or group_id")
	}
	if rm.UserId != "" || rm.GroupId != "" {
		return nil
	}
	// must have specified group or user, need to get underlying id
	if rm.UserName != "" {
		u, err := c.GetUserByName(rm.UserName, rm.Realm)
		if err != nil {
			return err
		}
		rm.UserId = u.Id
		return nil
	}
	if rm.GroupName != "" {
		g, err := c.GetGroupByName(rm.GroupName, rm.Realm)
		if err != nil {
			return err
		}
		rm.GroupId = g.Id
		return nil
	}
	// shouldn't get here. maybe multiple configs were set
	return errors.New("Can only set one of user, group, user_id, group_id")
}

func DeserializeRoleMapping(str string) (*RoleMapping, error) {
	rm := &RoleMapping{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&rm)
	if err != nil {
		return nil, err
	}
	return rm, nil
}

func (rm *RoleMapping) Serialize() string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	clone := *rm
	clone.roleMemo = nil
	e.Encode(clone)
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (rm *RoleMapping) role(c *KeycloakClient) (*Role, error) {
	if rm.roleMemo != nil {
		return rm.roleMemo, nil
	}
	role, err := c.GetRole(rm.RoleId, rm.Realm)
	rm.roleMemo = role
	return role, err
}

func (rm *RoleMapping) roleMapUrl(base, suffix string) string {
	ug := "users"
	ugId := rm.UserId
	if rm.UserId == "" {
		ug = "groups"
		ugId = rm.GroupId
	}
	if rm.ClientId == "" {
		suffix = fmt.Sprintf("realm%s", suffix)
	} else {
		suffix = fmt.Sprintf("clients/%s%s", rm.ClientId, suffix)
	}
	return fmt.Sprintf(roleMapBaseUri, base, rm.Realm, ug, ugId, suffix)
}

func (rm *RoleMapping) availableUrl(base string) string {
	return rm.roleMapUrl(base, "/available")
}

func (rm *RoleMapping) compositeUrl(base string) string {
	return rm.roleMapUrl(base, "/composite")
}

func (rm *RoleMapping) baseUrl(base string) string {
	return rm.roleMapUrl(base, "")
}

/** API client methods **/
func (c *KeycloakClient) GetAvailableRoles(rm RoleMapping) ([]Role, error) {
	url := rm.availableUrl(c.url)
	var roles []Role
	err := c.get(url, &roles)
	return roles, err
}

func (c *KeycloakClient) GetCompositeRoles(rm RoleMapping) ([]Role, error) {
	url := rm.compositeUrl(c.url)
	var roles []Role
	err := c.get(url, &roles)
	return roles, err
}

func (c *KeycloakClient) AddRoleMapping(rm RoleMapping) error {
	url := rm.baseUrl(c.url)
	role, err := rm.role(c)
	if err != nil {
		return err
	}
	body := []Role{*role}
	_, err = c.post(url, &body)
	return err
}

func (c *KeycloakClient) DeleteRoleMapping(rm RoleMapping) error {
	url := rm.baseUrl(c.url)
	role, err := rm.role(c)
	if err != nil {
		return err
	}
	body := []Role{*role}
	err = c.delete(url, &body)
	return err
}
