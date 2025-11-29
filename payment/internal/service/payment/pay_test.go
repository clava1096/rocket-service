package payment

import (
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestPay_Success() {
	transactionUUID := s.service.Pay(s.ctx)

	s.NotEmpty(transactionUUID)

	_, err := uuid.Parse(transactionUUID)
	s.NoError(err)
}
