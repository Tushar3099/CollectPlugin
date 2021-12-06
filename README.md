# CollectPlugin

Pluggin Integration with Collect System Architecture

# Design Document

You can find the design document here : https://docs.google.com/document/d/1pMkGttRhusGYT7Sk961c0LYz1vnuowx-sjouUJtxf8k/edit#

# Setting up

Follow the following steps:

1. Create your `.env` file and add you mongoDB URI in the variable name `MONGODB_URI`
2. You can add your port using `PORT` env variable. Bydefault it is 3000.
3. Type `go get` to fetch all the packages required.
4. Type `make run`, your server should start.

> **_NOTE:_**
> If error pops out which says plugin and application version are inconsistent, you will have to build the plugin again using this code :
> `go build -buildmode=plugin -o actions/<ACTION_NAME>/action.so actions/<ACTION_NAME>/action.go`
