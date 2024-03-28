# Advertising

## tool used in this project
- MongoDB
- CircleCI
- Docker

## System Design

![system design](_assets/System_Design.png)

- 此專案中除了main package 以外有四個主要package:
    - package ad （業務邏輯核心）
        - 定義domain object
            - Advertisement : 包含title, startAt, endAt, 與多個Conditions
            - Client : 包含各個condition的查詢, offset, limit, 與各個condition 是否為missing(也就是url中有無提供該查詢參數, 沒有就為false)
            - gender, country, platform 的 enum type
        - 定義此app的核心 UseCase interface：包含建立與投放廣告(Post, Get)
        - 定義 Repository interface：
            - CreateAdvertisement(Advertisement)
            - GetAdvertisementSlice(Client)
        - 以 Service struct實作 UseCase 業務邏輯：
            1. 驗證 api call 的參數是否符合規定（ex: gender 不能是 "J", age 要在 1~100 ...）
            2. 根據 repository interface 呼叫被注入進來的 repo method，以創建或投放廣告。 (建立與投放廣告使用的repo可以不一樣，增加更多彈性)
            3. 回傳執行結果

    - package controller (負責與外界 client (who call api)溝通 )
        - 把 核心業務邏輯(ad) 與 外部使用的通訊協定(http)、資料格式(json) 做隔離
        - 定義DataTransferer interface 
            - 處理從外部資料格式轉換成domain object (ad.Advertisement, ad.Client)
            - 處理從domain object ([]ad.Advertisement)轉換成 json 需求格式
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
