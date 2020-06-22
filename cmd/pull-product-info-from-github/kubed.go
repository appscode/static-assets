package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v32/github"
	"github.com/russross/blackfriday"
)

var kubed *api.Product

func addKubed() {
	kubed = &api.Product{
		Name: "Kubed",
		SocialLinks: map[string]string{
			"twitter":  "https://twitter.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/appscode/kubed",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/kubed`,
		},
		StripeProductID: "prod_FARVnpFkkXyggO",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "appscode", "kubed", nil)
	if err != nil {
		kubed.Description = nil
	} else {
		md, _ := description.GetContent()
		kubed.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
