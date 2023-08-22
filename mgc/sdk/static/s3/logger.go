package s3

import "fmt"

func LogSigningInfo(canonicalStr, strToSign string) {
	fmt.Printf(logSignInfoMsg, canonicalStr, strToSign)
}

const logSignInfoMsg = `Request Signature:
---[ CANONICAL STRING  ]-----------------------------
%s
---[ STRING TO SIGN ]--------------------------------
%s
-----------------------------------------------------
`
