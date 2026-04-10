## How to Run This Project (Local Demo)

Follow these steps to demonstrate the POC during your interview:

### 1. Run the Backend (Golang)
Open your terminal in the root directory and run:
```powershell
go run cmd/api/main.go
```
*   **API URL:** `http://localhost:8080`
*   **Action:** This starts the server and initializes the SQLite database (`files.db`).

### 2. Run the Frontend (Angular)
Open a **second** terminal, navigate to the `web` folder:
```powershell
cd web
npm install
npm start
```
*   **Frontend URL:** `http://localhost:4200`
*   **Action:** This serves the Angular UI.

### 3. Running with Docker
If you'd like to show containerization:
```powershell
docker-compose up --build
```

---

## 9. Project Overview & Working Logic (Short Points)

Here is a summary you can give at the end of your demo:

*   **Objective:** Provides a secure, scalable vault for uploading and managing documents.
*   **Architecture:** Follows **Clean Architecture**, making it highly maintainable and decoupling business logic from external frameworks.
*   **Auth Flow:** Implements **JWT (JSON Web Tokens)** for secure, stateless communication between Angular and Go.
*   **File Streaming:** The Go backend uses efficient multipart streaming to handle file uploads without high RAM consumption.
*   **Security:** Includes Angular **Route Guards** and Go **Auth Middleware** to ensure only logged-in users access the data.
*   **Database:** Uses **SQLite** for metadata persistence, ensuring file records are kept between restarts.
*   **AWS Compatibility:** Designed with interfaces so it can be deployed to **AWS S3, RDS, and ECS** with minimal code changes.