package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v32/github"
	"github.com/russross/blackfriday"
)

var searchlight *api.Product

func addSearchlight() {
	searchlight = &api.Product{
		Name: "Searchlight",
		SocialLinks: map[string]string{
			"twitter":  "https://twitter.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/searchlight/searchlight",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/searchlight`,
		},
		StripeProductID: "prod_FARWHv38XFKlXL",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "searchlight", "searchlight", nil)
	if err != nil {
		searchlight.Description = nil
	} else {
		md, _ := description.GetContent()
		searchlight.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
