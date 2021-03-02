package state

type Collector func(c *Credential) Credentials

func SignedByCollector() Collector {
	return func(c *Credential) Credentials {
		return Credentials{c.SignedBy}
	}
}

func SignsCollector() Collector {
	return func(c *Credential) Credentials {
		return c.Signs
	}
}

func SibilingsCollector() Collector {
	return func(c *Credential) Credentials {
		return c.Path.Versions
	}
}
