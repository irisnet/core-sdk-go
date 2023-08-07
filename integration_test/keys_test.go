package integration_test

import (
	"github.com/stretchr/testify/require"
)

func (s IntegrationTestSuite) TestKeys() {
	name, password := s.RandStringOfLength(20), s.RandStringOfLength(8)

	address, mnemonic, err := s.Add(name, password)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), address)
	require.NotEmpty(s.T(), mnemonic)

	address1, err := s.Show(name, password)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, address1)

	privKeyArmor, err := s.Export(name, password)
	require.NoError(s.T(), err)

	err = s.Delete(name, password)
	require.NoError(s.T(), err)

	address2, err := s.Import(name, password, privKeyArmor)
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, address2)

	err = s.Delete(name, password)
	require.NoError(s.T(), err)

	// test Recover
	//address3, err := s.Key.Recover(name, password, mnemonic)
	//require.NoError(s.T(), err)
	//require.Equal(s.T(), address, address3)

	// test Recover With HD Path
	address4, err := s.RecoverWithHDPath(name, password, mnemonic, "")
	require.NoError(s.T(), err)
	require.Equal(s.T(), address, address4)
}
