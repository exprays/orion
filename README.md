
# Orion

An Opensource insanely-fast, in-memory database built with Go. Built upon a foundation of lightning-fast key-value pair operations, Orion goes beyond conventional databases by seamlessly supporting a wide array of data types including strings, hashmaps, arrays, and more. Whether you're handling complex data structures or requiring rapid access and retrieval, Orion empowers you with the agility and efficiency.

Currently Orion is in its active developmet phase and I am doing it in buildspace s5.



## Documentation

[Will update soon](https://linktodocumentation)


## Features

| Feature              | Status
| :------------------- | :----------: 
| In-memory KV store   | ✅     
| Strings              | ✅      
| Persistence (AOF)    | ✅  
| Sets	               | ❌
| Sorted sets	       | ❌
| Hashes	           | ❌
| Streams              | ❌
| HyperLogLogs         | ❌
| Bitmaps	           | ❌
| Persistence	       | ❌
| Pub/Sub	           | ❌
| Transactions	       | ❌
| Lua scripting	       | ❌
| LRU eviction	       | ❌
| TTL	               | ✅
| Clustering           | ❌
| Auth                 | ❌
| IQO                  | ❌

This project is being actively developed so you will see these features soon in the project.


## Run Locally

> [!IMPORTANT]  
> Make sure to have GO installed in your computer.


Clone the project

```bash
  git clone https://github.com/exprays/orion
```

Go to the project directory

```bash
  cd orion/cmd/server
```

Start the server

```bash
  go run orion.go --port=6379
```

Go to the Hunter client directory in another terminal

```bash
  cd orion/cmd/hunter
```

Start the client

```bash
go run hunter.go connect
```

> [!NOTE]  
> You can exit the client by pressing CTRL + C.

& you are ready to go. You can now store key-values in the database and also try all commands and use them efficiently.

> [!TIP]
> Orion stores AOF in an appendonly.orion file. So to when you don't need the aof, you can manually delete it to save memory.



# 💻 Tech Stack:
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
