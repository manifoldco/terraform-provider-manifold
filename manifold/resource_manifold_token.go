package manifold

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"

	manifold "github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/integrations"
)

func resourceManifoldToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceManifoldTokenCreate,
		Read:   resourceManifoldTokenRead,
		Update: resourceManifoldTokenUpdate,
		Delete: resourceManifoldTokenDelete,

		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The role this key should assume.",
			},

			"team": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The team you want to link the token to.",
			},

			"self": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Assign the token to user that is linked to the configuration token.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Generated by terraform.",
				Description: "The description to give to this key.",
			},

			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceManifoldTokenCreate(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	req := &manifold.APITokenRequest{
		Description: d.Get("description").(string),
		Role:        d.Get("role").(string),
		TeamID:      cl.TeamID,
	}

	if d.Get("self").(bool) {
		self, err := cl.Self.Get(ctx)
		if err != nil {
			return err
		}
		req.UserID = &self.ID
	}

	if team := d.Get("team").(string); team != "" {
		teamID, err := teamID(ctx, cl.Client, team)
		if err != nil {
			return err
		}

		req.TeamID = teamID
	}

	token, err := cl.Tokens.Create(ctx, req)
	if err != nil {
		return err
	}

	d.SetId(token.ID.String())
	return d.Set("token", *token.Body.Token)
}

func resourceManifoldTokenRead(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	req := &manifold.TokensListOpts{
		TeamID: cl.TeamID,
	}

	if d.Get("self").(bool) {
		req.Me = func(b bool) *bool { return &b }(true)
	}

	tokensList := cl.Tokens.List(ctx, "api", req)
	defer tokensList.Close()

	var token *manifold.APIToken
	for tokensList.Next() {
		st, err := tokensList.Current()
		if err != nil {
			return err
		}

		if st.ID.String() == d.Id() {
			token = st
		}
	}

	if token == nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceManifoldTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceManifoldTokenDelete(d, meta); err != nil {
		return err
	}

	return resourceManifoldTokenCreate(d, meta)
}

func resourceManifoldTokenDelete(d *schema.ResourceData, meta interface{}) error {
	cl := meta.(*integrations.Client)
	ctx := context.Background()

	return cl.Tokens.Delete(ctx, d.Id())
}

func teamID(ctx context.Context, cl *manifold.Client, label string) (*manifold.ID, error) {
	teamsList := cl.Teams.List(ctx)
	defer teamsList.Close()

	for teamsList.Next() {
		team, err := teamsList.Current()
		if err != nil {
			return nil, err
		}

		if team.Body.Label == label {
			return &team.ID, nil
		}
	}

	return nil, integrations.ErrTeamNotFound
}
