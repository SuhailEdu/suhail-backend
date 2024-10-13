## README: Suhail Backend

### Overview

**Suhail** is a built to be robust and scalable platform designed to facilitate the creation, participation, and sharing of E-exams with students.

### Technology Stack

  - **Golang** 
  - **PostgreSQL** 
  - **pgx** 
  - **Go Melody** 
  - **Echo Router** 
  - **sqlc**

    ### This repository is designed to act as a backend for our [Suhail-Frontend](https://github.com/SuhailEdu/suhail-frontend) repository.

### Installation and Usage

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/SuhailEdu/suhail-backend
    ```
2.  **Set up Environment Variables:**
    Copy the `.env.example` to `.env` file and set the necessary environment variables.
3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```
4.  **Run the Application:**
    ```bash
    go run ./cmd/api
    ```