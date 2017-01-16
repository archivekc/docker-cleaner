# Summary
Docker cleaner is a simple go command line program that find all exited container of a certain age and delete them. It then look for all unused images and remove them as well.

# Options
* -simulate : true will not remove anything only log (may list less images than real results as it does not list images of listed containers). Default is false.
* -days : number of days to wait (since exited date) before removing container. Default is 7.
* -clean-images : true will clean all unused images. Default is true.
* -frequency-seconds : the amount of second between every clean. Default is 3600 (one hour) and zero or negative number means only once.
* -filters : a set of text filter applied on image repoTags to exclude them from removal. Multiple filters can be specified using # to separate them. For example: alpine#mongo:3, define two filter: alpine and mongo:3.

# Run in a docker container
This tools is available as a docker images. It needs to have access to docker socket.

>docker run -d --privileged -v /var/run/docker.sock:/var/run/docker.sock jdfischer/docker-cleaner:latest
