# go mongo micro service

             ,_---~~~~~----._         
      _,,_,*^____      _____``*g*\"*, 
     / __/ /'     ^.  /      \ ^@q   f 
    [  @f |    +    | |  +    l  0 _/  
     \`/   \~____ / __ \_____/    \   
      |           _l__l_           I   
      }          [______]           I  
      ]            | | |            |  
      ]             ~ ~             |  
      |                            |   
       |                           |   


Gongo => Go + Mongo, refers to the Go application above that is dockerized and ran as a microservice.

An example of how to use Go, MongoDB, and Docker to create a microservice that serves an API. 
This project would be a great starter for anyone looking to write an API microservice in Go with a MongoDB database.

An example of unit testing is also provided. I have mocked the ProduceService using its IProduceService interface. The unit test impl was done levearging Testify.

The mock object can be found in the mocks folder.

It is important to note that there is no full coverage. This is left up to you. My goal was only to test the routes and provide you an example accordingly.

I have also dockerized the Gongo application. Please take a look at the docker file. Note that I use the public.ecr.aws/amazonlinux/amazonlinux:latest image.
This image is managed by AWS in the Elastic Container Registry, where you can deploy custom Gongo images. The idea then would be, upon code changes, a versioned image is deployed to a registry repository. You then leverage the image to deploy your micro service to Kubernetes or ECS (for example). One can reason then, that you do not have to use this image. An alternative I recommend is => docker pull golang:buster. Note that upon changing the image the dockerfile commands may need to change accordingly.

Levaraging container orchestration systems can be tough at times depending on your experience lvl or needs. So if you are a newbie, I recommend using Heroku, as it is a great place to experiment at low to no cost at all.

Note Ihave included two json documents with usefull docker and go commands to help you run the service.

-->> cmds.json => cmds for running go tests && cmds for building and spinning up the docker containers that make up Gongo

-->> docker-cmd.json => useful docker cmds ; If you are like me, then you like cheat sheets.

### -->> gongo.postman_collection.json => Postman collection to test the API with - this serves as the API documentaion, one can also use things such as swag to document.

--debug-instructions.json => Instructions on how to debug Gongo 

# Install
### -->> Go
### -->> Docker
### -->> VSCODE -- You do not have to use it, but I have uploaded the VSCODE launch config for debugging.
**NOTE: You do not need to install MongoDB**

# Run Gongo

## docker-compose up --build 
**This command spins up the Gongo, MongoDB, MongoExpress containers. You do not need to do anything further.**
**MongoExpress is a UI/UX that allows you to explore the data in the MongoBD instance - the one we spin up with docker-compose up.**


