# Backend Overview

## API

| Endpoint | Action | Parameters | Description |
|----------|--------|------------|-------------|
|/expenses| GET | from=\<date\>, to=\<date\>, date=\<date\>, dates=[\<date\>], search=\<query\> , classification=\<c_id\>  unconfirmed=\<bool\> temporary=\<bool\>| Gets List of expenses as described by the parameters |
|/expenses| OPTIONS | - | Returns the list of options available for all expenses |
|/expenses/\<id\> | GET | - | Returns JSON instance of specific expense |
|/expenses/\<id\> | POST | - | Attempts to create new expense. Will merge with existing temporary expenses, and fail when exisiting non-temporary expenses already exists |
|/expenses/\<id\> | PUT | - | Replaces the expense with the instance in the payload |
|/expenses/\<id\> | PATCH | - | Replace the provided fields in the expenses |
|/expenses/\<id\> | MERGE | { id: \<id2>\, parameters: 'commission' } | Merges the expenses specificed by id1 with id2. If the parameters includes 'commission' it merges id2 as a commission expense |
