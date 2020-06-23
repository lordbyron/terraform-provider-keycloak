// This file provides a Terraform resource for Keycloak clients
// The client resource is documented at http://www.keycloak.org/docs-api/3.1/rest-api/index.html#_clientrepresentation

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lordbyron/terraform-provider-keycloak/keycloak"
)

func resourceClient() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceClientRead),
		Create: schema.CreateFunc(resourceClientCreate),
		Update: schema.UpdateFunc(resourceClientUpdate),
		Delete: schema.DeleteFunc(resourceClientDelete),

		// Keycloak clients are importable by ID, but the realm must also be provided by the user.
		Importer: &schema.ResourceImporter{
			State: importClientHelper,
		},

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"client_authenticator_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "client-secret",
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "openid-connect",
			},
			"public_client": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"bearer_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_accounts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"direct_access_grants_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"implicit_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"standard_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"web_origins": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"root_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"admin_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"full_scope_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// Computed fields (i.e. things looked up in Keycloak after client creation)
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"service_account_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// only available if protocol is `saml`
			"saml_idp_descriptor_xml": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func importClientHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	realm, id, err := splitRealmId(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	d.Set("realm", realm)

	resourceClientRead(d, m)

	return []*schema.ResourceData{d}, nil
}

func resourceClientRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	client, err := c.GetClient(d.Id(), realm(d))
	if err != nil {
		// Nothing was found, so just always recreate (instead of erroring)
		d.SetId("")
		return nil
	}

	clientToResourceData(client, d)

	/** Get computed fields **/
	// Look up client secret in addition
	secret, err := c.GetClientSecret(d.Id(), realm(d))
	if err != nil {
		return err
	}
	d.Set("client_secret", secret.Value)

	// Look up service account user ID (if enabled)
	if client.ServiceAccountsEnabled {
		user, err := c.GetClientServiceAccountUser(d.Id(), realm(d))
		if err != nil {
			return err
		}

		d.Set("service_account_user_id", user.Id)
	}

	if client.Protocol == "saml" {
		installation, err := c.GetClientInstallationSamlDesc(d.Id(), realm(d))
		if err != nil {
			return err
		}

		d.Set("saml_idp_descriptor_xml", installation)
	}

	return nil
}

func resourceClientCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	client := resourceDataToClient(d)
	created, err := apiClient.CreateClient(&client, realm(d))

	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceClientRead(d, m)
}

func resourceClientUpdate(d *schema.ResourceData, m interface{}) error {
	client := resourceDataToClient(d)
	apiClient := m.(*keycloak.KeycloakClient)
	return apiClient.UpdateClient(&client, realm(d))
}

func resourceClientDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	return apiClient.DeleteClient(d.Id(), realm(d))
}

// Turns resource.tf files into the Client struct
func resourceDataToClient(d *schema.ResourceData) keycloak.Client {
	redirectUris := []string{}
	webOrigins := []string{}
	attributes := map[string]interface{}{}

	for _, uri := range d.Get("redirect_uris").([]interface{}) {
		redirectUris = append(redirectUris, uri.(string))
	}

	rawOrigins, present := d.GetOk("web_origins")
	if present {
		for _, origin := range rawOrigins.([]interface{}) {
			webOrigins = append(webOrigins, origin.(string))
		}
	}

	rawAttributes, present := d.GetOk("attributes")
	if present {
		for k, v := range rawAttributes.(map[string]interface{}) {
			attributes[k] = v
		}
	}

	c := keycloak.Client{
		ClientId:                  d.Get("client_id").(string),
		Enabled:                   d.Get("enabled").(bool),
		ClientAuthenticatorType:   d.Get("client_authenticator_type").(string),
		RedirectUris:              redirectUris,
		Protocol:                  d.Get("protocol").(string),
		PublicClient:              d.Get("public_client").(bool),
		BearerOnly:                d.Get("bearer_only").(bool),
		ServiceAccountsEnabled:    d.Get("service_accounts_enabled").(bool),
		DirectAccessGrantsEnabled: d.Get("direct_access_grants_enabled").(bool),
		ImplicitFlowEnabled:       d.Get("implicit_flow_enabled").(bool),
		StandardFlowEnabled:       d.Get("standard_flow_enabled").(bool),
		WebOrigins:                webOrigins,
		RootUrl:                   d.Get("root_url").(string),
		AdminUrl:                  d.Get("admin_url").(string),
		BaseUrl:                   d.Get("base_url").(string),
		FullScopeAllowed:          d.Get("full_scope_allowed").(bool),
		Attributes:                attributes,
	}

	if !d.IsNewResource() {
		c.Id = d.Id()
	}

	return c
}

// Turns the struct (from a GET) into the internal representation
func clientToResourceData(c *keycloak.Client, d *schema.ResourceData) {
	d.Set("client_id", c.ClientId)
	d.Set("enabled", c.Enabled)
	d.Set("client_authenticator_type", c.ClientAuthenticatorType)
	d.Set("redirect_uris", c.RedirectUris)
	d.Set("protocol", c.Protocol)
	d.Set("public_client", c.PublicClient)
	d.Set("bearer_only", c.BearerOnly)
	d.Set("service_accounts_enabled", c.ServiceAccountsEnabled)
	d.Set("direct_access_grants_enabled", c.DirectAccessGrantsEnabled)
	d.Set("implicit_flow_enabled", c.ImplicitFlowEnabled)
	d.Set("standard_flow_enabled", c.StandardFlowEnabled)
	d.Set("web_origins", c.WebOrigins)
	d.Set("root_url", c.RootUrl)
	d.Set("admin_url", c.AdminUrl)
	d.Set("base_url", c.BaseUrl)
	d.Set("full_scope_allowed", c.FullScopeAllowed)
	d.Set("attributes", c.Attributes)
}

func defaultClientAttributes() map[string]interface{} {
	return map[string]interface{}{
		"saml.assertion.signature": true,
		"xray":                     "foobar",
	}
}
