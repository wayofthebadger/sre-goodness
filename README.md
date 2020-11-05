# sre-goodness
SRE build fun times

This is my repo for the SRE funtimes test

To build this please first ensure you've got docker installed on your local system

Ensure you're in the root of this Git repo

Run the following command to build the image

docker build . -t latest

Run the next command to spawn the container

docker run -p 8080:8080 wayofthebadger/sre-goodness

Go to http://localhost:8080

Be (very very) mildy impressed/horrified at an IT system's engineer's first ever real try at a GoLang app

Many thanks
