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



