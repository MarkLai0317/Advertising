# Advertising

## tool used in this project
- MongoDB
- CircleCI
- Docker

## System Design

![system design](_assets/System_Design.png)
# Project Structure Overview

This project consists of the main package and four other key packages:

## `package ad` (Core Business Logic)

- **Domain Objects**:
  - `Advertisement`: Contains `title`, `startAt`, `endAt`, and multiple `Conditions`.
  - `Client`: Contains queries for various `conditions`, `offset`, `limit`, and a flag indicating whether each condition is missing (i.e., if the query parameter was not provided by client, it's `false`).
  - Enum types for `gender`, `country`, `platform`.

- **Interfaces**:
  - Core **UseCase interface** for the app includes creating and getting (advertising) ads (`Post`, `Get`).
  - **Repository interface**:
    - `CreateAdvertisement(Advertisement)`
    - `GetAdvertisementSlice(Client)`

- **Service Struct Implementation**:
  Implements business logic for UseCases:
  1. Validates API call parameters for compliance (e.g., `gender` cannot be "J", `age` must be between 1 to 100, etc.).
  2. Calls the injected repository method to create or get advertisements.
  3. Returns execution results.

## `package controller` (Handles Communication with External Clients)

- Isolates core business logic (`ad`) from the external communication protocol (`http`), data format (`json`).
- **DataTransferer Interface**:
  - Converts from external data formats to domain objects (`ad.Advertisement`, `ad.Client`).
  - Converts from domain objects (`[]ad.Advertisement`) to the required JSON format.

            ```json
            {
                "items": [
                    {
                        "title": "active 897",
                        "endAt": "2024-04-02T00:00:00.000Z"
                    },
                    {
                        "title": "active 121",
                        "endAt": "2024-04-05T00:00:00.000Z"
                    }
                ]
            }
            ```
        - 定義 AdvrtisementJSON, ClientJSON 處理從json轉乘
    - package repository
    - package router
