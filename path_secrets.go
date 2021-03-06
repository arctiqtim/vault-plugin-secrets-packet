package packet

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathSecrets() *framework.Secret {
	return &framework.Secret{
		Type: secretType,
		Fields: map[string]*framework.FieldSchema{
			"api_key_token": {
				Type:        framework.TypeString,
				Description: "API token",
			},
		},
		Renew:  b.operationRenew,
		Revoke: b.operationRevoke,
	}
}

func (b *backend) operationRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	idRaw, ok := req.Secret.InternalData["api_key_id"]
	if !ok {
		return nil, fmt.Errorf("secret is missing ID of the API token")
	}
	keyID := idRaw.(string)
	client, err := b.Client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	resp, err := client.APIKeys.Delete(keyID)

	if (err != nil) && (resp.Response.StatusCode != 404) {
		return nil, err
	}

	return nil, nil
}

func (b *backend) operationRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	defaultLease, maxLease := b.getDefaultAndMaxLease()
	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = defaultLease
	resp.Secret.MaxTTL = maxLease
	return resp, nil
}

func (b *backend) getDefaultAndMaxLease() (time.Duration, time.Duration) {
	maxLease := b.system.MaxLeaseTTL()
	defaultLease := b.system.DefaultLeaseTTL()

	if defaultLease > maxLease {
		maxLease = defaultLease
	}
	return defaultLease, maxLease

}
