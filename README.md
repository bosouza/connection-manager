# connection-manager

This is a school project intended as a demonstration of the factory method pattern from GoF. The idea here is that we have a connection manager (`cnnmngr.ConnManager`) that knows about clients. It offers the method `ConnManager.AddClient()` for registering new clients and associating them with `cnnmngr.ConnectionFactory`s. When someone wants to contact one of those clients, they use `ConnManager.ConnectTo(client string)` method to get a connection to the client. In order to connect to the client the `ConnManager` will use the previously registered `ConnectionFactory` to create a connection to the client, handing it back to the caller.
The `ConnectionFactory` allows decoupling the connections from the connection manager. As an example I wrote `cnnchan.ChannelCnn`, a connection that simply uses regular go channels to connect both sides, but this could easily be expanded to other kinds of connections without requiring no code change to `ConnManager`. Examples of other kinds of connections would be websocket connection, tcp socket connection, bluetooth connections, smoke signal connections...

For running the example program just run
```
go run cmd/connnmngr/main.go
```