package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lordbyron/terraform-provider-keycloak/keycloak"
)

func dataSourceClient() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceClientRead,
		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceClientRead(d *schema.ResourceData, m interface{}) error {

	id, idExists := d.GetOk("guid")
	name, nameExists := d.GetOk("client_id")

	if idExists && nameExists {
		return fmt.Errorf("guid and client_id arguments can't be used together")
	}
	if !idExists && !nameExists {
		return fmt.Errorf("Either of guid or client_id must be set")
	}

	//var client keycloak.Client
	if idExists {
		d.SetId(id.(string))
		return resourceClientRead(d, m)
	} else {
		// Find client by name
		c := m.(*keycloak.KeycloakClient)
		clients, err := c.ListClients(realm(d))
		if err != nil {
			return err
		}

		for _, client := range clients {
			if client.ClientId == name {
				d.SetId(client.Id)
				return resourceClientRead(d, m)
			}
		}
	}
	return nil
}
