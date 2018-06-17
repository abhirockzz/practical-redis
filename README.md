# Practical Redis

> #### Head over to [Leanpub](https://leanpub.com/practical-redis) to grab a PDF version of this book

## Redis ?

[Redis](https://redis.io/) is an open source in-memory data structure server

- **Data Structure Server** - simply put, Redis is a database but differs from traditional ones since it directly exposes core data structures - strings, lists, sets, hashes, sorted sets, geospatial indexes, hyperloglogs, bitmaps etc.
	- It is possible to add new data structures and capabilities using [Redis Modules](https://redis.io/modules)
- **Other key features** - Transactions, Pub Sub messaging, Queuing, Lua scripting, ability to process infinite *Streams* of data
- **Redis also offers** - Tunable persistence mechanisms, Sentinel for High Availability and Redis Cluster for data sharding/partitioning
- **Usage patterns** include (but are not limited to) key-value database, a cache server, message broker, session store, analytics engine etc.

## About **Practical Redis**
As the name suggests, **Practical Redis** is a hands-on, code-driven guide to Redis. It's key characteristics are

**Practical** (obviously!) - Each chapter is based on an application (simple to medium complexity) which demonstrates the usage of Redis and its capabilities (data structures, modules etc.)

**Versatile** - you will get exposed to more than *just one way* of doing things

- the applications in the book are based on **Java** and [Golang](https://golang.org)
	- Redis client libraries for Java which are demonstrated - [Jedis](https://github.com/xetorthio/jedis), [Redisson](https://github.com/redisson/redisson), [Lettuce](https://lettuce.io/)
	- Redis client library for Go - [go-redis](https://github.com/go-redis/redis)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) are used to deploy the applications (along with Redis)

> The goal is to for you to run/test the applications quickly (single command) and see things in action rather than spending time setting things up

## Summary of chapters

> please note that the book is still in beta phase

Here is a quick outline of the book contents

* [Hello Redis](news-sharing-app.md) - intro to basic Redis data structures with a Java based news sharing app using the [Jedis](https://github.com/xetorthio/jedis) client
* [Extending Redis with Redis Modules](redis-modules.md) - learn about the basics of [Redis Modules](https://redis.io/modules) and make use of [ReBloom](https://github.com/RedisLabsModules/rebloom/) in a Golang based recommendation service

> ### Coming soon ...

The below chapters are work in progress on and will be made available in subsequent releases of this book

* Pipelines and Transactions in Redis
* Real time Tweet analysis service - ingest, push and analyze tweets with a combination of Java and Golang consumer services
* Scalable chat application - Use Redis PubSub and Websocket to create a chat service
* Stream processing with Redis - watch Redis Streams in action (new data structure in Redis 5.0)
* Redis based Tracking service - thanks to combination of the Geo data type and Lua scripting capability in Redis
* Redis for distributed computing - explore interesting use cases made possible by [Redisson](https://redisson.org/)
* Data partitioning in Redis - practical examples highlighting data sharding strategies in Redis
* Redis high availability


## Who is this book suitable for ?

Although I would love for this book to be used by everyone, but, it is most suitable developers who are looking to learn Redis in hands-on style

- Beginners will get a good idea of the breadth of Redis capabilities and pick up Redis faster by seeing it in action and getting their hands dirty with actual application code
- Someone who knows Redis but has not had enough hands-on experience (maybe you have been called upon to use it in a project at work ?)  
- Last but not the least, Redis ninjas might use this book as a 'refresher' if need be (although there are tons of other resources including the excellent [Redis commands documentation](https://redis.io/commands)) 

## What it is not supposed to be ?

- A theoretical discourse about the nuts and bolts of Redis e.g. it's protocol etc.
- Deep dive into administrative features of Redis
- Details of various frameworks, client libraries (e.g. Jedis, go-redis etc.) used for building the applications
- Cover extensive Docker related topics

{% creativecommons type="by-nc-nd" %}
{% endcreativecommons %}
