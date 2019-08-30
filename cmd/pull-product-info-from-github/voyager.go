package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v28/github"
	"github.com/russross/blackfriday"
)

var voyager *api.Product

func addVoyager() {
	voyager = &api.Product{
		Name: "Voyager",
		SocialLinks: map[string]string{
			"forum":    "discourse.appscode.com",
			"twitter":  "https://twitter.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/appscode/voyager",
			"youtube":  "https://www.youtube.com/channel/UCxObRDZ0DtaQe_cCP-dN-xg",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://appscode.com/products/voyager`,
		},
		StripeProductID: "prod_FARYkUPIIB0Ocx",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "appscode", "voyager", nil)
	if err != nil {
		voyager.Description = nil
	} else {
		md, _ := description.GetContent()
		voyager.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
