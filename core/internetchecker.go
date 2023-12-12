package core

import "sync"

type InternetChecker interface {
	HasInternetAccess() bool
}

type InternetCheckerImpl struct {
	HttpChecker           HttpChecker
	TimeoutInMilliseconds int
}

func (c *InternetCheckerImpl) HasInternetAccess() bool {
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
