// this example have one endpoint /register

// when our server is run board maker reciver listen on topic="user.register"
// in user register api after registered user we publish it in user.register topic

// in board maker package give message from gochannelpubsub and start to make board for user
package main
