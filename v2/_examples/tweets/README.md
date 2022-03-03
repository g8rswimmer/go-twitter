# Twitter v2 Tweets Examples
This directory contains examples for the APIs under `Tweets` in the Developer Platform.

## Examples
The examples can be run my providing some options, including the authorization token.

### [Tweets Lookup](https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/introduction)

* [Returns a variety of information about the tweet specified by the list of ids](./lookup/tweet-lookup)

### [Manage Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/introduction)

* [Creates a tweet on behalf of the user](./manage/tweet-create)
* [Allows a user to delete a tweet](./manage/tweet-delete)

### [Timelines](https://developer.twitter.com/en/docs/twitter-api/tweets/timelines/introduction)

* [Returns most recent tweets composed by a user](./timeline/user-tweet-timeline)
* [Returns most recent tweet mentioning a user](./timeline/user-mention-timeline)

### [Search Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/search/introduction)

* [Search for Tweets published in the last 7 days](./search/tweet-recent-search)

### [Tweet Counts](https://developer.twitter.com/en/docs/twitter-api/tweets/counts/introduction)

* [Receive a count of Tweets that match a query in the last 7 days](./counts/tweet-recent-counts)

### [Filtered Stream](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)

* [Add rules from your stream](./filtered-stream/tweet-search-stream-add-rule)
* [Delete rules from your stream](./filtered-stream/tweet-search-stream-delete-rules)
* [Retrieve your stream's rules](./filtered-stream/tweet-search-stream-rules)
* [Connect to the stream](./filtered-stream/tweet-search-stream)

### [Volume Streams](https://developer.twitter.com/en/docs/twitter-api/tweets/volume-streams/introduction)

* [Streams about 1% of all Tweets in real-time](./volume-stream/tweet-sample-stream)

### [Retweets](https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/introduction)

* [Users who have Retweeted a Tweet](./retweets/user-retweet-lookup)
* [Allows a user ID to Retweet a Tweet](./retweets/user-retweet)
* [Allows a user ID to undo a Retweet](./retweets/user-delete-retweet)

### [Likes](https://developer.twitter.com/en/docs/twitter-api/tweets/likes/introduction)

* [Users who have liked a Tweet](./likes/user-likes-lookup)
* [Tweets liked by an user](./likes/tweet-likes-lookup)
* [Allows a user ID to like a Tweet](./likes/user-like-tweet)
* [Allows a user ID to unlike a Tweet](./likes/user-unlike-tweet)

### [Hide Replies](https://developer.twitter.com/en/docs/twitter-api/tweets/hide-replies/introduction)

* [Hides or unhides a reply to a Tweet](./hide-replies/tweet-hide-replies)
