## Consensus breaking
Two changes were introduced to potentially break consensus: 
1. Timeout configuration for CometBFT was changed  
Changing the timeout parameters used for consensus would cause nodes to not be synchronised. 
Nodes having different timeouts might break consensus as some nodes might receive more information than other nodes.

2. Changing address scheme
Without the same addressing scheme, nodes would not be able to communicate with each other to achieve consensus
