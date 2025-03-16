package handler

import (
	"golang.org/x/exp/slices"
	"strings"
)

// Project represents a project with a title, description, and URL
type Portfolio struct {
	Title       string
	ImagePath   string
	Date        string
	Description string
	URL         string
	Tags        []string
}

var imgPath = "assets/imgs/portfolio_imgs/"

// projects holds the list of projects
var portfolioList = []Portfolio{
	{
		Title:       "Homelab",
		ImagePath:   imgPath + "homelab.png",
		Description: "Personal Homelab hosting networking systems and multiple services, ranging from virtual machines, to highly available microservices.",
		URL:         "https://example.com/project-one",
		Tags:        []string{"Project", "Homelab"},
	},
	{
		Title:       "3 Tier Web App",
		ImagePath:   imgPath + "3tierWeb.jpg",
		Description: "Created a highly available word press website using AWS services. Available in Terraform for ease of deployment.",
		URL:         "https://github.com/Brandhang34/ThreeTierWebApp/tree/main",
		Tags:        []string{"Project", "AWS"},
	},
	{
		Title:       "Serverless CI/CD Resume",
		ImagePath:   imgPath + "crc.png",
		Description: "Static website stored in S3 Bucket and served by a CDN in AWS. Includes DevOps work flows such as CI/CD pipelines and IaC.",
		URL:         "https://github.com/Brandhang34/Serverless-CICD-Resume",
		Tags:        []string{"Project", "AWS"},
	},
	{
		Title:       "ByteSender",
		ImagePath:   imgPath + "blank.jpg",
		Description: "CLI tool allows you to upload files to S3 and share files.",
		URL:         "https://example.com/project-three",
		Tags:        []string{"Project", "Golang", "Docker", "AWS"},
	},
	{
		Title:       "Portfolio",
		ImagePath:   imgPath + "blank.jpg",
		Description: "Website written in golang, htmx, and tailwindcss.",
		URL:         "https://example.com/project-three",
		Tags:        []string{"Project", "Golang", "Docker"},
	},
	{
		Title:       "Minecraft Server",
		ImagePath:   imgPath + "minecraft.jpg",
		Description: "Self hosted Minecraft sever, utilzing networking tunnels and handlers written as systemd services.",
		URL:         "https://example.com/project-three",
		Tags:        []string{"Wiki", "Homelab"},
	},
	{
		Title:       "Dotfiles",
		ImagePath:   imgPath + "blank.jpg",
		Description: "My terminal setup/configurations, including Ghostty terminal, Starship Prompt, Helix text editor and Fish shell.",
		URL:         "https://example.com/project-three",
		Tags:        []string{"Project"},
	},
	// {
	// 	Title:       "Kubernetes Cluster",
	// ImagePath:   imgPath + "blank.jpg",
	// 	Description: "Website written in golang, htmx, and picocss. Hosted in a highly available kubernetes cluster and deployed using docker image",
	// 	URL:         "https://example.com/project-three",
	// 	Tags:        []string{"Wiki", "Homelab"},
	// },
	{
		Title:       "Cloud Windows Gaming VM",
		ImagePath:   imgPath + "cloud-gaming.png",
		Description: "Self hosted VDI with a GPU passed through for any video rendering or remote gaming.",
		URL:         "https://example.com/project-three",
		Tags:        []string{"Wiki", "Homelab"},
	},
	{
		Title:       "AWS Cloud Practitioner",
		ImagePath:   imgPath + "AWS_CCP.png",
		Description: "Certified in AWS Cloud Practitioner",
		URL:         "https://www.credly.com/badges/7252e2a1-5f01-4590-961a-5dbb9671a556/public_url",
		Tags:        []string{"Certification"},
	},
	{
		Title:       "AWS Solutions Architect - Associate",
		ImagePath:   imgPath + "AWS_SAA.jpg",
		Description: "Certified in AWS Solutions Architect Associate",
		URL:         "https://www.credly.com/badges/7cd3a6be-ba3d-41f0-8431-1a439e39c63b/public_url",
		Tags:        []string{"Certification"},
	},
	// {
	// 	Title:       "AWS Developer - Associate",
	// ImagePath:   imgPath + "blank.jpg",
	// 	Description: "Website written in golang, htmx, and picocss. Hosted in a highly available kubernetes cluster and deployed using docker image",
	// 	URL:         "https://example.com/project-three",
	// 	Tags:        []string{"Certification"},
	// },
	{
		Title:       "Hashicorp Terraform Associate",
		ImagePath:   imgPath + "HCTA_03.jpg",
		Description: "Certified in Hashicorp Terraform Associate",
		URL:         "https://www.credly.com/badges/65b3ae0d-cf3f-4f6d-8140-8e7426b98722/public_url",
		Tags:        []string{"Certification"},
	},
}

func LoadAllPortfolioItems() []Portfolio {
	return portfolioList
}

func GetProjectTag() []Portfolio {
	taggedProjects := []Portfolio{}
	for _, project := range portfolioList {
		if slices.Contains(project.Tags, "Project") {
			taggedProjects = append(taggedProjects, project)
		}

	}
	return taggedProjects
}

func GetCertificationTag() []Portfolio {
	certificationsArray := []Portfolio{}
	for _, project := range portfolioList {
		if slices.Contains(project.Tags, "Certification") {
			certificationsArray = append(certificationsArray, project)
		}
	}
	return certificationsArray
}

func GetHomelabTag() []Portfolio {
	homelabArray := []Portfolio{}
	for _, project := range portfolioList {
		if slices.Contains(project.Tags, "Homelab") {
			homelabArray = append(homelabArray, project)
		}

	}
	return homelabArray
}

func SearchProjectsHandler(searchQuery string, tagsQuery string) []Portfolio {
	list := []Portfolio{}
	filteredPortfolioList := []Portfolio{}
	searchQuery = strings.ToLower(searchQuery)

	// Filter the list by tags
	if tagsQuery == "projects" {
		list = GetProjectTag()
	} else if tagsQuery == "certifications" {
		list = GetCertificationTag()
	} else if tagsQuery == "homelab" {
		list = GetHomelabTag()
	} else {
		list = portfolioList
	}

	// Perform Search functionality
	for _, portfolioItem := range list {
		searchByTitle := strings.Contains(strings.ToLower(portfolioItem.Title), searchQuery)
		if searchByTitle {
			filteredPortfolioList = append(filteredPortfolioList, portfolioItem)
		}
	}
	return filteredPortfolioList
}
