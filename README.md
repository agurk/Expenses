# An Expenses and Document/Receipt Tracking System

## Frontend

A vue.js based front end. Located in the directory [frontend](https://github.com/agurk/Expenses/tree/master/frontend)

## Backend

The backend is written in Go and interfaces with the rest of the world via a RESTful interface.

It is located in the directory [backend](https://github.com/agurk/Expenses/tree/master/backend)

## Integration

Integration with external data sources is through the RESTful intefrace, with loaders being required to transform the data. Examples for various banks and a Doxie scanner are to be found in the [loaders](https://github.com/agurk/Expenses/tree/master/loaders) directory
