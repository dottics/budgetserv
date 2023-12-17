## Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2023-12-17

### Added

- The `SetupBudget` method to be able to set up a new budget with defaults:
  - income group.
  - expenses group.
    - norm item and category where unclassified transactions are mapped to.
    - subscriptions item and category. To show how to set up events.
  - investments group.

### Updated

- The `Item` and `Category` type and related payloads to reflect the migration
  that a category is related to a specific item.
- When a `Group` is deleted, then all the removed subgroups, items and
  categories are also deleted. However, only the removed categories are returned
  as they are required to remap transactions back to the default category. Or
  to reclassify them.

## [0.1.5] - 2023-10-27

### Added

- The `Category` type to define the category of an item.
- The `GetBudgetCategories` method to be able to get all categories.
- The `CreateCategory` method to be able to create a new category.
- The `UpdateCategory` method to be able to update a category.
- The `DeleteCategory` method to be able to delete a category.

### Updated

- The `Item` type to include the `Category` field.
  - And the methods to handle the new field.

## [0.1.4] - 2023-10-10

### Added

- The `Item`'s `category` field.
  - Allowing items to be categorised.

### Updated

- The `Group` and associated structs to omit the `GroupUUID` and `BudgetUUID`
  fields when they are not present.

## [0.1.3] - 2023-10-04

### Fixed

- The `Group` payload structs to include a `GroupUUID` field to allow groups
  to be nested.

## [0.1.2] - 2023-09-30

### Fixed

- The `GetGroups` method to also return a standard error.

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
