## Gr Hexagonal Architecture (DDD) Boilerplate ##

A production-ready, strictly typed implementation of **Hexagonal Architecture** (Ports and Adapters) using **Domain-Driven Design** (DDD) principles in Go.

This project focuses on type safety, separation of concerns, and long-term maintainability, moving away from tightly coupled dependencies towards a robust domain-centric core.

## Architectural Philosophy ##

This project was made using **Hexagonal** structure to keep a strict separation between business rules and side effects. With **Domain-Driven Design (DDD)** adherence to the core business logic was enforced. The domain becomes the source of truth not the database schema. 

### Dependency rule ###

With dependency inversion all the dependencies only flow inward. The domain layer is plain go package, with no dependency with chi, postgres or third-party libraries. So the business logic is more or less "immortal". All the frameworks may change overtime and the implementations may have breaking changes but that doesn't affect the core business logic. Also allows a plug and play like nature where any component of the project can be switched without disrupting the main business logic. 

## Technology ##

<p align="left">
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=go,postgres,docker" />
  </a>
</p>

**Core Libraries**

| Libraries | Purpose | Repo |
|-----------|---------|------|
|Chi        |Routing  |[link](https://github.com/go-chi/chi)|
|Validator  |DTO Validation| [link](https://github.com/go-playground/validator)|
|pgx        |Postgres Drover| [link](https://github.com/jackc/pgx)|
|sqlx       |Data binding from DB to domain| [link](https://github.com/jmoiron/sqlx)|
|godotenv   |ENV management | [link](https://github.com/joho/godotenv)|
|swaggo     |Swagger documentation| [link](https://github.com/swaggo/swag)|


## Progress ##

✅ Implemented | 🔄 In Progress/Planned

|Aread        |Features and Best Practices         |Status        |
|---------------|-----------------------------------|--------------|
|API Design & Architecture | RESTful API design<br> Domain Driven Design, Hexagonal architecture <br> Open API 2.0 specifications<br> Task Queue (Asynq)   |✅<br> ✅<br> ✅ <br> 🔄 |
|Database       | PostgreSQL <br> Raw SQL quries for performance <br> SQL version control and schema Migrations <br> Base ERROR maping <br> Optimized indexing | ✅ <br> ✅ <br> ✅ <br> ✅<br> ✅|
|Security       | Parameterized sql queries to prevent SQL injection <br> DTO for controlled client data<br> User input and query param validation<br> JWT-ES256 ECDSA asymmetric key pairs| ✅<br> ✅<br> ✅<br> 🔄  |
