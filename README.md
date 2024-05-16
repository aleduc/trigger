## Task
Your team has been tasked with building a server that initiates and maintains
outbound phone calls. For example, it would be used in an application where

users requested a call-back from their utility company, or implementing reach-
out from a GP clinic for scheduling in patients who signed up for their flu jabs.

As part of this server, you have been tasked with implementing the “call trigger”
functionality. The server will expose an endpoint
POST /trigger which clients can call to initiate an outbound call. The
endpoint takes the following JSON payload:
{
"phone_number": "+44788888888",
"virtual_agent_id": "TTFD_UDFNuhdeuhUHUwd"
}
This API will respond with a call id, which can then be used to check the status
of the call. (This is
out of scope — but just giving an idea of the usage patterns).

To initiate the outbound calls we rely on a third-party API POST /originate_call
that makes the phone call between the client and the virtual agent. If we make a
valid request to this API, then it may return the following responses:

* 200 OK - the call was initiated and successfully answered.
* 429 Too Many Requests - the call was rejected by the API, as we’re being
rate-limited by the API. It will not return a Retry-After header.
The API rate limit is set at 25 API requests every 10 seconds.
Hitting the rate limit will introduce a substantial backoff (e.g. 30s)
The /originate_call API latency is not negligible. It can take up to 5-10s for a
call to be established (think about how long it takes for a call to ring before
pickup). We want our API wrapper to respond much faster than that.

Your objective is to:

* Implement the trigger route, together with tests. The signature of
the function you’re expected to provide is essentially:
function trigger(phone_number, virtual_agent_id) → call_id
* Maximise call throughput while minimising the number of rate limit
errors from the third-party API.
* This API should respond instantly, it shouldn’t block and it should
not return rate limiting errors. The status of the call will be queried
through a different route (this other route is out of scope)

## Principles
* Avoid third-party libs, if it doesn't consume too much time.
* Error handling is kind of ignored(random messages).
* Somewhere comments don't follow standards. In obvious places skipped, in tough places too large.
* Naming doesn't take into account future code.
* Time-consuming, premature optimizations were skipped.


## CMD
**cmd/test_trigger/main.go** - main run, with "google" URL for API.
**cmd/mocked_trigger/main.go** - same as main, but with mocked externalAPI call.

## General description
This implementation is based on producer/consumer pattern(pub/sub) + worker pool.

**Handler** - is producer.

**Storage** - shared storage for producer/consumer.

**Worker** - is consumer.

**Pool** - implements worker pool.


## Additional packages
realtime, logger - auxiliary packages, useful for tests.
http_wrapper - is not the best name, simple http wrapper for requests.

## Limiter
Sliding window limiter: I've tried to implement it with a circular/ring buffer. 
https://en.wikipedia.org/wiki/Circular_buffer.

I could use a timer for cleaning values each second, but ideally, the timer should be adjusted each run.
It makes sense for minutes, for seconds can be overwhelming. 

**Example**:

current time: 14:00:01.103
next second will be ~ 14:00:02.103


## TODO, ways to improve
First of all, this task should be implemented in 2 services + message broker + storage.

**first service**: Handle request, saves/produces messages to broker/queue.

**second service**: Consume and process messages from a queue.

**storage**: Some RDBMS/ NoSQL database, depending on purposes.

**message broker**: Kafka, PubSub(GCP), Apache ActiveMQ, RabbitMQ, etc

**Also:**
* use https://github.com/mailru/easyjson for static json.Marshal/unmarshal 
* use https://github.com/valyala/fasthttp for requesting/handling (should be tested)
* https://pkg.go.dev/golang.org/x/sync/errgroup can be used for graceful shutdown, but it is a matter of taste.



