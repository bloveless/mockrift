# Mockrift

What is mockrift. Contract testing for developers and CI/CD environments.

Turn this into a sales pitch.

---

Essentially it is a contract testing/mocking framework for systems that require use of external/internal APIs.
So when you are a developer you can turn on mockrift and put it in proxy mode. It will record all the requests you send to it and all the responses that come back from the backend. It can then do matching based on method, url, body, and headers and return responses in leu of the backend.

So as a developer you can build interesting responses that test specific scenarios in your UI and just have mockrift replay those scenarios.

It can record multiple responses for every request and the developer can just pick the active one and it will serve up the active response over and over again.

Think of Younique when we had to use postman to insert garbage into the request so we could test errors and stuff like that. Well in mockrift you can just record it, share it, and replay it.

Then there is the CI/CD part. You can use those responses and convert them into contracts. Mockrift will play requests to your backend and verify that the response matches the contract. That part is the headless part. You’d run mockrift in your CI/CD environment and tell it to play a scenario. It would hit your backend with a suite of requests and verify that every response that comes back exactly matches your contract. Other wise the test fails.

The idea is that you’d be able to commit your contracts as something like a git submodule within both projects (say you are building two apis that communicate) then mockrift can function on both ends to verify the neither project has broken their contract.

And neither project would ever actually have to talk to each other.

It is also built as a single docker container with no external dependencies (I.E. MySQL or Mongo), so you can include it super easily.

You are welcome to run the project using straight go if you’d like. Maybe as a zip/binary that can just be run to make it easier... not sure about that one.

---

Companion app for testing [Shopping Cart Demo](https://github.com/bloveless/demo-shopping-cart)
GraphQL first then REST later.

ToDo For Standard Development:
- [x] Proxy Requests to backend
- [x] Store/respond with Requests/Responses from backend and serve those responses
- [ ] Save requests/responses to the file system and reload them when mockrift starts
- [ ] Store multiple requests/responses from backend and allow user to switch between responses. (1)
- [ ] Create a UI that can be used by developers to pick which response to server from a request on-the-fly.
- [ ] Create an algorithm for comparing stored requests which will create a match percentage and match requests by their
highest match percentage (2)
- [ ] Allow the user to duplicate/edit stored responses from the UI
- [ ] All exceptions are log.Fatal right now. Figure out the best way to handle exceptions rather than dying.

1) In the beginning the user will be able to pass in a specific header in order to match responses with requests.
Eventually there will be a UI that developers will be able to use, but the header will continue in order to configure
requests for CI/CD.
_The header may not even need to be anything special since it the header will just be used as part of the comparison
when searching for a response to return._

2) For example a request that matches a URL, Method, Body exactly may return a 95% match because some of the headers are
missing or different. Different requests would likely have different criteria since some requests may be more concerned
about specific things I.E. headers may be very important for a match. Or only specific headers may be very important.
Anyway, there should be some algorithm for matching requests to responses based on URL, Method, Body, and Headers.

ToDo For Contract Testing:
- [ ] Mockrift should be able to do shape deviation testing (1)
- [ ] Mockrift may need to store test plans in order to do CI/CD. I.E. Send request A to backend (verify contract) then
send Request B to the backend (verify contract)
- [ ] Mockrift should be able to be called from the CLI with a test plan and execute the plan without starting the server

1) Shape deviation testing is when the shape of the body (request or response) changes in some way. Maybe there is a key
where there shouldn't be one. Maybe a value was once and integer and is now a string. These would fail the contract
test.
