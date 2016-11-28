# comment-server-go

## Design Notes
The system takes in events and commands. 
* Events are what happened in the past. 
* Commands are changes requested by users. These can fail and an error will be returned to the user. When a command succeeds, an event is created

A comment thread can be uniquely identified by the domain and title of a comment thread. 
* is the page url not that useful then?
* should a user be allowed to have the same comment thread on multiple pages?
* should a user be allowed to use the comment thread id as identifier?

There should be 2 clients interfacing with the backend. 
* An admin view where users can register domains to create comment threads on 
* A client view that is displayed on web pages displaying comments
