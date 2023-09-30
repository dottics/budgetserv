## Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [0.1.1] - 2023-09-29
### Added
- The `CreateBudget` method to be able to create a new budget.
- The `UpdateBudget` method to be able to update a budget.
- The `GetGroup` method to be able to get a specific group.
- The `CreateGroup` method to be able to create a new group.
- The `UpdateGroup` method to be able to update a group.
- The `DeleteGroup` method to be able to delete a group.
- The `GetItems` method to be able to get all items. This retrieves all the
  items that are related to a group based on the group UUID.
- The `CreateItem` method to be able to create a new item.
- The `UpdateItem` method to be able to update an item.
- The `DeleteItem` method to be able to delete an item.

### Updated
- The `CreateEvent`, `UpdateEvent` and `DeleteEvent` methods to be able to
  handle the new `Event` type.

## [0.1.0] - 2023-08-18
### Added
- A `Config` type to define the configuration of the Budget microservice.

### Updated
- The `Service` and `NewService` functions to accept a `Config` type.
- The `micro` package to `v0.2.0`.

## [0.0.0] - 2022-04-15
### Added
- Service to make it a mock and testable integration for any microservice 
or API gateway.
- Models all the data structures that can be expected as responses from
the Budget microservice.
- Budget, Group, Item and Event are all exchange functions that format 
and process HTTP exchanges between the Microservice Package (MSP) and the
Budget Microservice.

## [Released]
