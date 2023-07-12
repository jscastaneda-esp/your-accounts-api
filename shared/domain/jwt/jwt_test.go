package jwt

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) TestJwtGenerateSuccess() {
	require := require.New(suite.T())

	token, _, err := JwtGenerate(1, "test", "test")

	require.NoError(err)
	require.NotEmpty(token)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
