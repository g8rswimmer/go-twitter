# Twitter v2 Users Examples
This directory contains examples for the APIs under `Users` in the Developer Platform.

## Examples
The examples can be run my providing some options, including the authorization token.

### [Users Lookup](https://developer.twitter.com/en/docs/twitter-api/users/lookup/introduction)

* [Retrieve multiple users with IDs](./lookup/user-lookup/main.go)
* [Retrieve multiple users with usernames](./lookup/username-lookup/main.go)
* [Returns the information about an authorized user](./lookup/auth-user-lookup/main.go)

### [Follows](https://developer.twitter.com/en/docs/twitter-api/users/follows/introduction)

* [Lookup following of a user by ID](./follows/user-following-lookup/main.go)
* [Lookup followers of a user by ID](./follows/user-followers-lookup/main.go)
* [Allows a user ID to follow another user](./follows/user-follows/main.go)
* [Allows a user ID to unfollow another user](./follows/user-delete-follows/main.go)