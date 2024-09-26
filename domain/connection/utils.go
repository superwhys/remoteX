package connection

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateConnectionID(source, target string) string {
	rawID := fmt.Sprintf("%s-%s/%d", source, target, time.Now().Unix())
	
	hash := sha256.New()
	hash.Write([]byte(rawID))
	
	return hex.EncodeToString(hash.Sum(nil))
}
