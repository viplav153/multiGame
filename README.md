# social-media

This project is simple social media

### friends urls
- /createUser -> creates new user payload require `` { "name" : "name"}``
- /sendRequest -> send friend request, only header require  [userId, friendId]
- /addFriend -> add friend from received friend request list, only header require [userId, friendId // userId is your id on which you received request and friend id is id from which you received request
- /rejectFriend -> reject friend from received friend request list, only header require [userId, friendId]
- /removeFriend -> remove friend from your friend list, only header require [userId, friendId]
- /friends -> view all  friend,  only header require [userId]
- /profile -> view you profile info, only header require [userId]


### Party urls
- /party -> creates Party, only header require  [userId] // userId of user who is going to create party
- /partyInfo -> get party by partyId, only header require  [partyId]
- /partyInvite -> sendParty invite to friend ,  only header require  [userId, friendId, partyId] // userId(any member of party),send party invite to your friend
- /partyRespond ->respond to your party Invite either ``{"ACCEPTED", "REJECTED"}``, only headers require [partyId, userId, response]
- /partyAdmin -> make other member of party Admin , only header require  [userId, friendId, partyId]// userId is current admin, friendId is friend we want to make admin
- /kick -> admin can kick out other member of party, only header require  [userId, friendId, partyId]// userId is  admin, friendId is member to kick out
- /leaveParty -> member can leave the party, only Header require [userId, partyId]


## Pre-requisite for running the project

- latest go version should be installed
- Run the docker command for installing postgres sql on port 5433
 ```shell
       docker run -d --name yugabyte -p7000:7000 -p9000:9000 -p5433:5433 -p9042:9042 -v ~/yb_data:/home/yugabyte/var yugabytedb/yugabyte:2.14.7.0-b51 bin/yugabyted start --daemon=false --ui=false
 ```
- Create Database in YSQL
  ```shell
  username@OS:/tmp/desktop$ docker exec -it yugabyte /home/yugabyte/bin/ysqlsh
  yugabyte=# CREATE DATABASE test;
  yugabyte=# exit
- run redis on docker port 6379
- ``go mod tidy``
- ``go run  main.go``
  Now all set to go !