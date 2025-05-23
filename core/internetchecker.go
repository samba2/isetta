package core

import "sync"

type InternetChecker struct {
	HttpChecker           HttpChecker
	TimeoutInMilliseconds int
}

func NewInternetChecker(HttpChecker HttpChecker) InternetChecker {
	return InternetChecker{
		HttpChecker: HttpChecker,
		TimeoutInMilliseconds: 1000,
	}
}

func (c *InternetChecker) HasInternetAccess() bool {
	var wg sync.WaitGroup

	ch := make(chan bool, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		ch <- c.HttpChecker.HasDirectInternetAccess(c.TimeoutInMilliseconds)
	}()

	go func() {
		defer wg.Done()
		ch <- c.HttpChecker.HasInternetAccessViaProxy(c.TimeoutInMilliseconds)
	}()

	wg.Wait()
	close(ch)

	result1, result2 := <-ch, <-ch
	return result1 || result2
}
