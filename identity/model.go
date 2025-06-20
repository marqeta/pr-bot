package identity

type CallerIdentity struct {
	Arn     string `json:"Arn"`
	Account string `json:"Account"`
	UserID  string `json:"UserId"`
}

type GetCallerIdentityJSON struct {
	GetCallerIdentityResponse struct {
		GetCallerIdentityResult CallerIdentity `json:"GetCallerIdentityResult"`
	} `json:"GetCallerIdentityResponse"`
}
