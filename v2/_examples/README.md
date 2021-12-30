# go-twitter v2 Examples
Within this directory there are basic examples for the twitter v2 APIs.

In order to run these examples, an authorization token will need to be provided.

## Tweet Lookup
This [example](./tweet-lookup) demonstrates the tweets lookup API call. 

```
go run *.go -token=YOUR_API_TOKEN -ids=1261326399320715264,1278347468690915330
```

## Create Tweet
This [example](./tweet-create) demonstrates the create tweet API call. 

```
go run *.go -token=YOUR_API_TOKEN -text='Hello World'
```

## Delete Tweet
This [example](./tweet-delete) demonstrates the delete tweet API call. 

```
go run *.go -token=YOUR_API_TOKEN -id=YOUR_TWEET_TO_DELETE
```

## Recent Tweet Count
This [example](./tweet-recent-counts) demonstrates the recent tweet counts API call.
```
go run &.go -token=YOUR_API_TOKEN -query=YOUR_SEARCH_QUERY
```

## User Lookup
This [example](./user-lookup) demonstrates the users lookup API call.

```
go run *.go -token=YOUR_API_TOKEN -ids=2244994945,18080585
```

## Username Lookup
This [example](./username-lookup) demonstrates the usernames lookup API call.

```
go run *.go -token=YOUR_API_TOKEN -ids=TwitterDev,MongoDB
```

## Authorized User Lookup
This [example](./auth-user-lookup) demonstrates the authorized user lookup API call.

```
go run *.go -token=YOUR_API_TOKEN
```

## User Retweet
This [example](./user-retweet) demonstrates the user retweet API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=2244994945 -tweet_id=1228393702244134912
```

## Delete User Retweet
This [example](./user-delete-retweet) demonstrates the deleting user retweet API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=2244994945 -tweet_id=1228393702244134912
```
