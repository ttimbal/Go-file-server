# Go-file-server
TCP server to transfer files between clients.
The clients can subscribe to channels to send and receive files.


The client and server use a custom CLI.

### Server
Enter the server folder and follow the next steps.
  * Run
  ```
  go run .
  ```
  * Start
  ```
  server start
  ```
  * Stop
  ```
  server stop
  ```
    


### Client
Enter the client folder and follow the next steps.
  * Run
  ```
  go run .
  ```
  * Subscribe channel
  ```
  subscribe channel:name
  ```
  * Unsubscribe channel
  ```
  unsubscribe channel:name
  ```
  * Send file
  ```
  send channel:name file:path
  ```  
