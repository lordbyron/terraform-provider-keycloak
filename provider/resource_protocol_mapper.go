package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lordbyron/terraform-provider-keycloak/keycloak"
)

func resourceProtocolMapper() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceProtocolMapperRead),
		Create: schema.CreateFunc(resourceProtocolMapperCreate),
		Update: schema.UpdateFunc(resourceProtocolMapperUpdate),
		Delete: schema.DeleteFunc(resourceProtocolMapperDelete),

		// ProtocolMappers are importable by ID
		Importer: &schema.ResourceImporter{
			State: importProtocolMapperHelper,
		},

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol_mapper": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consent_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"consent_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func importProtocolMapperHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	realm, client_id, id, err := splitRealmClientId(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	d.Set("realm", realm)
	d.Set("client_id", client_id)

	resourceProtocolMapperRead(d, m)

	return []*schema.ResourceData{d}, nil
}

func resourceProtocolMapperRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	pm, err := c.GetProtocolMapper(d.Id(), realm(d), client(d))
	if err != nil {
		// Nothing was found, so just always recreate (instead of erroring)
		d.SetId("")
		return nil
	}

	protocolMapperToResourceData(pm, d)

	return nil
}

func resourceProtocolMapperCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	pm := resourceDataToProtocolMapper(d)
	created, err := c.CreateProtocolMapper(&pm, realm(d), client(d))

	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceProtocolMapperRead(d, m)
}

func resourceProtocolMapperUpdate(d *schema.ResourceData, m interface{}) error {
	pm := resourceDataToProtocolMapper(d)
	c := m.(*keycloak.KeycloakClient)
	return c.UpdateProtocolMapper(&pm, realm(d), client(d))
}

func resourceProtocolMapperDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	return c.DeleteProtocolMapper(d.Id(), realm(d), client(d))
}

func resourceDataToProtocolMapper(d *schema.ResourceData) keycloak.ProtocolMapper {
	config := map[string]interface{}{}

	rawConfig, present := d.GetOk("config")
	if present {
		for k, v := range rawConfig.(map[string]interface{}) {
			config[k] = v
		}
	}

	pm := keycloak.ProtocolMapper{
		Name:            d.Get("name").(string),
		Protocol:        d.Get("protocol").(string),
		ProtocolMapper:  d.Get("protocol_mapper").(string),
		ConsentRequired: d.Get("consent_required").(bool),
		ConsentText:     d.Get("consent_text").(string),
		Config:          config,
	}

	if !d.IsNewResource() {
		pm.Id = d.Id()
	}

	return pm
}

// Turns the struct into the internal representation
func protocolMapperToResourceData(pm *keycloak.ProtocolMapper, d *schema.ResourceData) {
	d.Set("name", pm.Name)
	d.Set("protocol", pm.Protocol)
	d.Set("protocol_mapper", pm.ProtocolMapper)
	d.Set("consent_required", pm.ConsentRequired)
	d.Set("consent_text", pm.ConsentText)
	d.Set("config", pm.Config)
}
