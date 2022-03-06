# Twitter v2 Lists Examples
This directory contains examples for the APIs under `Lists` in the Developer Platform.

## Examples
The examples can be run my providing some options, including the authorization token.

### [List Lookup](https://developer.twitter.com/en/docs/twitter-api/lists/list-lookup/introduction)

* [Lookup a specific list by ID](./lookup/list-lookup/main.go)
* [Lookup a user's owned List](./lookup/user-list-lookup/main.go)

### [Manage Lists](https://developer.twitter.com/en/docs/twitter-api/lists/manage-lists/introduction)

* [Creates a new List on behalf of an authenticated user](./manage/list-create/main.go)
* [Deletes a List the authenticated user owns](./manage/list-delete/main.go)
* [Updates the metadata for a List the authenticated user owns](./manage/list-update/main.go)

### [List Tweets Lookup](https://developer.twitter.com/en/docs/twitter-api/lists/list-tweets/introduction)

* [Lookup Tweets from a specified List](./tweets-lookup/list-tweet-lookup/main.go)

### [List Members](https://developer.twitter.com/en/docs/twitter-api/lists/list-members/introduction)

* [Returns a list of members from a specified List](./members/list-members/main.go)
* [Returns all Lists a specified user is a member of](./members/list-memberships/main.go)
* [Add a member to a List that the authenticated user owns](./members/list-add-member/main.go)
* [Removes a member from a List the authenticated user owns](./members/list-remove-member/main.go)