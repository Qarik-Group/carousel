package state

func (s *state) Credentials(filters ...filter) []*Credential {
	certs := s.credentials.Select(func(_, v interface{}) bool {
		for _, fn := range filters {
			if !fn(v.(*Credential)) {
				return false
			}
		}
		return true
	})
	out := make([]*Credential, 0, certs.Size())
	for _, cert := range certs.Values() {
		out = append(out, cert.(*Credential))
	}
	return out
}
