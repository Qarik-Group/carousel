package state

func pathComparator(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Path)
	c2 := b.(*Path)

	switch {
	case c1.Name > c2.Name:
		return 1
	case c1.Name < c2.Name:
		return -1
	default:
		return 0
	}
}

func credentialComparator(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Credential)
	c2 := b.(*Credential)

	switch {
	case c1.ID > c2.ID:
		return 1
	case c1.ID < c2.ID:
		return -1
	default:
		return 0
	}
}

func deploymentComparator(a, b interface{}) int {
	// Type assertion, program will panic if this is not respected
	c1 := a.(*Deployment)
	c2 := b.(*Deployment)

	switch {
	case c1.Name > c2.Name:
		return 1
	case c1.Name < c2.Name:
		return -1
	default:
		return 0
	}
}
