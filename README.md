#Discussion

###Next Steps: Production Ready

1. Dockerize the application
2. Provide health check scripts to make sure the application is up and running
3. Separate the API layer and the DataBase access layer.
4. The above step could help in providing a caching abstraction which could speed up the requests.
5. I would aim for eventual consistency, hence would go for a NOSQL DB with proper indexes to make the query faster.
6. The above consideration would also help in horizontal scaling when the data increase.
7. The Database access layer would be faster when implemented as a GraphQL server, again for large workloads
8. Robust logging and analytics measurements provided by the infrastructure for debugging in case of issues.
9. Easy flex up and flex down capabilities which will be facilitated by dockerizing the app.
10. Make sure the app is Disaster recovery ready, by deploying it in two or more isolated geo locations with the help of a robust load balancing mechanism.


###Performance, Reliability, Maintainability ... Limitations

1. I would prefer to have a much better logging mechanism middleware [I'm currently using the stock option]
2. The application is single threaded, would prefer to have a connection pool mechanism and a concurrency manager to serve large number of requests.
3. Swagger implementation would provide a good mechanism for integration with the API layer for frontend or other clients using the application.
4. Provide more documentation and have a regression/functional test suites.
5. Provide metrics on code coverage and other static analysis reports.
6. Authentication mechanism to make sure POST calls are authenticated and data is not mutated in an uncontrollable way.


###APIs exposed

1. /engineers --> Gives you the list of distinct engineers from the DB.
2. /actions/engineer/{engineer} --> Gives you the entire list of actions done by the engineer from the DB. {engineer} is the variable
3. /events/from/{from}/to/{to} --> Gives the list of actions that was performed `from` to `to` duration, inclusive, only accepts `unix timestamp`
4. /eventStats/{time} --> Gives the stats of the different events that was done during the day within which `time` falls under
**eg: /eventStats/1509192354  would give you the stats for October 28th 2017, since 1509192354 falls on that date**
5. /eventStats --> Would give the entire list of stats from the DB

P.S. For high volumes data, this should be paginated so that we could still keep processing the request on demand. 