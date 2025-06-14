# Decentralized Distributed File System (DFS)

This project implements a decentralized, distributed file system using a DHT-based architecture. It allows efficient and fault-tolerant storage, retrieval, and management of files across a peer-to-peer network.

---

## Client Commands

| Command         | Description                                |
|-----------------|--------------------------------------------|
| `UPLOAD <file>` | Uploads the file to the DFS                |
| `DOWNLOAD <file>` | Downloads the file from DFS               |
| `DELETE <file>` | Deletes the file from the DFS              |

---

## Node Functionality API

These are low-level messages exchanged between DFS nodes in the network.

| Message                        | Description |
|-------------------------------|-------------|
| `PING`                        | Heartbeat message to check node liveness. Response: `PONG` |
| `unknown`                     | Sent in case of invalid or malformed requests. Response: `STOP` |
| `SEND_NODE_ID <ID> <ADDR>`   | Request to add a node to routing table (bucket) |
| `CLOSEST <NODE_ID>`          | Returns a list of closest nodes to the given `NODE_ID` |
| `STORE <hash>`               | Begins the storage protocol. Responds with `OKAY SEND`, after which the chunk is sent |
| `DONE`                       | Ends an operation. Response: `STOP` |
| `DOYOUHAVE <NODE_ID> DOWNLOAD/CHECK` | Checks if a node has the required data or else they send the chunk |
| `KEEPALIVE <NODE_ID>`        | Keeps the node active and prevents its removal from buckets |

---

## 🔧 File Storage Flow

1. **Upload**:
    - File is compressed using gzip (with deterministic settings).
    - Split into chunks, each chunk hashed (SHA-256).
    - Chunks are distributed across nodes using DHT lookup.

2. **Download**:
    - Client sends request to retrieve file.
    - DFS locates relevant nodes storing required chunks.
    - Chunks are fetched, decompressed, and reassembled.

3. **Delete**:
    - Deletes metadata and triggers a cleanup of distributed chunks.

---

## 🧪 Development Notes

- All file I/O is performed under a port-specific directory (`files/<PORT>/...`).
- Compression uses gzip with `ModTime = Unix(0)` to ensure identical outputs for same inputs.
- Communication is message-based using TCP.