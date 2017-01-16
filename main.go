package main

import (
	"flag"
	"fmt"
	"time"

	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

// Parameters variable
var daysStoppedBeforeClosed, frequency *int
var autoCleanImage, simulate *bool
var filtersEntry *string

func loadParameter() {
	daysStoppedBeforeClosed = flag.Int("days", 7, "Number of days to wait before removing stopped container")
	simulate = flag.Bool("simulate", false, "Do not perform actions, only display them.")
	autoCleanImage = flag.Bool("clean-images", true, "Automatically remove image that does not have any container using it.")
	frequency = flag.Int("frequency-seconds", 3600, "Number of seconds to wait between every execution.")
	filtersEntry = flag.String("filters", "", "image name to never remove. Define exclude filters from the clean image, separate filter with # character.")
	flag.Parse()
	return
}

func cleanContainers(client *docker.Client) {
	fmt.Println("Cleaning containers...")
	refTime := time.Now().Add(time.Duration(*daysStoppedBeforeClosed*24*-1) * time.Hour)

	// list all stop container and remove them
	options := docker.ListContainersOptions{
		All:     true,
		Size:    false,
		Limit:   100,
		Filters: map[string][]string{"status": {"exited"}},
	}
	containers, err := client.ListContainers(options)
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		inspectRes, err := client.InspectContainer(container.ID)
		if err != nil {
			panic(err)
		}
		if inspectRes.State.FinishedAt.Before(refTime) {
			fmt.Println("Container", container.ID, "(", container.Names, ")", "is older than", *daysStoppedBeforeClosed, "days since last exiting, removing...")
			if !(*simulate) {
				err := client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
				if err != nil {
					panic(err)
				}
			}
		}
	}
	fmt.Println("Cleaned containers.")
}

func isFilterImage(repoTags []string, filters []string) bool {
	result := false
	for _, filter := range filters {
		if len(filter) != 0 {
			for _, repoTag := range repoTags {
				result = (result || strings.Contains(repoTag, filter))
			}
		}
	}
	return result
}

func cleanImages(client *docker.Client, filters []string) {
	fmt.Println("Cleaning images...")
	images, err := client.ListImages(docker.ListImagesOptions{Filters: map[string][]string{"dangling": {"true"}}})
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		if isFilterImage(image.RepoTags, filters) {
			fmt.Println("Image", image.ID, "(", image.RepoTags, ")", "match filter, it is not removed.")
		} else {
			fmt.Println("Image", image.ID, "(", image.RepoTags, ")", "is unused, removing...")
			if !(*simulate) {
				err := client.RemoveImage(image.ID)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	fmt.Println("Cleaned images.")
}

func main() {
	fmt.Println("Starting docker cleaner")
	loadParameter()
	filters := strings.Split(*filtersEntry, "#")

	fmt.Println("Config is:")
	fmt.Println("\tdays:", *daysStoppedBeforeClosed)
	fmt.Println("\tsimulate:", *simulate)
	fmt.Println("\tclean-images:", *autoCleanImage)
	fmt.Println("\tfrequency-seconds:", *frequency)
	fmt.Println("\tfilters:", filters)

	// connect docker
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	cleanContainers(client)
	cleanImages(client, filters)
	if *frequency > 0 {
		for true {
			time.Sleep(time.Duration(*frequency) * time.Second)
			cleanContainers(client)
			cleanImages(client, filters)
		}
	}
}
