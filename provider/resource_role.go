package provider

import (
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/lordbyron/terraform-provider-keycloak/keycloak"
)

func resourceRole() *schema.Resource {
  return &schema.Resource{
    // API methods
    Read:   schema.ReadFunc(resourceRoleRead),
    Create: schema.CreateFunc(resourceRoleCreate),
    Update: schema.UpdateFunc(resourceRoleUpdate),
    Delete: schema.DeleteFunc(resourceRoleDelete),

    // Roles are importable by ID
    Importer: &schema.ResourceImporter{
      State: importRoleHelper,
    },

    Schema: map[string]*schema.Schema{
      "realm": {
        Type:     schema.TypeString,
        Required: true,
      },
      "name": {
        Type:     schema.TypeString,
        Required: true,
        ForceNew: true,
      },
      "container_id": {
        Type:     schema.TypeString,
        Optional: true,
      },
      "description": {
        Type:     schema.TypeString,
        Optional: true,
      },
      "scope_param_requierd": {
        Type:     schema.TypeBool,
        Optional: true,
      },
      // Computed
      "client_role": {
        Type:     schema.TypeBool,
        Computed: true,
      },
    },
  }
}

func importRoleHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
  realm, id, err := splitRealmId(d.Id())
  if err != nil {
    return nil, err
  }

  d.SetId(id)
  d.Set("realm", realm)

  resourceRoleRead(d, m)

  return []*schema.ResourceData{d}, nil
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
  c := m.(*keycloak.KeycloakClient)

  role, err := c.GetRole(d.Id(), realm(d))
  if err != nil {
    return err
  }

  roleToResourceData(role, d)

  return nil
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
  c := m.(*keycloak.KeycloakClient)
  role := resourceDataToRole(d)
  var created *keycloak.Role
  var err error
  if role.ClientRole {
    created, err = c.CreateClientRole(&role, realm(d), role.ContainerId)
  } else {
    created, err = c.CreateRealmRole(&role, realm(d))
  }

  if err != nil {
    return err
  }

  d.SetId(created.Id)

  return resourceRoleRead(d, m)
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
  role := resourceDataToRole(d)
  c := m.(*keycloak.KeycloakClient)
  return c.UpdateRole(&role, realm(d))
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
  c := m.(*keycloak.KeycloakClient)
  return c.DeleteRole(d.Id(), realm(d))
}

func resourceDataToRole(d *schema.ResourceData) keycloak.Role {
  role := keycloak.Role{
    Name:               d.Get("name").(string),
    ClientRole:         d.Get("client_role").(bool),
    ContainerId:        d.Get("container_id").(string),
    Description:        d.Get("description").(string),
    ScopeParamRequired: d.Get("scope_param_requierd").(bool),
  }

  if role.ContainerId != "" {
    role.ClientRole = true
  }

  if !d.IsNewResource() {
    role.Id = d.Id()
  }

  return role
}

func roleToResourceData(role *keycloak.Role, d *schema.ResourceData) {
  d.Set("name", role.Name)
  d.Set("client_role", role.ClientRole)
  d.Set("container_id", role.ContainerId)
  d.Set("description", role.Description)
  d.Set("scope_param_requierd", role.ScopeParamRequired)
}
