# Full-Stack Go, Angular & AWS Interview Preparation Guide

This document is a comprehensive guide to help you prepare for your technical interview based on your File Upload Proof of Concept (POC) project, as well as general questions on Go, AWS, and Angular.

## 1. The Use Case of this Project

### Secure Document Management Portal (or Secure File Vault)
* **Summary:** It is a centralized, secure platform where authenticated users can upload, view, and manage sensitive files. 
* **Enterprise Application:** This POC simulates an internal company portal (e.g., HR document submission system, a financial invoice vault, or an employee resource portal). It demonstrates the ability to securely stream files via a backend API, track file metadata in a database, and ensure only authorized personnel have access via stateless JWT authentication.

## 2. Tech Stack Used (With Descriptions for Interview)

* **Frontend: Angular (v21)** - Used to build a reactive, Single Page Application (SPA). Provides a fast, dynamic user interface with route caching and state management capabilities.
* **Backend: Golang (Go)** - Chosen for extreme performance, low memory footprint, and powerful concurrency model (goroutines). Go handles multipart file streams safely and efficiently without blocking the event loop.
* **Authentication: JWT (JSON Web Tokens)** - Used for stateless, secure authorization. Once a user logs in, the Go server issues a token that Angular passes via the `Authorization` HTTP header for API requests.
* **Architecture Pattern: Clean Architecture** - The Go code is divided into `Handlers` (HTTP logic), `Usecases` (Business logic), and `Infrastructure` (Database/Storage). This allows swapping dependencies (like local storage to AWS S3) without rewriting core logic.
* **Database (POC): SQLite** - Used for rapid prototyping to track file metadata. *(Interview Tip: Explain that in AWS, this is easily swapped to **Amazon RDS / PostgreSQL**).*
* **Storage (POC): Local File System** - Used to store binary files locally. *(Interview Tip: Explain that in an AWS production environment, this would be replaced with **Amazon S3**).*

## 3. How the Project Works (Step-by-Step Overview)

1. **Authentication Flow:** The user opens the Angular app and submits login credentials. The Go API handles the `POST /login` request, verifies credentials, and returns a signed JWT.
2. **Authorized State:** Angular saves the JWT (usually in memory or secure storage) and attaches it to the HTTP Headers via an HTTP Interceptor.
3. **Upload Flow:** The user selects a file and clicks Upload. Angular sends a `POST /upload` request as `multipart/form-data`. Go intercepts it, validates the JWT, strictly streams the file to the disk (or S3), and saves the file's metadata into the database.
4. **Data Retrieval:** When viewing the dashboard, Angular sends a `GET /files` request. Go validates the JWT, queries the database for the list of files, and returns a JSON payload for Angular to render.

---

## 4. Questions Based Directly on This Project

1. **Why did you choose Golang over Node.js or Python for the backend?**
   > *Answer:* For heavy I/O operations like file streaming, Go's concurrency model (Goroutines) and static typing make it extremely performant, safer, and less resource-heavy than Node or Python.
