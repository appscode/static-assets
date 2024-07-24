package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v63/github"
	"github.com/russross/blackfriday"
)

var kubeform *api.Product

func addKubeform() {
	kubeform = &api.Product{
		Name: "Kubeform",
		SocialLinks: map[string]string{
			"twitter":  "https://x.com/Kubeform",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/kubeform",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.kubeform.com`,
		},
		StripeProductID: "",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "kubeform", "cli", nil)
	if err != nil {
		kubeform.Description = nil
	} else {
		md, _ := description.GetContent()
		kubeform.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
