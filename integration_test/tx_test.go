package integration_test

func (s IntegrationTestSuite) TestTx() {
	block, err := s.Client.QueryBlock(666)
	s.T().Error(err)
	s.T().Error(block)
}
