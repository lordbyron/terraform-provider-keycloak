package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lordbyron/terraform-provider-keycloak/keycloak"
)

func resourceRoleMapping() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceRoleMapRead),
		Create: schema.CreateFunc(resourceRoleMapCreate),
		Delete: schema.DeleteFunc(resourceRoleMapDelete),

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ForceNew:      true,
				ConflictsWith: []string{"group", "user_id", "group_id"},
			},
			"group": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ForceNew:      true,
				ConflictsWith: []string{"user", "user_id", "group_id"},
			},
			"user_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ForceNew:      true,
				ConflictsWith: []string{"user", "group", "group_id"},
			},
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ForceNew:      true,
				ConflictsWith: []string{"user", "group", "user_id"},
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
		},
	}
}

func resourceDataToRoleMap(d *schema.ResourceData) keycloak.RoleMapping {
	return keycloak.RoleMapping{
		Realm:     d.Get("realm").(string),
		RoleId:    d.Get("role_id").(string),
		UserName:  d.Get("user").(string),
		GroupName: d.Get("group").(string),
		UserId:    d.Get("user_id").(string),
		GroupId:   d.Get("group_id").(string),
		ClientId:  d.Get("client_id").(string),
	}
}

func resourceRoleMapRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	rm := resourceDataToRoleMap(d)
	err := rm.Validate(c)
	if err != nil {
		return err
	}

	roles, err := c.GetCompositeRoles(rm)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if role.Id == rm.RoleId {
			// nothing to do, but no error
			return nil
		}
	}
	return fmt.Errorf("No role mapping found for realm: %s, user: %s, group: %s, client: %s, role: %s", rm.Realm, rm.UserId, rm.GroupId, rm.ClientId, rm.RoleId)
}

func resourceRoleMapCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	rm := resourceDataToRoleMap(d)
	err := rm.Validate(c)
	if err != nil {
		return err
	}
	d.SetId(rm.Serialize())
	return c.AddRoleMapping(rm)
}

func resourceRoleMapDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	rm := resourceDataToRoleMap(d)
	err := rm.Validate(c)
	if err != nil {
		return err
	}
	return c.DeleteRoleMapping(rm)
}
