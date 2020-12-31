# go-twitter v2 Examples
Within this directory there are basic examples for the twitter v2 APIs.

In order to run these examples, an authorization token will need to be provided.

## Tweet Lookup
This [example](./tweet-lookup) demonstrates the tweets lookup API call. 

```
go run *.go -token=YOUR_API_TOKEN -ids=1261326399320715264,1278347468690915330
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
