#Get A Redis

A one click, docker based, auto scaling, Redis host implemented in Go and hosted on Digitalocean.

![Homepage](https://raw.githubusercontent.com/MohamedBassem/getaredis/master/imgs/GetARedisHomePage.gif)

##Tags

- Docker
- Redis
- Digitalocean
- Service Discovery
- Auto Scalability

##Why?

I started the project to enhance my Go skills which I started to learn few weeks ago. Then I found that the idea may be useful for hackathons and proof of concept projects.

##Technical Details

###System Components

####Go Server

Running [martini](https://github.com/go-martini/martini) to serve the webpage and accepts starting new instance requests. Nginx is installed on this machine to act as a reverse proxy for martini.

####Go Jobs

Currently two jobs are scheduled to run every certain amount of time. The first one is a job to kill the containers that has been running for a preconfigured amount of time. The second one is a job to spin up and tear down digitalocean droplets based on the current load of the other servers.

####Docker Hosts

Digitalocean droplets that are used to host redis containers. Nginx is installed on those machines to act as a reverse proxy to the running docker daemon with HTTP authentication.

####Redis

Redis is used only for service discovery. Service discovery will be explained in a later section.

####MySQL Database

To store the details of the running containers, such as the container host, port, id, state and the creator IP. All of the containers' information could be collected from running docker hosts but the database is mainly used for throttling the number of containers per IP.

###System Architecture

![System Architecture](https://raw.githubusercontent.com/MohamedBassem/getaredis/master/imgs/SystemArchitecture.png)

The server, running nginx, listens for requests on port 80 and forwards those requests to Go running [Martini](https://github.com/go-martini/martini) as a web framework. Go queries redis, which will be explained later, for active docker hosts. Go then tries to schedule the new container on one of the hosts based on a certain criteria. The current criteria is to try to schedule as much containers as possible on the host to reduce the running costs, since it's currently a free service :grimacing:. The maximum number of containers per host is configurable. The server then schedules the container on the chosen host and insert these data into the database. The details of the scheduled container (host, port, redis password) are returned back to the user. The user can then connect to the redis container directly. A background job runs every 20 minutes to kill containers that have passed their maximum allowed number of hours, currently 12 hours.

####Auto Scaling

Another job runs in the background every 10 minutes to check the load of the docker hosts. Since the containers have a maximum memory of 5MB, the load is estimated by the number of containers running on each host. As we mentioned before scheduling tries to add more containers to the busiest hosts as long as they can hold more. Another host is waiting in a standby state. Whenever the standby host starts getting some containers, the job will start another container. Whenever the job detects that there are more that one host containing no containers, it kills them until one is left.

####Service Discovery

Service discovery is needed because we have automatically scaling hosts that the scheduler needs to detect. The are many tools that can be used for service discovery, such as Apache Zookeeper, etcd and Consul. I needed a very simple discovery service so I decided to implement my own.

Redis has a command to expire some key after a certain amount of time. Whenever this command is called it resets the timeout. Using this idea, docker hosts can add a key for themselves in redis and constantly refreshes the timeout. If the key times out, this means that the host didn't send a heartbeat which means that it got disconnected.

The code is as simple as this:
```bash
#!/bin/bash
(
 PRIVATE_IP=$(curl http://169.254.169.254/metadata/v1/interfaces/private/0/ipv4/address)
 echo "AUTH <REDIS_PASSWORD>";
 while true; do
 NUMBER_OF_CONTAINERS=$(($(docker ps | wc -l) - 1))
 echo "SET server:$NODE_NAME '{\"PrivateIP\":\"$PRIVATE_IP\",\"NumberOfContainers\":$NUMBER_OF_CONTAINERS}'";
 echo "EXPIRE server:$NODE_NAME 10";
 sleep 4;
 done
 ) | telnet REDIS_IP REDIS_PORT
```

##TODO
- ~~Open Docker port on hosts for the master.~~
- ~~Pull redis image on new hosts.~~
- ~~Authenticate redis master when connecting to docker hosts.~~
- Supporting more containers other than redis.
- Monitoring console.
- Better deployment method.
- Better documentation.

##Note
The project is still in beta and not 100% stable.

##Contribution
Your contributions and ideas are welcomed through issues and pull requests.

##License
Copyright (c) 2015, Mohamed Bassem. (MIT License)

See LICENSE for more info.

