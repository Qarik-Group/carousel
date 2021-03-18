package main

import (
	"github.com/starkandawyne/carousel/resource"
	"github.com/cloudboss/ofcourse/ofcourse"
)

func main() {
	ofcourse.Check(&resource.Resource{})
}
