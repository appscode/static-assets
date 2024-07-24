package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v63/github"
	"github.com/russross/blackfriday"
)

var guard *api.Product

func addGuard() {
	guard = &api.Product{
		Name: "Guard",
		SocialLinks: map[string]string{
			"twitter":  "https://x.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/appscode/guard",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/guard`,
		},
		StripeProductID: "prod_FARVnkJx8nelJk",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "appscode", "guard", nil)
	if err != nil {
		guard.Description = nil
	} else {
		md, _ := description.GetContent()
		guard.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
