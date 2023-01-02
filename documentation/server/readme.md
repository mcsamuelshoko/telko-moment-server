# TELKO-MOMENT :  `SERVER-APP` 🖥

## About

> - nodejs server 
> - split into 2 apps
>   1. ***chat-server***
>   2. ***media-server***
> - uses REST-API (json)0
> - Uses web-sockets and uses the [STOMP protocol](https://stomp-js.github.io/stomp-websocket/)
> - server-side STOMP [documentation for nodejs](https://stomp-js.github.io/stomp-websocket/codo/extra/docs-src/Introduction.md.html)
> - Jeff Mesnil's documentation of STOMP [j.mesnil-doc](http://jmesnil.net/stomp-websocket/doc/)

> also uses **expressJs** [visit](https://expressjs.com/)
> and also **feathersJs** [visit](https://feathersjs.com/) .
> There will be an inclusion of an ORM namely ***prisma*** [visit](https://www.prisma.io/)
> code testing will be done using ***cypress*** [visit](https://www.cypress.io/)
> as for the APIs we will test them with ***postman*** [visit](https://www.postman.com/), it is also good with stomp tests on version 8+.

> ### <ul>ExpressJs</ul>
> Express Js will be the framework of choice for nodejs and will be the one implementing the ***media-server*** since the media server will not have much to deal with in terms of the databases handled and the new connections count and actions.
> Here we will also be dealing directly with the DB without the use of an ORM so as to get more of a grip of the non-ORM usage.


> ### <ul>FeathersJs</ul>
> Feathers Js will be also another framework of choice built on top of expressJs famous satisfactory results.
> feathersJs will be used on the ***chat-server*** since it will deal with most of the chats and most calls will be done here.
> Also the ***chat-server*** handles all the functions besides media, so this means user handling will also be done here.
> There also will be a very high usage of the ***prisma-orm*** so as to speed up the movement and the feature additions too.

> As for the **Database** used will have considered a NoSql database and our choice was mongodb [check out mongodb](https://www.mongodb.com/)

## NPM packages

> ### Chat-server :
>   1. FeathersJs 
>       + [visit](https://feathersjs.com/)
> 
```cmd
$ npm install -g @feathersjs/cli

$ mkdir my-app

$ cd my-app

$ feathers generate app

$ npm start
```

>   2. Prisma
>       +   [visit](https://www.prisma.io/docs/getting-started/setup-prisma/start-from-scratch/relational-databases-node-postgres)
> 
```cmd
$ mkdir hello-prisma 
$ cd hello-prism

$ npm init -y
$ npm install prisma --save-dev

$ npx prisma

$ npx prisma init
```

>   3. StompJs
>       +   [visit](https://www.npmjs.com/package/@stomp/stompjs)
```cmd
$ npm i @stomp/stompjs
```


<br/><br/>

>  ### Media-Server :
>   1. Expressjs
>       +   [visit](https://expressjs.com/)
```cmd
$ npm install express --save
```

>   2. StompJs
>       +   [visit](https://www.npmjs.com/package/@stomp/stompjs)
>       +   [visit-old-version](https://www.npmjs.com/package/stompjs)
```cmd
$ npm i @stomp/stompjs
```

