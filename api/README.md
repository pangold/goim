#### business

It's about your business implements. Let me think about how to do with yours data.

There are two kinds of data: 

* Message: that is what you dispatch(process).
* Session: that is who you dispatch to.

##### Message

How to process the message you just received? Do you want to process it yourself? Or just pass it to another services to process it? 

There are so many ways to process it, What I implemented are:

* Grpc dispatcher: business/system/dispatch.go
* Point to point chat: business/system/chat.go

##### Session

About session, it comes from token. So how to extract your info from token, that depends on how you implement it. 

As Default, I'm using Jwt. You can just custom it.

#### codec

It uses to handle huge message for passing files, images, video. But with GRPC stream for backend service, I think it may not needs it, so far.

#### grpc

GRPC server for backend service. It can also for end-client, but not recommend.

It designs for cluster.

#### http

Http server for end-client. It's for a simple system/requirement and doesn't need backend service.

#### session

A table that stores sessions. Map user and token.

#### api.go

Entrance.