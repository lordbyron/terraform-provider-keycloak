package keycloak

import (
	"fmt"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"username"`
}

const (
	usersUriSearch = "%s/auth/admin/realms/%s/users?search=%s"
)

func (c *KeycloakClient) GetUserByName(name, realm string) (*User, error) {
	url := fmt.Sprintf(usersUriSearch, c.url, realm, name)

	var users []User
	err := c.get(url, &users)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		// make sure it's an exact match
		if name == u.Name {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("Exact match search failed to find %s", name)
}
