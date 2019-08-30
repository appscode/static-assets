package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v28/github"
	"github.com/russross/blackfriday"
)

var kubeci *api.Product

func addKubeCI() {
	kubeci = &api.Product{
		Name: "KubeCI",
		SocialLinks: map[string]string{
			"forum":    "discourse.appscode.com",
			"twitter":  "https://twitter.com/appscode",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/kube-ci/engine",
			"youtube":  "https://www.youtube.com/channel/UCxObRDZ0DtaQe_cCP-dN-xg",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/kube-ci`,
		},
		StripeProductID: "prod_FFi3FtILZ6xXV0",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "kube-ci", "engine", nil)
	if err != nil {
		kubeci.Description = nil
	} else {
		md, _ := description.GetContent()
		kubeci.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.Run([]byte(md))),
		}
	}
}
