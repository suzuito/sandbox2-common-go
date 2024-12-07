// このパッケージはsandbox2のサービスで使用される定義群を提供する
package services

import (
	"github.com/google/uuid"
)

// サービスを指すユニークID
type ServiceID uuid.UUID

func (s *ServiceID) UUID() uuid.UUID {
	return uuid.UUID(*s)
}

// サービスとは、何らかのストーリーを実現するソフトウェアを指す
type Service struct {
	// サービスを指すユニークID
	ID ServiceID
	// サービスを指すユニークな名前
	UniqueName string
}

var (
	Blog Service = registerService(Service{
		ID:         ServiceID(uuid.MustParse("c90db9d1-e1b3-45c9-8016-56e3c7355db8")),
		UniqueName: "blog3",
	})
)

var availableServices = []Service{}

func registerService(s Service) Service {
	availableServices = append(availableServices, s)
	return s
}
