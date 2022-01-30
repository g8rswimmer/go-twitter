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

## User Retweet Lookup
This [example](./user-retweet-lookup) demonstrates looking up user that have retweeted a tweet API call.

```
go run *.go -token=YOUR_API_TOKEN -tweet_id=1228393702244134912
```

## User Blocks Lookup
This [example](./user-blocks-lookup) demonstrates looking up user that have been blocked API call.

```
go run *.go -token=YOUR_API_TOKEN -id=2244994945
```
## User Blocks
This [example](./user-blocks) demonstrates the user blocking API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=6253282 -target_user_id=2244994945
```

## User Delete Blocks
This [example](./user-delete-blocks) demonstrates the delete a user block API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=6253282 -target_user_id=2244994945
```

## User Mutes Lookup
This [example](./user-mutes-lookup) demonstrates looking up user that have been muted API call.

```
go run *.go -token=YOUR_API_TOKEN -id=2244994945
```
## User Mutes
This [example](./user-mutes) demonstrates the user muting API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=6253282 -target_user_id=2244994945
```

## User Delete Mutes
This [example](./user-delete-mutes) demonstrates the delete a user mute API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=6253282 -target_user_id=2244994945
```

## User Followers Lookup
This [example](./user-followers-lookup) demonstrates looking up user's followers API call.

```
go run *.go -token=YOUR_API_TOKEN -id=2244994945
```

## User Following Lookup
This [example](./user-following-lookup) demonstrates looking up who the user is following API call.

```
go run *.go -token=YOUR_API_TOKEN -id=2244994945
```

## User Follows
This [example](./user-follows) demonstrates the user follows API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=6253282 -target_user_id=2244994945
```

## User Delete Follows
This [example](./user-delete-follows) demonstrates the delete a user follows API call.

```
go run *.go -token=YOUR_API_TOKEN -user_id=6253282 -target_user_id=2244994945
```

## User Likes Lookup
This [example](./user-like-lookup) demonstrates user tweet likes API.

```
go run *.go -token=YOUR_API_TOKEN -id=2244994945
```

## Tweet Likes Lookup
This [example](./tweet-like-lookup) demonstrates tweet user likes API.

```
go run *.go -token=YOUR_API_TOKEN -id=2244994945
```

## Tweet Sample Stream
This [example](./tweet-sample-stream) demonstrates tweet sample stream API.

```
go run *.go -token=YOUR_API_TOKEN -output=output.txt
```

## Tweet Search Stream Add Rules
The [example](./tweet-search-stream-add-rule) demostrates tweet search stream add rules API.

```
go run *.go -token=YOUR_API_TOKEN -rule='cat has:images' -tag 'cat with images'
```

## Tweet Search Stream Rules
The [example](./tweet-search-stream-rules) demostrates tweet search stream rules API.

```
go run *.go -token=YOUR_API_TOKEN
```

## Tweet Search Stream Delete Rules
The [example](./tweet-search-stream-delete-rule) demostrates tweet search stream delete rules API.

```
go run *.go -token=YOUR_API_TOKEN -ids=14858,122345
```

## Tweet Search Stream
This [example](./tweet-search-stream) demonstrates tweet sample stream API.

```
go run *.go -token=YOUR_API_TOKEN -output=output.txt
```

## List Lookup
This [example](./list-lookup) demonstrates list lookup API.

```
go run *.go -token=YOUR_API_TOKEN -id=84839422
```
