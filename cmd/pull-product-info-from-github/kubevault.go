package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v54/github"
	"github.com/russross/blackfriday"
)

var kubevault *api.Product

func addKubeVault() {
	kubevault = &api.Product{
		Name: "KubeVault",
		SocialLinks: map[string]string{
			"twitter":  "https://twitter.com/kubevault",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/kubevault",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://kubevault.com`,
		},
		StripeProductID: "prod_FiawKtgLGy6R9v",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "kubevault", "docs", nil)
	if err != nil {
		kubevault.Description = nil
	} else {
		md, _ := description.GetContent()
		kubevault.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
