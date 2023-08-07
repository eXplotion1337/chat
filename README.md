# go-chat

See the improvised version [here](https://github.com/alextanhongpin/go-chat.v2).



## Requirements

**Functional requirements:**

- system should allow users to have a private conversation
- system should allow users to have a group conversation
- system should not expose the user identifier
- system should notify users in group when member is online
- system should notify users in group when member is offline
- system should notify users when message is read
- system should display count of unread messages
- system should mark messages as read
- system should allow only authenticated users to chat
- system should allow users to login with multiple sessions
- system should handle stickiness of the user's session

**Non-functional requirements**

- system should be reliable (messages should be stored)
- system should be available (users should not be disconnected)
- system should be secure (authentication)

**Extended requirements**
- system should handle validation on the messages sent to avoid spam
- system should allow users to block unwanted chat requests
- system should send notification when user is not online

## References
- https://github.com/alextanhongpin/go-chat
- https://www.thepolyglotdeveloper.com/2016/12/create-real-time-chat-app-golang-angular-2-websockets/
- https://devcenter.heroku.com/articles/go-websockets
- https://www.jonathan-petitcolas.com/2015/01/27/playing-with-websockets-in-go.html
- https://blog.arnellebalane.com/sending-data-across-different-browser-tabs-6225daac93ec
