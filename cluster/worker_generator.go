package main

import (
	"flag"
	"os"
	"path/filepath"
	"text/template"
)

const regions = 16
const deploymentPath = "templates/worker-deployment.template.yaml"
const servicePath = "templates/worker-service.template.yaml"
const hpaPath = "templates/worker-hpa.template.yaml"

func main() {
	version := flag.String("version", "0.0.1", "Specifies new version of the image")
	flag.Parse()

	deploymentTemplate := template.Must(template.ParseFiles(deploymentPath))
	serviceTemplate := template.Must(template.ParseFiles(servicePath))
	hpaTemplate := template.Must(template.ParseFiles(hpaPath))

	deployments, err := os.Create("workers-deployment.yaml")
	if err != nil {
		panic(err)
	}
	defer deployments.Close()
	services, err := os.Create("workers-service.yaml")
	if err != nil {
		panic(err)
	}
	defer services.Close()
	hpa, err := os.Create("workers-hpa.yaml")
	if err != nil {
		panic(err)
	}
	defer hpa.Close()

	for i := 0; i < regions; i++ {
		data := struct {
			Region  int
			Version string
		}{i, *version}

		if err := deploymentTemplate.ExecuteTemplate(deployments, filepath.Base(deploymentPath), data); err != nil {
			panic(err)
		}
		if err := serviceTemplate.ExecuteTemplate(services, filepath.Base(servicePath), data); err != nil {
			panic(err)
		}
		if err := hpaTemplate.ExecuteTemplate(hpa, filepath.Base(hpaPath), data); err != nil {
			panic(err)
		}

		if _, err := deployments.WriteString("\n---\n"); err != nil {
			panic(err)
		}
		if _, err := services.WriteString("\n---\n"); err != nil {
			panic(err)
		}
		if _, err := hpa.WriteString("\n---\n"); err != nil {
			panic(err)
		}
	}
}
