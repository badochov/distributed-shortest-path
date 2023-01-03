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

const deploymentLocalPath = "templates/worker-deployment.local.template.yaml"
const hpaLocalPath = "templates/worker-hpa.local.template.yaml"

func main() {
	version := flag.String("version", "0.0.1", "Specifies new version of the image.")
	local := flag.Bool("local", false, "If local deployments should be generated.")
	flag.Parse()

	dp := deploymentPath
	hp := hpaPath
	if *local {
		dp = deploymentLocalPath
		hp = hpaLocalPath
	}

	deploymentTemplate := template.Must(template.ParseFiles(dp))
	serviceTemplate := template.Must(template.ParseFiles(servicePath))
	hpaTemplate := template.Must(template.ParseFiles(hp))

	deployments, err := os.Create("generated/workers-deployment.yaml")
	if err != nil {
		panic(err)
	}
	defer deployments.Close()
	services, err := os.Create("generated/workers-service.yaml")
	if err != nil {
		panic(err)
	}
	defer services.Close()
	hpa, err := os.Create("generated/workers-hpa.yaml")
	if err != nil {
		panic(err)
	}
	defer hpa.Close()

	for i := 0; i < regions; i++ {
		data := struct {
			Region  int
			Version string
		}{i, *version}

		if err := deploymentTemplate.ExecuteTemplate(deployments, filepath.Base(dp), data); err != nil {
			panic(err)
		}
		if err := serviceTemplate.ExecuteTemplate(services, filepath.Base(servicePath), data); err != nil {
			panic(err)
		}
		if err := hpaTemplate.ExecuteTemplate(hpa, filepath.Base(hp), data); err != nil {
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
