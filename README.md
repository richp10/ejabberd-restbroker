# ejabberd-restbroker

This is a 'broker' which accepts jobs via Redis using sidekiq protocol and uses these to construct RESTful 
commands to interact with the ejabberd API. 

We use use goroutines for concurrent processing of jobs, so a high volume of messages can be processed very quickly. The app also re-uses the http connection with ejabberd, which is more efficient than re-creating this for every request. It increases resilience by allowing commands to be buffered and retried if they fail.  It also allows the sending applications to 'fire and forget', without waiting for the response from Rest. Because it uses the sidekiq protocol, any sidekiq compatible library can be used for the 'sender'.   

Optionally, a request can return the response from ejabberd via a different Redis queue. This potentially allows async processing of the restful interaction with ejabberd.     

### Configuration
Create env.json with the mandatory settings: 

- "Token" - security token required with all requests 
- "JabberDomain" - domain of the jabber service
- "RestUser" - jabber account with permission to use rest  
- "RestPass" - password for this account
- "RestQueue" - name of the redis queue to get messages

Other settings have sane defaults and _might_ work unchanged (eg. via localhost)

- "RestURL" - url to the rest endpoint for jabber
- "RedisHost" - host url and port for redis
- "ResponseQueue" - Queue name to post responses if indicated
- "MaxIdleConnections" - for http connection
- "RequestTimeout" - for http connection
- "RedisHost" - optionally including port, defaults localhost:6379,
- "RedisDatabase" - defaults to 0,
- "RedisPoolCount" - connection pool to redis, defaults 15,
- "RedisNamespace" - usually best not to use, defaults ""
- "ProcessID" - allows worker id if restarted, defaults 1
- "Concurrency" - how many workers to run at the same time, defaults 20
- "Stats" - display api with summary of worker statistics in json, defaults false
- "StatsPort" - if Stats is true, which http port for the result, defaults 8080 

### Requests
Clients add messages to the redis queue with three mandatory and a fourth optional parameter.  

- The first param is assumed to be the endpoint, added to the end of the rest url. This corresponds to the
ejabberdctl command. 

- The second param is assumed to contain any arguments that would be passed to the ejabberd command. This must be an array of key value pairs, which is then converted to json before being passed to ejabberd (as the body POSTed to the endpoint). 

- The third param is assumed to contain the security token.  We check this matches the configured token as a security check before running the command. 

- Optionally, the fourth param is an id for the response. If this is provided, any response from Ejabberd will be sent to the configured ResponseQueue with this id, providing an async response channel. 

### Security
This app potentially opens up a big holes in the security of your ejabberd server.  Make sure you have configured ejabberd to only expose commands that are essential and consider how best to secure the token that is used to authorize messages.

### Posting Requests
Any sidekiq compatible library can be used to post messages to the job queue. For example, if your need was PHP you could use: https://github.com/spinx/sidekiq-job-php 

    $redis = new Predis\Client('tcp://127.0.0.1:6379');
    $client = new \SidekiqJob\Client($redis);
    $token = 'security-token-random';
    
    $endpoint = 'user_sessions_info';
    $args = [
        'user' => 'richard',
        'host' => 'jabberdomain.org'
    ];
    
    $client->push('Add', [$endpoint, $args, $token, 'responseQueuename' ], false, 'queuename');

#### Installation
Uses golang deps:   

    1) go get -u github.com/golang/dep/cmd/dep
    2) dep init
    3) dep ensure

#### Todo
- No tests yet
- Response queue has not been tested 

#### Mea Culpa
This is my first production golang app and the coding is likely a bit crap. By all means fork and submit PR's... 

#### Thanks to..
The backbone of the app is https://github.com/jrallison/go-workers
