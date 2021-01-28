# Karbo P2P
Node it is copy of the blockchain data that shares it is data to outer world.
Karbo synchronize nodes using peer-to-peer network.
Means each node it is a client and server same time.


## Listen
First we start listening connection on defined port.
One of the main jobs done by the node, it is a handle new connections to node and response to them.

## Handshake with seed nodes
Peer must have base list of other peers, to connect to when it is starts.
This list called "seed nodes".
The list is hardcoded into specific karbo version.

On node start we must send [handshake]() request to seed nodes to get full list of available nodes.
After handshake response we may test connection to new nodes and establish connection with handshake with them. 
