First time using go, so it's a bit slow going. 200 lines of code ain't much.

Based on initial impression Go takes care of 90% of the work a server needs to do. handle requests, format the info sent/received goroutines. I think this will be either a good basic lesson to learn what standard functions Go uses or just a template to know what other points to learn for a lower level integration. Would be good practice for string formatting.

Like... What is a server loop looking like? With Go handling so much you don't really learn any of that. What does it look like to wait for a request? Don't want it running full tilt like a game loop, so is it waiting for input something like a TUI? Does the kernel handle sending the messages that something new came into the network? Would it change based on hardware or OS?

As a template for what a server could need to handle it is a nice starting point.
Get
Post
data formatting
querying another DB
Bundling
Webhooks
Authentication

Breaking it down like this makes it easier to see where the different points of 'task completed' are for another lower level implementation.

Quick glance at Odin's net package and bindings to OS seems like a net package is a wrapper on the OS's network APIs.
https://github.com/odin-lang/Odin/blob/master/core/sys/windows/ws2_32.odin#L199
https://learn.microsoft.com/en-us/windows/win32/api/winsock2/nf-winsock2-listen