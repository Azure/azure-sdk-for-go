package azfile

import (
	"github.com/stretchr/testify/assert"
)

func (s *azfileTestSuite) TestShareCreateRootDirectoryURL() {
	_assert := assert.New(s.T())
	testName := s.T().Name()
	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
	if err != nil {
		s.Fail("Unable to fetch service client because " + err.Error())
	}

	shareClient := createNewShare(_assert, generateShareName(testName), svcClient)
	defer deleteContainer(_assert, shareClient)
}
