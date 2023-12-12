package core

type InternetChecker struct {
	HttpChecker HttpChecker
}

func (c *InternetChecker) HasInternetAccess() bool {
	
	return c.HttpChecker.HasDirectInternetAccess() || c.HttpChecker.HasInternetAccessViaProxy()
}
