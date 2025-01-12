# Liphium: Decentralization for everyone

Welcome to everything powering the magic experience you get in the [Liphium app](https://github.com/Liphium/chat_interface). Station powers all of Liphium's features and is it's backbone. Just as a quick note, we call all servers running station "town", because that makes explaining decentralization a lot easier. If you here "town", you can assume we mean a server running station.

Anyway, here's a quick overview over this repository:
- In the **backend** folder you'll find the service handling file management, authentication, the admin panel, friend requests and general account management.
- **Chat server** is the service handling all of the conversations and message sending.
- **Space station** is responsible for our Spaces feature where you can share ports and also our digital table.
- **Pipes** is Liphium's event loop abstraction for being able to handle decentralization and sending events without complicated code.
- **Pipeshandler** is a self-built WebSocket framework that just accepts "actions" (things sent from the client), handles them like a HTTP framework and passes back events over pipes as the client's responses. It's used for all of the WebSocket connections throughout Station.
- **Main** is the most boring folder, it's where shared types and logic between Chat server and Space station are stored and it also contains start logic to start all of these services through just one command.  

You can find a lot of information about the server [in our town documentation](https://docs.liphium.com). Be sure to check it out, and if your question isn't answered, you can always contact us through Discord or Email (available on our website).

## The goal

With Liphium, the attempt is to hide all of the decentralization magic and give the user a completely normal experience like on any other platform that is as good, if not better, than the other product. I don't want to make *another* open-source end-to-end encrypted messenger that brags about how privacy focused it is. All of that should be obvious. I want to build something people *actually* want to use because of what it is and not just because of privacy promises or the promise of digital freedom. I want to show the world how you can build a really cool and useful platform on top of a both end-to-end encrypted and privacy-focused backend while not even mentioning it to the normal if they don't want to hear about it. None of this has been achieved, but I've been building it for the last 2 years and know that it *can* be achieved. That's why I'm still working on this project after such a long time.

## This repository

Everything you can see in this repository are services that, together, power the Liphium app. For using Liphium, you need to have this installed on some VPS in the cloud or even your own home server. How to do that and what to look out for is available in [our admin documentation](https://docs.liphium.com/setup/docker.html). There is currently only a guide for installing it with Docker. But since that's the best way for now, you'll just have to get used to Docker.

## Contributing

Hey, I appreciate that you want to contribute, but I'm sorry to report that all of the guides are just not there yet. If you are really curious about the app, you can reach out to me over on our Discord server (you can find the invite on [our website](https://liphium.com)). Contribution guides and more will be available at some point, but for now I just wanna develop this app alone and make sure there is some documentation because I don't want your experience making contributions to be really bad and I also want you to know where what is and why it's that way. The code quality is also still really bad and rough in some places, I'd also like to fix all of that before allowing contributions. So, just wait for now and let me cook.
