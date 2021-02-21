package state

import "bytes"

func (s *state) getCredentialBySubjectKeyId(keyId []byte) (*Credential, bool) {
	_, foundValue := s.credentials.Find(func(index interface{}, value interface{}) bool {
		if value.(*Credential).Certificate != nil {
			return bytes.Compare(value.(*Credential).Certificate.SubjectKeyId, keyId) == 0
		}
		return false
	})
	if foundValue != nil {
		return foundValue.(*Credential), true
	}
	return nil, false
}

func (s *state) eachPath(fn func(*Path)) {
	s.paths.Each(func(_, v interface{}) {
		fn(v.(*Path))
	})
}

func (s *state) getPath(name string) (*Path, bool) {
	i, found := s.paths.Get(name)
	if found {
		return i.(*Path), true
	}
	return nil, false
}

func (s *state) getOrCreateDeployment(name string) *Deployment {
	i, found := s.deployments.Get(name)
	if found {
		return i.(*Deployment)
	}
	d := &Deployment{
		Name:     name,
		Versions: make([]*Credential, 0),
	}
	s.deployments.Put(name, d)
	return d
}

func (s *state) getCredential(id string) (*Credential, bool) {
	i, found := s.credentials.Get(id)
	if found {
		return i.(*Credential), true
	}
	return nil, false
}

func (s *state) eachCredential(fn func(*Credential)) {
	s.credentials.Each(func(_, v interface{}) {
		fn(v.(*Credential))
	})
}

func (s *state) eachDeployment(fn func(*Deployment)) {
	s.deployments.Each(func(_, v interface{}) {
		fn(v.(*Deployment))
	})
}
