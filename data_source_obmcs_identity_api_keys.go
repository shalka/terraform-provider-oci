// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"context"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oracle/bmcs-go-sdk"

	"github.com/oracle/terraform-provider-oci/crud"

	"bitbucket.aka.lgl.grungy.us/golang-sdk2/common"
	"bitbucket.aka.lgl.grungy.us/golang-sdk2/identity"
)

func APIKeyDatasource() *schema.Resource {
	return &schema.Resource{
		Read: readAPIKeys,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     APIKeyResource(),
			},
		},
	}
}

func readAPIKeys(d *schema.ResourceData, m interface{}) (e error) {
	// myclient := identity.IdentityClient{}

	client := m.(*OracleClients)
	sync := &APIKeyDatasourceCrud{}
	sync.D = d
	sync.Client = client.client
	return crud.ReadResource(sync)
}

type APIKeyDatasourceCrud struct {
	crud.BaseCrud
	Response identity.ListApiKeysResponse
}

func (s *APIKeyDatasourceCrud) Get() (e error) {
	userID := s.D.Get("user_id").(string)
	client := identity.IdentityClient{}
	request := identity.ListApiKeysRequest{
		common.String(userID),
	}
	s.Response, e = client.ListApiKeys(context.Background(), request)
	return
}

func (s *APIKeyDatasourceCrud) SetData() {
	if s.Res != nil {
		s.D.SetId(time.Now().UTC().String())
		resources := []map[string]interface{}{}
		for _, v := range s.Res.Keys {
			res := map[string]interface{}{
				"fingerprint":  v.Fingerprint,
				"id":           v.KeyID,
				"key_value":    v.KeyValue,
				"state":        v.State,
				"time_created": v.TimeCreated.String(),
				"user_id":      v.UserID,
			}
			resources = append(resources, res)
		}
		s.D.Set("api_keys", resources)
	}
	return
}