2. **How does your JWT implementation work? What happens if the token is stolen?**
   > *Answer:* Explain stateless authentication (the server doesn't store session data). For stolen tokens, mention you would normally use short-lived access tokens combined with a refresh token, or check a Redis-based revocation list (blacklist) on every request.
3. **How are you handling CORS (Cross-Origin Resource Sharing)?**
   > *Answer:* Explain the role of your Go `corsWrapper` middleware, which attaches headers like `Access-Control-Allow-Origin` to allow the Angular app running on a different port/domain to communicate with the API securely.
4. **How does Go handle large file uploads without running out of memory?**
   > *Answer:* Go parses and streams the multipart data directly to the disk/S3 instead of loading the entire heavy file into RAM at once using packages like `multipart.Reader`.
5. **How are routes protected in Angular?**
   > *Answer:* Using Route Guards (specifically `CanActivate`). If a user attempts to access the dashboard without a valid JWT, the Guard intercepts the navigation and redirects them to the login page.

---

## 5. Golang: Basic, Intermediate & Architecture Questions

### Basic Level
1. **What are Goroutines and how do they differ from OS Threads?**
   > *Answer:* Goroutines are lightweight virtual threads managed by the Go runtime, not the OS. They start with a tiny stack (around 2KB) and grow as needed, allowing you to spawn hundreds of thousands of them efficiently.
2. **What is the `defer` keyword?**
   > *Answer:* It delays the execution of a function until the surrounding function returns. It’s typically used to ensure resources like database connections or file handles are properly closed.
3. **How does Go handle object-oriented programming (OOP)?**
   > *Answer:* Go uses composition instead of inheritance. It relies on `structs` to hold data and *implicit* `interfaces` for behavior, meaning you don't explicitly use an `implements` keyword.

### Intermediate Level
4. **Explain Channels and the difference between Unbuffered and Buffered channels.**
   > *Answer:* Channels are pipes that connect concurrent goroutines so they can safely share data. **Unbuffered channels** block the sender until the receiver is ready (synchronous). **Buffered channels** have a capacity, allowing the sender to send multiple values without blocking until the buffer is full.
5. **How does the Go Garbage Collector (GC) work?**
   > *Answer:* Go leverages a concurrent, tri-color mark-and-sweep garbage collector. It is heavily optimized for extremely low latency (pause times), meaning it sacrifices some throughput to ensure the program never freezes randomly.
6. **What is a Race Condition and how do you detect it?**
   > *Answer:* A race condition happens when two or more goroutines access the same memory concurrently, and at least one is writing. You detect this by running your code with the race detector flag: `go run -race main.go`.

### Architecture & Structure Level
7. **What is "Clean Architecture" in Go?**
   > *Answer:* It's a structural pattern that separates a project into layers (e.g., Handlers, Usecases/Services, and Repository/Infrastructure). The inner layers contain the business logic and do not know anything about the outer layers (like HTTP handlers or DB drivers). This uses dependency injection (passing interfaces) to keep the code decoupled and testable.
8. **How would you structure a production-grade enterprise Golang Application?**
   > *Answer:* 
   > * `/cmd`: Main applications (entry points).
   > * `/internal`: Private application and business logic (e.g., `/internal/handlers`, `/internal/usecases`, `/internal/repository`). Code here cannot be imported by external projects.
   > * `/pkg`: Library code that is okay to be consumed by other projects.
   > * `/api`: OpenAPI/Swagger specs or protocol definitions.
   > * `/deployments`: Dockerfiles and CI/CD configurations.
9. **How do you handle dependency injection at an architectural level without frameworks?**
   > *Answer:* Via constructor functions. For example, a `Usecase` requires a database repository to fetch data. You define an interface `Repository` in the Usecase package. Then, the `main.go` instantiates an outer DB layer (like SQLite) and passes it into `NewUsecase(dbRepo)` satisfying the interface.

---

## 6. AWS: Basic, Intermediate & Architecture Questions

> **💡 POC ARCHITECTURE MIGRATION (Crucial for passing):**
> Remind the interviewer: "I built this POC using Local Storage and SQLite for speed, but architecturally, it's designed so the `Infrastructure` layer can be instantly swapped to AWS. The local storage becomes **S3**, the SQLite DB becomes **RDS**, and the Go binary runs on an **ECS Container**."

### Basic Level
1. **What is the difference between Amazon S3 and Amazon EBS?**
   > *Answer:* S3 (Simple Storage Service) is highly scalable object storage accessed via HTTP APIs—perfect for storing user uploaded files or static websites. EBS (Elastic Block Store) is a virtual hard drive attached directly to a single EC2 instance.
2. **What are IAM Users vs. IAM Roles?**
   > *Answer:* IAM Users represent long-lived credentials (like a human developer). IAM Roles provide temporary credentials intended to be assumed by an AWS service (like an EC2 instance assuming a role to get permission to write to S3).

### Intermediate Level
3. **What is a VPC, and what are Subnets (Public vs Private)?**
   > *Answer:* A VPC (Virtual Private Cloud) is your isolated network in AWS. A **Public Subnet** has a route to the Internet Gateway (e.g., for Load Balancers). A **Private Subnet** has no direct internet access (e.g., where you securely place your RDS Database).
4. **How do you securely give your Golang application access to an S3 bucket?**
   > *Answer:* You NEVER hardcode AWS keys in your Go code (`.env` files are insecure for production). Instead, you assign an **IAM Task Role** (if using ECS) or **Instance Profile** (if using EC2) with the specific `s3:PutObject` policy. The AWS SDK in Go automatically retrieves these credentials securely at runtime.

### Architecture & Structure Level
5. **Architect a highly available and scalable environment for your Go & Angular application.**
   > *Answer:*
   > * **Frontend:** Compile the Angular app into static files and host them in an **S3 Bucket**. Put **CloudFront (CDN)** in front of it to cache it globally and provide a scalable HTTPS endpoint.
   > * **Load Balancing:** Place an **Application Load Balancer (ALB)** in public subnets to receive API traffic.
   > * **Compute (Backend):** Dockerize the Go backend and run it on **Amazon ECS with Fargate** (serverless containers). These containers sit in private subnets, completely hidden from the public internet. The ALB securely routes traffic to them.
   > * **Database:** Use **Amazon RDS (PostgreSQL)** in private subnets across multiple Availability Zones (Multi-AZ) for redundancy.
6. **How would you handle sudden massive traffic spikes (Scaling)?**
   > *Answer:* The ALB distributes the load. I would attach an **Auto Scaling Group (ASG)** to my ECS service. If average CPU utilization hits 70%, the ASG will spin up more Golang containers automatically to handle the load, and spin them down when traffic dies to save costs.
7. **How do you handle database security at the infrastructure level?**
   > *Answer:* The RDS instance lives in a private subnet. I would configure its **Security Group** to aggressively deny all inbound traffic *except* traffic originating from the specific Security Group attached to my Golang ECS containers.

---

## 7. General Angular Interview Questions (Beyond the Project)

1. **What are Angular Lifecycle hooks?**
   > Crucial hooks include `ngOnInit` (invoked once after data-bound properties are initialized), `ngOnChanges` (invoked when input properties change), and `ngOnDestroy` (used for critical cleanup to prevent memory leaks, like unsubscribing from active Observables).
2. **Explain the difference between Observables and Promises.**
   > Promises handle a single asynchronous event and cannot be cancelled. Observables (using RxJS) process continuous streams of data over time, can easily emit multiple values, can be chained with operators (map, filter), and can be cancelled manually using `.unsubscribe()`.
3. **What is Dependency Injection (DI) in Angular?**
   > A design pattern Angular relies on heavily to provide components with necessary services automatically. Instead of a component manually instantiating a class, the Angular injector provides the proper singleton instance via the constructor (e.g., `constructor(private authService: AuthService)`).
4. **What is the difference between Template-driven forms and Reactive forms?**
   > Template-driven forms utilize standard HTML directives for logic (best for simple scenarios). Reactive forms place the logic in the component class (TypeScript), which allows incredibly structured, highly scalable, and completely testable form manipulation.
5. **What are Standalone Components (Angular v14+)?**
   > They allow developers to build Angular apps entirely without `NgModules`. A component specifically outlines its own dependencies directly within the `@Component` decorator, vastly reducing boilerplate code and making the component independently reusable.

## 5. General Golang Interview Questions (Beyond the Project)

1. **What are Goroutines vs OS Threads?**
   > Goroutines are significantly lighter than standard OS threads. They are managed entirely by the Go runtime, and start with only ~2KB of stack space. You can easily and efficiently spawn hundreds of thousands of them.
2. **What is the difference between Concurrency and Parallelism?**
   > Concurrency is *dealing* with many things at once (structuring a program to run independently). Parallelism is *executing* many things at once (running simultaneously on multiple multi-core CPUs).
3. **What are Channels and when do you use them?**
   > Channels are established pipes that connect concurrent goroutines. They allow goroutines to safely pass data between them and synchronize execution without explicitly using complex locks/mutexes.
4. **Explain the `defer` keyword.**
   > It schedules a function call to be run immediately before its surrounding function returns. It is heavily used to close database connections, file handles, or network sockets predictably.
5. **How do interfaces work in Go?**
   > Go uses *implicit* interface implementation. If a struct implements all the methods defined by an interface, it automatically satisfies that interface. You do not explicitly use an `implements` keyword.
6. **How does Garbage Collection (GC) work in Go?**
   > Go primarily relies on a concurrent, tri-color mark-and-sweep garbage collector optimized heavily for minimal pause times (low latency).
7. **What is a Race Condition and how do you detect it in Go?**
   > A race condition occurs when concurrent operations access the same memory/variable simultaneously, and at least one is a write. You detect it using the `-race` flag during testing or building (e.g., `go run -race main.go`).

---

## 6. General AWS Architecture & Interview Questions

> **💡 INTERVIEW STRATEGY:** If asked how to deploy this exact project on AWS, propose this modern architecture:
> *   **Frontend:** Hosted on **Amazon S3** (Static Website Hosting), distributed globally via **Amazon CloudFront** (CDN).
> *   **Backend API:** Containerized via Docker, deployed on **Amazon ECS with Fargate** (Serverless containers, meaning you manage containers, not EC2 instances).
> *   **Database:** **Amazon RDS (PostgreSQL)** secured inside a Private Subnet.
> *   **Storage:** Backend streams uploaded files directly into an **Amazon S3 Bucket**.