<p align="center">
<image width="150" height="50" src="images/kube-go.png"></image>&nbsp;
<image width="260" height="410" src="images/golang-logo.png">
&nbsp;<image width="130" height="50" src="images/tls-logo.png"></image>
</p><br/>
<br/>

# Go TCP Server/Client Plugin Modules

It defines modules used as default provy modules from Go TCP Server and clients

## How does it work a module

Module is composed by 2 parts:

* Server side, implementing [Commander](https://github.com/hellgate75/go-tcp-server/blob/master/common/types.go) interface

* Client side, implementing [Sender](https://github.com/hellgate75/go-tcp-client/blob/master/common/types.go) interface


## How to develop an external Linux Plugin Module


External Linux plugins must implement proxy function, for:

* Server side: [GetCommander](https://github.com/hellgate75/go-tcp-server/blob/master/server/proxy/proxy.go)

* Client side: [GetSender and Help](https://github.com/hellgate75/go-tcp-client/blob/master/client/proxy/proxy.go)

Some implementation are available,

For the client:

* [shell plugin](/client/proxy/shell/shell.go)
* [tranfer plugin](/client/proxy/transfer/transfer.go)

For the server:

* [shell plugin](/server/proxy/shell/shell.go)
* [tranfer plugin](/server/proxy/transfer/transfer.go)


## Rules

Any of the plugins must incorporate code for client, server or client and server together.
* Client plugin will be responsible of communication with server.
* Client plugin will be responsible of executing commands on the server.

In the GetSender and GetCommander you define match between command name and developed components.

Connectivity and other features are available in the repositories :

* [go-tcp-common](https://github.com/hellgate75/go-tcp-common)

## References

Here list of linked repositories:

[Server Repository](https://github.com/hellgate75/go-tcp-server)

[Client Repository](https://github.com/hellgate75/go-tcp-client)



Enjoy the experience.



## License

The library is licensed with [LGPL v. 3.0](/LICENSE) clauses, with prior authorization of author before any production or commercial use. Use of this library or any extension is prohibited due to high risk of damages due to improper use. No warranty is provided for improper or unauthorized use of this library or any implementation.

Any request can be prompted to the author [Fabrizio Torelli](https://www.linkedin.com/in/fabriziotorelli) at the following email address:

[hellgate75@gmail.com](mailto:hellgate75@gmail.com)


