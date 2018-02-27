package keycloak

import (
	"fmt"
	"log"
)

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

const (
	groupsUri = "%s/auth/admin/realms/%s/groups"
)

func (c *KeycloakClient) GetGroupByName(name, realm string) (*Group, error) {
	url := fmt.Sprintf(groupsUri, c.url, realm)
	log.Println(url)

	var groups []Group
	err := c.get(url, &groups)
	if err != nil {
		return nil, err
	}
	for _, u := range groups {
		// make sure it's an exact match
		if name == u.Name {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("Exact match search failed to find %s", name)
}
