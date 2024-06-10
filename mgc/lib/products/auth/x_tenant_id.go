package auth

type XTenantIDParameters struct {
	Key string
}

func (a XTenantIDParameters) GetXTenantID() string {
	return a.Key
}
