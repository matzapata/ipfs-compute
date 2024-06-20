package services

import "github.com/wabarc/ipfs-pinner/pkg/pinata"

type IpfsService struct {
	pnt pinata.Pinata
}

func NewIpfsService(apiKey string, secret string) *IpfsService {
	return &IpfsService{
		pnt: pinata.Pinata{Apikey: apiKey, Secret: secret},
	}
}

func (s *IpfsService) PinFile(file string) (string, error) {
	return s.pnt.PinFile(file)
}
