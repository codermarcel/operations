package src

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

//Profile contains key value pairs
type Profile map[string]interface{}

type Service struct {
	Profiles map[string]Profile
}

func (s *Service) GetProfileNames() []string {
	res := []string{}
	for key := range s.Profiles {
		res = append(res, key)
	}

	return res
}

func (s *Service) HasProfile(name string) bool {
	_, found := s.Profiles[name]
	return found
}

type serviceOption func(*Service) error

func ServiceFromFile(path string) serviceOption {
	return func(b *Service) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		fn := ServiceFromBytes(data)

		return fn(b)
	}
}

func ServiceFromBytes(data []byte) serviceOption {
	return func(b *Service) error {
		svc := new(Service)
		err := yaml.Unmarshal(data, svc)
		if err != nil {
			return err
		}
		b.Profiles = svc.Profiles
		return nil
	}
}

func NewService(options ...serviceOption) (*Service, error) {
	b := &Service{Profiles: map[string]Profile{}}

	for _, v := range options {
		err := v(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}
