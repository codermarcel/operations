package main

import (
	"flag"
	"fmt"
	"log"
	"operations/src"
	"os"
)

func main() {
	var serviceFile, templateFile, profileName, outFile string

	flag.StringVar(&serviceFile, "service", "service.yaml", "the path of the service.yaml file")
	flag.StringVar(&profileName, "profile", "default", "the name of the profile to use")
	flag.StringVar(&templateFile, "template", "template.yaml", "the path of the template file")
	flag.StringVar(&outFile, "out", "out.yaml", "the output filename")
	flag.Parse()

	svc, err := src.NewService(src.ServiceFromFile(serviceFile))

	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := src.NewBaseTemplate(src.TemplateFromFile(templateFile))

	if err != nil {
		log.Fatal(err)
	}

	if !svc.HasProfile(profileName) {
		log.Fatalf("profile with name %s not found.", profileName)
	}

	f, err := os.Create(outFile)

	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteTo(svc.Profiles[profileName], f)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Success!")
	os.Exit(0)
}
