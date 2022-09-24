package main

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"strings"

	staticassets "github.com/appscode/static-assets"
	"github.com/appscode/static-assets/api"
)

func main() {
	addGuard()
	addKubed()
	addKubeDB()
	addPharmer()
	addSearchlight()
	addServiceBroker()
	addStash()
	addSwift()
	addVoyager()
	addKubeShield()
	addKubeCI()
	addKubeVault()
	addKubeform()

	fsys := staticassets.Product()
	fs.WalkDir(fsys, "", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return err
		}

		key := strings.ReplaceAll(name, ".json", "")

		var p api.Product
		err = json.Unmarshal(data, &p)
		if err != nil {
			log.Fatalln("failed to deserialize", name)
		}

		switch key {
		case "guard":
			p.SocialLinks = guard.SocialLinks
			p.SupportLinks = guard.SupportLinks
			p.Description = guard.Description
			p.StripeProductID = guard.StripeProductID
		case "kubeci":
			p.SocialLinks = kubeci.SocialLinks
			p.SupportLinks = kubeci.SupportLinks
			p.Description = kubeci.Description
			p.StripeProductID = kubeci.StripeProductID
		case "kubed":
			p.SocialLinks = kubed.SocialLinks
			p.SupportLinks = kubed.SupportLinks
			p.Description = kubed.Description
			p.StripeProductID = kubed.StripeProductID
		case "kubedb":
			p.SocialLinks = kubedb.SocialLinks
			p.SupportLinks = kubedb.SupportLinks
			p.Description = kubedb.Description
			p.StripeProductID = kubedb.StripeProductID
		case "kubeshield":
			p.SocialLinks = kubeshield.SocialLinks
			p.SupportLinks = kubeshield.SupportLinks
			p.Description = kubeshield.Description
			p.StripeProductID = kubeshield.StripeProductID
		case "kubevault":
			p.SocialLinks = kubevault.SocialLinks
			p.SupportLinks = kubevault.SupportLinks
			p.Description = kubevault.Description
			p.StripeProductID = kubevault.StripeProductID
		case "pharmer":
			p.SocialLinks = pharmer.SocialLinks
			p.SupportLinks = pharmer.SupportLinks
			p.Description = pharmer.Description
			p.StripeProductID = pharmer.StripeProductID
		case "searchlight":
			p.SocialLinks = searchlight.SocialLinks
			p.SupportLinks = searchlight.SupportLinks
			p.Description = searchlight.Description
			p.StripeProductID = searchlight.StripeProductID
		case "service-broker":
			p.SocialLinks = serviceBroker.SocialLinks
			p.SupportLinks = serviceBroker.SupportLinks
			p.Description = serviceBroker.Description
			p.StripeProductID = serviceBroker.StripeProductID
		case "stash":
			p.SocialLinks = stash.SocialLinks
			p.SupportLinks = stash.SupportLinks
			p.Description = stash.Description
			p.StripeProductID = stash.StripeProductID
		case "swift":
			p.SocialLinks = swift.SocialLinks
			p.SupportLinks = swift.SupportLinks
			p.Description = swift.Description
			p.StripeProductID = swift.StripeProductID
		case "voyager":
			p.SocialLinks = voyager.SocialLinks
			p.SupportLinks = voyager.SupportLinks
			p.Description = voyager.Description
			p.StripeProductID = voyager.StripeProductID
		}

		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(p)
		if err != nil {
			log.Fatalln("failed to serialize", name)
		}

		filename := "/home/tamal/go/src/github.com/appscode/static-assets/data/products/" + name
		err = os.WriteFile(filename, buf.Bytes(), 0o644)
		if err != nil {
			log.Fatalln("failed to write", filename)
		}

		return nil
	})
}
