# Decentralized Chat

## Step by Step

- Spin up zookeeper instances with `docker-compose up`
- Bash into a container instace like `docker exec -it decentralized-chat-zoo1-1 bash`
- Inside this container run `./bin/zkCli.sh -server zoo1:2181` to run the client
- On the client create the 3 initial nodes **users**, **channels** and **conn** like `create /users 0` for the **users** node
- On the root directory of the program run `make` to compile it
- Run instances of the chat like `./bin/client -username=<username> -port<port>`, make sure to have unique usernames and ports for each instance you run

## Commands

To issue a chat command you must first precede it by using the dollar sign like so `$command <params>`.
The chat client has commands to list, create, join and disconnect from a channel, those are:

- `$list`: list all channels
- `$create <channel-name>`: creates a new channel
- `$join <channel-name>`: joins and connects to all peers inside a channel
- `$leave`: leaves the current channel

## Sending Messages

To send messages to other peers inside a channel you must first join a channel. By not preceding your input with the dollar sign you will broadcast your input to all available peers inside you channel