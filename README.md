![](https://img.shields.io/endpoint?url=https%3A%2F%2Ftwbadges.glitch.me%2Fbadges%2Fv2)
# go-twitter
This is a go library for Twitter v2 API integration.

In order to use or demo this library, you need a developer account with twitter.  If you do not have an account, please go [here](https://developer.twitter.com/en).  Please understand that by using this library, you are using under the terms which twitter has defined.

## API Reference
The library is based off of [version 2](https://developer.twitter.com/en/docs/twitter-api/early-access) of Twitters API.  At the time of creation, the APIs are early access.  Please be aware that as the APIs are built out, the library may lag behind.

## Version 2

Currently, version 2 is in development. Please refer to [here](./v2/README.md) for more information.  Please note, that version 1 is still going to be maintained for some time after version 2 is released.

## Version 1

```
go get -u github.com/g8rswimmer/go-twitter
```

### Examples
To run all examples, the user is required to provide their developer account credentials.  The library does not share any credentials.

#### Tweet
The following examples demostrate the tweet APIs.

##### Lookup
The tweet lookup API example is located [here](./_examples/tweet/lookup).  

##### Recent Search 
The tweet recent search API example is locate [here](./_examples/tweet/recent-search).  

##### Filtered Search 
The tweet filtered search API example is locate [here](./_examples/tweet/filtered-search).  

##### Sampled Search 
The tweet sampled search API example is locate [here](./_examples/tweet/sampled-search).  

##### Hide Replies
The tweet hide replies API example is locate [here](./_examples/tweet/hide).  

#### User
The following example demostrate the user APIs

##### Lookup by ID
The user lookup API example is located [here](./_examples/user/lookup-id)

##### Lookup by User Name
The user lookup API example is located [here](./_examples/user/lookup-name)

##### Following by User Id
The user following API example is located [here](./_examples/user/following)

##### Followers by User Id
The user followers API example is located [here](./_examples/user/followers)

##### Tweet Timeline by User Id
The user tweet timeline API example is located [here](./_examples/user/tweets)

##### Mention Timeline by User Id
The user tweet timeline API example is located [here](./_examples/user/mentions)