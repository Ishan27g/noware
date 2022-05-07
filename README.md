# noware
no(op-middle)ware for Go &amp; Ts

The `noop middleware` 
- propogates a `noop` context over http requests between various http services. 
- allows injection of `action events` into the context, enabling a log/event style response formulation.

In a nutshel this allows the ability to **test a live service/endpoint on a per request basis**, without the need to set any environment variables
or pass any url params in the request.

This middleware has none to no impact for regular/user based requests, but can provide an enriched log of the services when triggered by the test suite.

- Requests that do not contain the `noop` context are handled as normal. Everything is null, so it has minimal impact on performance or memory.
- The requests that do contain the `noop` context, will be send from a test suite. The `noop context` 

With `noop` We can assume the criteria of success for the method as everything it would normally done up till the final file persistence operation, db writes or external calls (assuming the external services are available)

A simple test suite can then send `noop` requests to these service's http wrappers, for which the receiving service would run the usual logic and return before the final external service calls.
The data that it would have normally sent to the external service can be sent back to the test suite. In this way, the need to generate mock data for testing each service or to mock a downstream service 
is removed. The test suite can then trigger the next service with this data.

The `noop middleware` can be specially useful in `async-systems`, where a message sent to the system may be delivered at a later time.
An `async` system that is deployed and live can then be sent `noop` requests to get back an enriched log of the requests handling.

The `async` methods can be tested in a `sync` way by also running them as http handlers using the `noop middleware`, 
where triggering the http method is akin to the async method receiving a message from some async-queue. 
When triggered with a `noop` context, these services can operate as usual up till the moment they have to publish a message that ends up in some queue. 
Instead, they skip publishing the message and return back the message to the test suite.

This timely response from a deployed and live aysnc-service allows testing the service without incurring the queue delays or dropped messages 
that would normally be present in live and high traffic systems. 


#### All concepts about the async testing pattern are from this GopherCon talk https://docs.microsoft.com/en-us/events/gophercon-2021/rethinking-how-we-test-our-async-architecture. 

- ❌ It does not propogate the `context` over http requests, rather a url param is used to indicate `noop`
- ❌ `Events` are not propogated with the context, which means the test suite is left with the responsibility of building up the action of events it gets back.

Building from the concepts, `noware` propogates the same `context` over http-requests and allows injecting the `action-events` into this context (opentelemetry anyone?), This enables downstream services to extract the previous event and get the expected input, without the test suite having to manipulate anything.
