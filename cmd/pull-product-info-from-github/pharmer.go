package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v54/github"
	"github.com/russross/blackfriday"
)

var pharmer *api.Product

func addPharmer() {
	pharmer = &api.Product{
		Name: "Pharmer",
		SocialLinks: map[string]string{
			"twitter":  "https://x.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/pharmer/pharmer",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/pharmer`,
		},
		StripeProductID: "prod_FARVvB8RAs9eaU",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "pharmer", "pharmer", nil)
	if err != nil {
		pharmer.Description = nil
	} else {
		md, _ := description.GetContent()
		pharmer.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
