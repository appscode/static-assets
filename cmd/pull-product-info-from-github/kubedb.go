package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v28/github"
	"github.com/russross/blackfriday"
)

var kubedb *api.Product

func addKubeDB() {
	kubedb = &api.Product{
		Name: "KubeDB",
		SocialLinks: map[string]string{
			"forum":    "discourse.appscode.com",
			"twitter":  "https://twitter.com/KubeDB",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/kubedb/cli",
			"youtube":  "https://www.youtube.com/channel/UCxObRDZ0DtaQe_cCP-dN-xg",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.kubedb.com`,
		},
		StripeProductID: "prod_FARWKSytMzmLb9",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "kubedb", "cli", nil)
	if err != nil {
		kubedb.Description = nil
	} else {
		md, _ := description.GetContent()
		kubedb.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.Run([]byte(md))),
		}
	}
}
