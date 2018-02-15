package keycloak

import (
  "fmt"
)

type ProtocolMapper struct {
  Id              string                 `json:"id,omitempty"`
  Name            string                 `json:"name"`
  Protocol        string                 `json:"protocol,omitempty"`
  ProtocolMapper  string                 `json:"protocolMapper,omitempty"`
  ConsentRequired bool                   `json:"consentRequired,omitempty"`
  ConsentText     string                 `json:"consentText,omitempty"`
  Config          map[string]interface{} `json:"config,omitempty"`
}

const (
  mappersUri = "%s/auth/admin/realms/%s/clients/%s/protocol-mappers/models"
  mapperUri = "%s/auth/admin/realms/%s/clients/%s/protocol-mappers/models/%s"
)

func (c *KeycloakClient) GetProtocolMapper(id, realm, clientId string) (*ProtocolMapper, error) {
  url := fmt.Sprintf(mapperUri, c.url, realm, clientId, id)

  var pm ProtocolMapper
  err := c.get(url, &pm)

  if err != nil {
    return nil, err
  }

  return &pm, nil
}

func (c *KeycloakClient) CreateProtocolMapper(pm *ProtocolMapper, realm, clientId string) (*ProtocolMapper, error) {
  url := fmt.Sprintf(mappersUri, c.url, realm, clientId)

  mapperLocation, err := c.post(url, *pm)
  if err != nil {
    return nil, err
  }

  var createdMapper ProtocolMapper
  err = c.get(mapperLocation, &createdMapper)

  return &createdMapper, err
}

func (c *KeycloakClient) UpdateProtocolMapper(pm *ProtocolMapper, realm, clientId string) error {
  url := fmt.Sprintf(mapperUri, c.url, realm, clientId, pm.Id)
  return c.put(url, *pm)
}

func (c *KeycloakClient) DeleteProtocolMapper(id, realm, clientId string) error {
  url := fmt.Sprintf(mapperUri, c.url, realm, clientId, id)
  return c.delete(url, nil)
}
