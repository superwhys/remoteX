package hash

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestAdler32RollingUpdate(t *testing.T) {
	blockLength := 5
	buf := []byte("helloworld")
	sum := CheckAdlerSum(buf[:blockLength])
	rollingSum := RollingUpdate(sum, buf[0], buf[5], uint32(blockLength))
	actualRollingSum := CheckAdlerSum(buf[1 : 1+blockLength])
	t.Logf("origin sum: %v, rolling 1 byte(%s) sum: %v, actual rolling sum: %v", sum, buf[1:1+blockLength], rollingSum, actualRollingSum)
	
	assert.Equal(t, actualRollingSum, rollingSum)
}
