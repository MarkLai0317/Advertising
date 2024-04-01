# Advertising

## tool used in this project
- MongoDB ReplicaSet
- CircleCI
- Docker, docker-compose

## Go Server Design

![Go Server design](_assets/System_Design.png)
# Go server Structure Overview

This project consists of the `main` package and four other key packages:`router`, `ad`,` controller`, `repository `


### `package ad` (Core Business Logic) 
- **folder: /ad**
- **Domain Objects**: `/ad/advertisement.go`, `/ad/country_code.go`
  - `Advertisement`: Contains `title`, `startAt`, `endAt`, and multiple `Conditions`.
  - `Client`: Contains queries for various `conditions`, `offset`, `limit`, and a flag indicating whether each condition is missing (i.e., if the query parameter was not provided by client, it's `false`).
  - Enum types for `gender`, `country`, `platform`.
    - I put country enum in contry_code.go because it's very long.

- **Interfaces**: `/ad/advertisement.go`
  - Core **`UseCase interface`** for the app includes creating and getting (advertising) ads (`Post`, `Get`).
  - **`Repository interface`**:
    - isolate domain layer from any storage infrustrucutre -> will not need to change the core business logic due to DB changes.
    - `CreateAdvertisement(Advertisement)`
        - store advertisement in any storage DB or cache as long as it implement this method.
    - `GetAdvertisementSlice(Client)`
        - query advertisements from any storage DB, cache as long as it implement this method.

- **`Service Struct`**: `/ad/service.go`
  - **Implements business logic for UseCases**:
    1. **Validates API call parameters** (e.g., `gender` cannot be "J", `age` must be between 1 to 100, etc.).
    2. Calls the injected repository method to create or get advertisements.
    3. Returns execution results.

### `package controller` (Handles Communication with External Clients like API caller for the `Ad domain`)
- **folder : /ad/controller**
- Isolates core business logic (`ad`) from the external communication protocol (`http`), data format (`json`).
- **`DataTransferer Interface`**: `/ad/controller/controller.go`
  - Converts from external data formats to domain objects (`ad.Advertisement`, `ad.Client`).
  - Converts from domain objects (`[]ad.Advertisement`) to the required JSON format for user.

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
- **`Controller struct`**: `/ad/controller/controller.go`
    - accept `UseCase interface` and `DataTransferer interface` from `package ad`for easy unit testing with Mocking
    - **GET Advertisement API** workflow
        1. call `DataTransferer interface` to convert qury data from the input format to domain objeect
        2. call `UseCase interface`, which is implemented by `service struct`, to do the business logic.
        3. call `DataTranser interface` to conveert data from domain object (slice of advertisement) to the json format body
    - **POST Advertisement API** workflow
        1. call `DataTransferer interface` to convert qury data from the input format to domain objeect
        2. call `UseCase interface`, which is implemented by `service struct`, to do the business logic.
        3. return success or error
- **`AdDataTransferer struct`**: `/ad/controller/ad_data_transferer.go`
    - implement `DataTransferer interface`
    - define `AdvrtisementJSON`, `ClientJSON` for decode json to go object with `json tag` instead of using `json tag` on `Advertisement` and `Client` defined in `/ad/advertisement`
        - By doing this, **we can decouple domain entity from data format of api**. We only need to change the controller if the input is not json or the naming of `POST API` body changes. 

### `package repository`(Handles storage operation like Redis or MongoDB for `Ad domain`)
- **folder: /ad/repository**
- **implement `Repository interface`** defined by service
- can be injected to service

- **`Mongo struct`**
    - use connection pool to perform operation on MongoDB
    - define `AdvertisementMongo struct` and `ConditionsMongo` with json tag for Marshal and Unmarshal from MongoDB
       **can decouple the naming and format of data from domain object and the input, return format of API**.
    - implement `GetAdvertisement` and `CreateAdvertisements` that is called by `Service struct` to interact with MongoDB
- **`Redis Struct`**
    - can accept a `Repository interface` 
    - **implement `GetAdvertisements`**
        - if the cache doesn't exist call the regular `GetAdvertisements method` of injected `Repository interface`
         - if the cache exist return the data in Redis.
    - **implement `CreateAdvertisement`**
        - just call the `CreateAdvertisement` method of the accepted `Repository interface`
    - due to the Accepted Repoository interface, We can **reuse** this `Redis struct` **with any other DB** implements the `Repository interface` (i.e. MySQL, PostgreSQL, Mongo, etc)
    - **with this struct, we can construct custom repository for `Service struct`**
        -  ex:  Inject the original `Mongo struct` to `Redis struct`, and inject Redis struct to `Service struct`.
            -  now we have the Repoository that work with MongoDB and has Cache capability.

            - (we can also just inject Mongo sttruct to service without using cache)

### `package router`: (Handles routing of API and middleware)
- **folder:  /internal/router**
- **`WebFramework interface`**
    - methods
        - `Get(path string, handler http.HandlerFunc)`
            - method for setting the GET route of the API
        - `Post(path string, handler http.HandlerFunc)`
            - method for setting the POST route of the API
        - `Use(pathPrefix string, middleware func(http.Handler) http.Handler)`
            - can set route for the middleware 
        - `ListenAndServe(address string, options ...ServerOption) error`
            - let `main()` call to start server 
        - `GetHandler() http.Handler`
            - return http.Handler for starting server and testing
    - With this interface, we can implement the interface with any existing web framework.
    - **ensures decoupling from any specific web framework, thereby significantly enhancing its adaptability and maintainability across various team environments.**
      
- **`ChiRouter struct`**
    - use chi webframwork to implement `WebFramework interface`

### `package main` (contruct the Advertisement api with the package above)

- define `defineAPI(adRouter router.WebFramework, adController *controller.Controller)`
    - specify the route and handler for the API
        - ex: `adRouter.Get("/api/v1/ad", adController.Advertise)`
1. construct the struct that implement `Repository interface`:       
    1. `repsitory.Mongo` 
    2. `repository.Redis` with injection of `repository.Mongo`
2. construct `ad.Service` that implement `USECASE interface`, and inject one of the above repository to `ad.Service`.
3. construct `controller.AdDataTransferer`  that implement `DataTransferer inteface`
4. construct `controller.Controller` with injection of `controller.AdDataTransferer` and 'ad.Service'.
5. contruct `router.ChiRouter` that implement `WebFramework inteface`
6. call `defineAPI(adRouter router.WebFramework, adController *controller.Controller)` that define the route and handler of api with input of `router.ChiRouter` and `controller.Controller` contstructed in step 4 and 5
7. serve GET and POST service on URL
    - http://\<host>/api/v1/ad

# Test (with CircleCI)

### Table Driven Test
- For different testcases of each UnitTest and IntegrationTest, I define functions in `test_data` folder to generate testCases that return `map[string]TestCases`. And call them in `*_test.go`
    - **example**:
        - define testcases in file
            ```
            /test_data/unit_test_data/ad/controller/ad_data_transfer_test_cases/AdvertisementSliceToJSON_test_cases.go
            ```
            ```go
            func AdvertisementSliceToJSONTestCases() map[string]AdvertisementSliceToJSONTestCase {
                testCases := map[string]AdvertisementSliceToJSONTestCase{
                    "ValidSlice": {
                        AdSlice: []ad.Advertisement{
                            {
                                Title: "test1",
                                EndAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
                            },
                            {
                                Title: "test2",
                                EndAt: time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC),
                            },
                        },
                        Expects: AdvertisementSliceToJSONExpects{
                            AdResponse: []byte(`
                                {
                                    "items":[
                                        {
                                            "title": "test1",
                                            "endAt": "2024-01-03T00:00:00.000Z"
                                        },
                                        {
                                            "title": "test2",
                                            "endAt": "2024-02-03T00:00:00.000Z"
                                        }
                                    ]
                                }
                            `),
                        },
                    },
                }
                return testCases
            }
            ```
            - the key of map is the test name and the value of map is custom struct for the test **with input and Expect result of the test**
            
        - use test cases in **test.go**
            ```
            /ad/controller/ad_data_transferer_test.go
            ```
            **run subtest for every test case**
            ```go
            func (uts *DataTransferUnitTestSuite) TestAdvertisementSliceToJSON() {
                // set stub
                patches := gomonkey.NewPatches()
                defer patches.Reset()

                patches.ApplyMethod(reflect.TypeOf(controller.ResponseTime{}), "MarshalJSON", func(rt controller.ResponseTime) ([]byte, error) {
                    t := time.Time(rt)
                    formattedTime := t.Format("2006-01-02T15:04:05.000Z") // Adjust the layout according to your requirement
                    return []byte(fmt.Sprintf(`"%s"`, formattedTime)), nil
                })

                // get test case
                testCases := data_transfer_test_cases.AdvertisementSliceToJSONTestCases()
                // run sub test with TestSuite.Run
                for name, tc := range testCases {
                    uts.Run(name, func() {
                        dt := controller.AdDataTransferer{}
                        adResponse, err := dt.AdvertisementSliceToJSON(tc.AdSlice)
                        uts.Equal(nil, err, "should not be error")

                        // unmashal []byte of json string to map for ignoring indentation of json string
                        var expectResponse, actualResponse map[string]interface{}

                        err = json.Unmarshal(tc.Expects.AdResponse, &expectResponse)
                        uts.Equal(nil, err, "expect response err should be nil: check if test file is correct")

                        err = json.Unmarshal(adResponse, &actualResponse)
                        uts.Equal(nil, err, "ad response should be nil")

                        // compare unmashal json
                        uts.Equal(expectResponse, actualResponse, "response not equal")

                    })
                }
            }
            ```
        


### UnitTest
- use `github.com/agiledragon/gomonkey/v2` for **stubbing**
- use `github.com/stretchr/testify/mock` and `mockery` for auto generating **Mock** for interface
    - I **put mocks in seperate mocks dierectory** for more clean file structure in the source code
    - the mock file for each interface in `/mocks` is at the same folder where the interface originaly defined
        - **ex**
            - `DataTransferer inteface` is defined in `/ad/controller/controller.go` : then MockDataTransferer is defined in `/mocks/ad/countroller/` folder
            -  `Repository interface` is defined in `/ad/service.go` : then MockRepository is defined in `/mocks/ad/` folder

- with the use of **stubbing** and **mocking** we can maintain the **isolation of UnitTest**

- command used to unit test in CircleCI config.yml
    ```
    gotestsum --format testname --junitfile junit.xml -- -gcflags=all=-l -v -cover  -tags=unit ./...
    ```
    - use `-gcflags=all=-l` ensure the proper execution of `github.com/agiledragon/gomonkey/v2`
    - use `-tags=unit` to only run the unit test file with
        ``` go
        //go:build unit
        // +build unit
        ```
### Integration Test
- run MongoDB container, Redis container, and the Advertisement server container for Integration test
- before running integration test command, use `dockerize` to check the above containers are running in CircleCI
- use **`TestSuite`** for **setup and tear down** of integration test
- command used to integration test in CircleCI config.yml
    ```
    gotestsum --format testname --junitfile junit.xml -- -gcflags=all=-l -p 1 -v -cover -tags=integration ./...
    ```
    - use `-tags=integration` to only run test file with
        ``` go
        //go:build integration
        // +build integration
        ```

# System Design with CQRS pattern
![system design](_assets/CQRS.png)




### **Design Reasons**: 


- the main **bottleneck of system is Querying the advertisements from DB**
- if we use only one one Collection(database),  as the number of advertisements grows, the query will be very slow
- I use two MongoDB collections to serve as Write DB and Read DB instead of two different DBMS for simplicity.
- **the Active Request is less than 1000**
    - faster query to only store these 1000 Ads in seperate collection (database if using RDBMS)  instead of query from collection with hundreds of thousands of Ads.
- **number of Post Request won't be greater than 3000 for each day**
    - the stress and the PS of the Post request isn't the main consideration in this system 



### **Database** :
- I devide the Database into **2 Part**
    - **Write DB**:  store all Advertisements whether they are active or not 
        - (can Be used as history record or used by other business logic for other system)
        - create **compound index** for **fast query of Data Synchronizer**
            ``` javaScript
            db.all_advertisement.createIndex({
                "endAt": 1,
                "startAt": 1,
            })
            ```
        - can use other DB like MySQL or PostgresSQL in the future as long as Data Synchronizer implement the interface that Query New Active Ad from the Write DB of MySQL or PostgresSQL
    - **Read DB**: **serve for the Get Advertisements service that need 10000 rps**
        - create **coumpound index**for fast query from client with different combination of conditions (`age, gender, contry, platform`)
            ```JavaScript
            db.active_advertisement.createIndex({
                "conditions.countries": 1,
                "endAt": 1,
                "conditions.ageStart": 1,
                "conditions.ageEnd": 1,
            })

            db.active_advertisement.createIndex({
                "conditionos.genders":1,
                "endAt":1,
                "conditions.ageStart":1,
                "conditions.ageEnd":1
            })
            db.active_advertisement.createIndex({
                "conditions.platforms":1,
                "endAt":1,
                "conditions.ageStart":1,
                "conditions.ageEnd":1
            })
            db.active_advertisement.createIndex({
                "endAt":1,
                "conditions.ageStart":1,
                "conditions.ageEnd":1
            })

            db.active_advertisement.createIndex(
                { "endAt": 1 }, 
                { expireAfterSeconds: 0 }
            );
            ```
    
- **ReplicaSet**:
    - I deploy **ReplicaSet of MongoDB on 3 nodes** for loadbalancing which enhance the RPS
    - use Secondaryperfer when querying
-  **influence of indexes**: 
    - my indexes follow [ESR rule](https://www.mongodb.com/docs/manual/tutorial/equality-sort-range-rule/)
    - in Addtion to the compound index above, I also use **TTL index** on `condidtions.endAt` field so that the **Ad will be deleted in ReadDB(active_collection) after `conditions.endAt < Now`**
    - **comparison** of using and no using indexes in Read DB with `wrk2`
        |  | with index | wihtout index |
        |----------|----------|----------|
        | Requests/sec | **17014.89**   |2825.15   |
### **pros and cons:**
- pros:
    - **Performance**: It allows for different optimization for the collectino for writing and collection for reading 
    - **Scalability**: It allows for allocating different resources to Commands and Queries.
- cons:
    - **Complexity**: Implementing CQRS can significantly increase the complexity of your system
    - **Data Consistency Challenges**: the new Active Ads may not be available as soon as it is insert to writeDB (is available after Data Synchronizer put Ads from WriteDB to ReadDB) 
        - **for Advertisement application, I think this trade off is worth it for the performance boost**

# use of Redis

- I originally use Redis to cache the Advertisement for a short time for the same request Parameter. However, **this design has some disadvantages** and is not sutible for this project.

    1. **the 10000 RPS usually come from different client**, Redis won't help much.
    2. **Redis will consump resource of server** when multiple different clients call the public GET API which will affect the performance of main Go server
- as a result, I choose not to use it in the end.


# performance Testing

- use [wrk2](https://github.com/giltene/wrk2) for stress testing (please download it first)

1. **step 1**
    copy the code below to  `script.lua` file to simulate multiple different request
    ```lua
    math.randomseed(os.time())

    function random_choice(choices)
        local index = math.random(1, #choices)
        return choices[index]
    end

    function request()
        local offset = math.random(1, 5)
        local limit = math.random(1, 100)
        local age = math.random(1, 100)
        local gender_choices = {"M", "F"}
        local gender = random_choice(gender_choices)
        local country_choices = {"TW", "US", "JP", "CN", "CI", "CH"}
        local country = random_choice(country_choices)
        local platform_choices = {"ios", "android", "web"}
        local platform = random_choice(platform_choices)

        local path = "/api/v1/ad?offset=" .. offset ..
                    "&limit=" .. limit ..
                    "&age=" .. age ..
                    "&gender=" .. gender ..
                    "&platform=" .. platform ..
                    "&country=" .. country

        return wrk.format(nil, path)
    end

    ```
2. **step 2**
    - run command to use the script.lua
    ```bash
    wrk2 -t20 -c1000 -d30s -R20000 -s script.lua http://35.234.12.24:80
    ```
3. The stress testing output
    ```txt
    Running 30s test @ http://35.234.12.24:80
    10 threads and 1000 connections
    Thread calibration: mean lat.: 615.429ms, rate sampling interval: 2451ms
    Thread calibration: mean lat.: 685.330ms, rate sampling interval: 2527ms
    Thread calibration: mean lat.: 600.879ms, rate sampling interval: 2406ms
    Thread calibration: mean lat.: 618.258ms, rate sampling interval: 2420ms
    Thread calibration: mean lat.: 632.469ms, rate sampling interval: 2373ms
    Thread calibration: mean lat.: 699.937ms, rate sampling interval: 2588ms
    Thread calibration: mean lat.: 605.320ms, rate sampling interval: 2428ms
    Thread calibration: mean lat.: 732.010ms, rate sampling interval: 2742ms
    Thread calibration: mean lat.: 684.973ms, rate sampling interval: 2545ms
    Thread calibration: mean lat.: 759.632ms, rate sampling interval: 2764ms
    Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     2.47s   889.61ms   5.65s    65.01%
        Req/Sec     1.75k    16.26     1.78k    71.62%
    514264 requests in 30.00s, 0.98GB read
    Requests/sec:  17141.89
    Transfer/sec:     33.62MB
    ```
4. **rps is beyond 10000 with 1000 active Ads in DB** (the inactive Ads will not influence the performance of Public API since they are stored in different collection)


# how to run this project
run the following step by step
### 0. set up environment files in main server
(host1 is the ip of your main server where you will run docker-compose-all.yml, host2 and host3 are the server ip where you will run the docker-compose-replica.yml to prepare for the MongoDB replicaSet)
- remember to change the `[host1]`, `[host2]`, `[host3]`, `[username]` and `[yourpassword]`
- `.ad_env`
    ```
    DB_URL=mongodb://[username]:[yourpassword]@[host1]:27017,[host2]:27017,[host3]:27017/?replicaSet=rs0&readPreference=secondaryPreferred
    DB_NAME=advertising
    WRITE_COLLECTION=all_advertisement
    READ_COLLECTION=active_advertisement
    DB_TIMEOUT_SECOND=10
    DB_RETRIES=3
    REDIS_HOST=redis:6379
    REDIS_POOL_SIZE=10000
    USE_CACHE=FALSE
    PORT=80
    ```
- `.sync_env`
    ```
    DB_NAME=advertising
    COMMAND_DB_URL=mongodb://[username]:[yourpassword]@mongo1:27017
    COMMAND_DB_COL_NAME=all_advertisement
    QUERY_DB_URL=mongodb://[username]:[yourpassword]@mongo1:27017
    QUERY_DB_COL_NAME=active_advertisement
    INTERVAL_SYNC_DB=5
    ```
    

### 1. docker compose
- run the `docker-compose-all.yml` in `/docker-compose` directory
    ```
    docker-compose -f docker-compose/docker-compose-all.yml up -d
    ```

### 2. MongoDB ReplicaSet setup
```
docker exec -it mongo1 mongosh -u username -p yourpassword --authenticationDatabase admin --eval  "cfg = rs.conf(); cfg.members[0].priority = 1; cfg.members[1].priority = 0.05; cfg.members[2].priority = 0.05; rs.reconfig(cfg);"
```

    











    

        





