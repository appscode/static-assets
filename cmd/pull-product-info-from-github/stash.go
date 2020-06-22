package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v32/github"
	"github.com/russross/blackfriday"
)

var stash *api.Product

func addStash() {
	stash = &api.Product{
		Name: "Stash",
		SocialLinks: map[string]string{
			"twitter":  "https://twitter.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/stashed/stash",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/stash`,
		},
		StripeProductID: "prod_FARXQdMCCvjfZw",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "stashed", "stash", nil)
	if err != nil {
		stash.Description = nil
	} else {
		md, _ := description.GetContent()
		stash.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
