package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v28/github"
	"github.com/russross/blackfriday"
)

var swift *api.Product

func addSwift() {
	swift = &api.Product{
		Name: "Swift",
		SocialLinks: map[string]string{
			"twitter":  "https://twitter.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://kubepack.dev/swift",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/swift`,
		},
		StripeProductID: "prod_FARXloMWkxsKvP",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "appscode", "swift", nil)
	if err != nil {
		swift.Description = nil
	} else {
		md, _ := description.GetContent()
		swift.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
