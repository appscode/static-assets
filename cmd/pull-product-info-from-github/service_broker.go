package main

import (
	"context"

	"github.com/appscode/static-assets/api"
	"github.com/google/go-github/v32/github"
	"github.com/russross/blackfriday"
)

var serviceBroker *api.Product

func addServiceBroker() {
	serviceBroker = &api.Product{
		Name: "Service Broker",
		SocialLinks: map[string]string{
			"twitter":  "https://twitter.com/AppsCodeHQ",
			"facebook": "https://facebook.com/appscode",
			"linkedin": "https://www.linkedin.com/company/appscode",
			"github":   "https://github.com/appscode/service-broker",
			"youtube":  "https://www.youtube.com/c/appscodeinc",
		},
		SupportLinks: map[string]string{
			"Support URL": `https://appscode.freshdesk.com`,
			"Website URL": `https://www.appscode.com/products/service-broker`,
		},
		StripeProductID: "prod_FARXK4vNsfeLk3",
	}
	ctx := context.Background()
	client := github.NewClient(nil)

	description, _, err := client.Repositories.GetReadme(ctx, "appscode", "service-broker", nil)
	if err != nil {
		serviceBroker.Description = nil
	} else {
		md, _ := description.GetContent()
		serviceBroker.Description = map[string]string{
			"markdown": md,
			"html":     string(blackfriday.MarkdownCommon([]byte(md))),
		}
	}
}
