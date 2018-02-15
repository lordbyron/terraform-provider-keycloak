package provider

import (
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "strings"
)

func realm(d *schema.ResourceData) string {
  return d.Get("realm").(string)
}

func client(d *schema.ResourceData) string {
  return d.Get("client_id").(string)
}

func getOptionalBool(d *schema.ResourceData, key string) *bool {
  if v, present := d.GetOk(key); present {
    b := v.(bool)
    return &b
  }
  return nil
}

func setOptionalBool(d *schema.ResourceData, key string, b *bool) {
  if b != nil {
    d.Set(key, *b)
  }
}

func getOptionalInt(d *schema.ResourceData, key string) *int {
  if v, present := d.GetOk(key); present {
    i := v.(int)
    return &i
  }
  return nil
}

func setOptionalInt(d *schema.ResourceData, key string, i *int) {
  if i != nil {
    d.Set(key, *i)
  }
}

func getStringSlice(d *schema.ResourceData, key string) []string {
  var stringSlice []string = []string{}
  untyped, present := d.GetOk(key)

  if !present {
    return stringSlice
  }

  for _, value := range untyped.([]interface{}) {
    stringSlice = append(stringSlice, value.(string))
  }

  return stringSlice
}

// This function is used when importing realm-specific resources. The realm must be specified by the user when
// importing by using a `${realm}.${resource_id}` syntax.
func splitRealmId(raw string) (string, string, error) {
  split := strings.Split(raw, ".")

  if len(split) != 2 {
    return "", "", fmt.Errorf("Import ID must be specified as '${realm}.${resource_id}'")
  }

  return split[0], split[1], nil
}

func splitRealmClientId(raw string) (string, string, string, error) {
  split := strings.Split(raw, ".")

  if len(split) != 3 {
    return "", "", "", fmt.Errorf("Import ID must be specified as '${realm}.${client_id}.${resource_id}' (n.b. client_id is not client name)")
  }

  return split[0], split[1], split[2], nil
}
