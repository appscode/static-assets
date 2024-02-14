package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v54/github"
	"github.com/russross/blackfriday"
)

var kubeshield *api.Product

func addKubeShield() {
	kubeshield = &api.Product{
		Name: "KubeShield",
		SocialLinks: map[string]string{
			"twitter":  "https://x.com/kubeshield",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/kubeshield/scanner",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://kubeshield.app/`,
		},
		StripeProductID: "prod_FFhyh5Z56cmEJL",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "kubeshield", "scanner", nil)
	if err != nil {
		kubeshield.Description = nil
	} else {
		md, _ := description.GetContent()
		kubeshield.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
