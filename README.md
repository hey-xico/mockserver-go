# mockserver-go
A Fluent API that helps you create mocks and proxies using [MockServer](https://www.mock-server.com/) standalone or running on Docker

## MockServer
>**What is MockServer**
> 
>For any system you integrate with via HTTP or HTTPS MockServer can be used as:
>
>- a mock configured to return specific responses for different requests
>- a proxy recording and optionally modifying requests and responses
>- both a proxy for some requests and a mock for other requests at the same time
> 
>When MockServer receives a request it matches the request against active expectations that have been configured, if no matches are found it proxies the request if appropriate otherwise a 404 is returned.

Beware that this API might change before v1.

Up til now it supports only Mock with strict matchers.

```go
msAddress := "http://localhost:1080"

msClient := mockserver.NewMockServerClient(msAddress, &http.Client{})

expectation := msClient.
    Expect(t, mockserver.NewExpectation().
        When(mockserver.ARequest().
            WithMethod("PATCH").
            WithHeader(map[string]string{
                "Content-Type": "application/json"}).
            WithPath("/foo/{fooId}").
            WithPathVariable(map[string][]string{
                "fooId": {"987654321"}}).
            WithBody("{\"fooName\": \"John Doe\"}")).
        Respond(mockserver.AResponse().
            WithStatusCode(200)))
		
//test something

msClient.
    Verify(t, mockserver.NewVerifier().
        Expectation(expectation).
        ToBeCalled(mockserver.Once()).
        Retry(5))

		
```
