# Distributed File Sharing System

This project is a lightweight distributed file sharing system implemented in Go. It enables users to upload, split, share, and reconstruct files over a decentralized network, inspired by the core ideas of IPFS (InterPlanetary File System). The system is designed to demonstrate key concepts such as chunking, distributed storage, and Merkle tree-based integrity verification.

## 🚀 Features

- 📦 **Fixed-size Chunking**: Files are split into equal-sized parts for parallel storage and transmission.
- 🔐 **Merkle Tree Hashing**: Ensures content integrity and verification using cryptographic hash trees.
- 🌐 **Peer-to-Peer Distribution (Prototype)**: Simulates basic peer storage and retrieval behavior.
- 🔄 **Reconstruction**: Rebuilds the original file from distributed chunks.
- ✅ **Chunk Verification**: Validates chunks against the Merkle root.
- 🌍 **REST API (Gin Framework)**: Upload, download, and verify file chunks via API endpoints.

## 🗂️ Project Structure

```
distributed-file-sharing/
│
├── api/                     # Gin-based route handlers
├── chunk_maps/              # Metadata mapping (chunk info, hash trees)
├── chunks/                  # Stored file chunks
├── downloaded_chunks/       # Chunks downloaded from distributed peers
├── output/                  # Final reconstructed files
├── reconstructed/           # Intermediate reconstruction state
├── services/                # Business logic (chunking, Merkle tree, etc.)
├── static/                  # HTML/CSS UI assets (if any)
├── testfiles/               # Sample input files
├── verify_chunks.go         # Hash verification logic
├── main.go                  # Entry point
├── go.mod / go.sum          # Go modules
└── README.md                # Project documentation
```

## ⚙️ Requirements

- Go 1.16+
- Git

## 📦 Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/altafahmed2500/distributed-file-sharing.git
   cd distributed-file-sharing
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build and run:

   ```bash
   go run main.go
   ```

   Or to build the binary:

   ```bash
   go build -o distributed-file-sharing
   ./distributed-file-sharing
   ```

4. Access the service:

   Open your browser and go to: [http://localhost:8080](http://localhost:8080)

## 🛠️ API Endpoints

| Method | Endpoint                       | Description                         |
|--------|--------------------------------|-------------------------------------|
| POST   | `/chunk`                       | Upload and chunk a file             |
| POST   | `/upload-chunk`                | Upload a specific chunk             |
| GET    | `/get-chunk`                   | Retrieve a file chunk               |
| GET    | `/download/:fileName`          | Download full file (local)          |
| GET    | `/download-distributed/:fileName` | Download using simulated distribution |
| GET    | `/merkle/:fileName`            | Generate or view Merkle Tree        |

## 🔍 Algorithms Used

1. **Fixed-Size Chunking**
   - Time Complexity: O(N)
   - Space Complexity: O(chunk size)

2. **Merkle Tree Construction**
   - Time Complexity: O(M) (M = number of chunks)
   - Space Complexity: O(M)

3. **Chunk Verification**
   - Uses SHA-256 hashing to ensure tamper-proof chunks

## 📁 Sample Workflow

1. Upload a file to `/chunk`
2. Chunks stored in `/chunks/` and metadata in `/chunk_maps/`
3. Merkle Tree is built and root hash is saved
4. Download from `/download/:fileName` or `/download-distributed/:fileName`
5. Verify chunks via `/merkle/:fileName` and `verify_chunks.go`

## 🧪 Testing

Use files in the `testfiles/` directory to try uploading, chunking, and downloading.

## 📜 License

This project is licensed under the MIT License. Feel free to use, modify, and share it.

---

## 🙌 Author

**Altaf Ahmed**  
[GitHub](https://github.com/altafahmed2500)

---

> “Data belongs to everyone. This project is a small step toward decentralization.”
