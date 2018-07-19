package keycloak

import (
	"fmt"
	neturl "net/url"
)

// Does not implement composite roles at this timej
type Role struct {
	Id                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	ClientRole         bool   `json:"clientRole,omitempty"`
	ContainerId        string `json:"containerId,omitempty"`
	Description        string `json:"description,omitempty"`
	ScopeParamRequired bool   `json:"scopeParamRequired,omitempty"`
}

const (
	realmRoleCreateUri  = "%s/auth/admin/realms/%s/roles"
	clientRoleCreateUri = "%s/auth/admin/realms/%s/clients/%s/roles"
	roleUri             = "%s/auth/admin/realms/%s/roles-by-id/%s"
)

func (c *KeycloakClient) GetRole(id, realm string) (*Role, error) {
	url := fmt.Sprintf(roleUri, c.url, realm, id)

	var role Role
	err := c.get(url, &role)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (c *KeycloakClient) CreateRealmRole(role *Role, realm string) (*Role, error) {
	url := fmt.Sprintf(realmRoleCreateUri, c.url, realm)
	return c.createRole(role, url)
}

func (c *KeycloakClient) CreateClientRole(role *Role, realm, clientId string) (*Role, error) {
	url := fmt.Sprintf(clientRoleCreateUri, c.url, realm, clientId)
	return c.createRole(role, url)
}

func (c *KeycloakClient) createRole(role *Role, url string) (*Role, error) {
	roleLocation, err := c.post(url, *role)
	if err != nil {
		return nil, err
	}

	// if the roel name contains a slash, the location gets messed up (because
	// slashes dont get URL encoded by keycloak as they should. this will break
	// if the role name has a URL escaped character in it...
	suffix := roleLocation[len(url)+1:]
	suffix, err = neturl.QueryUnescape(suffix)
	if err != nil {
		return nil, err
	}
	suffix = neturl.QueryEscape(suffix)
	roleLocation = url + "/" + suffix

	var createdRole Role
	err = c.get(roleLocation, &createdRole)

	return &createdRole, err
}

func (c *KeycloakClient) UpdateRole(role *Role, realm string) error {
	url := fmt.Sprintf(roleUri, c.url, realm, role.Id)
	return c.put(url, *role)
}

func (c *KeycloakClient) DeleteRole(id, realm string) error {
	url := fmt.Sprintf(roleUri, c.url, realm, id)
	return c.delete(url, nil)
}
