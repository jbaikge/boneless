# Form Builder API & Submission Management

The Bonless CMS needs to collect form submissions which may or may not include personally identifiable information and credit card data (credit cards are processed with Authorize.net and only last four digits are stored for dispute resolution). For the sake of the target use case, these forms will capture event registrations and contact information in front of gated content.

Form processing may be its own API on its own domain supporting all Boneless CMS deployments. If advised to do so, each Boneless CMS deployment may have its own form processing deployment.

The client requests for most, if not all, features of a WordPress plugin, [GravityForms](https://www.gravityforms.com/), which is currently used in the above scenarios. These features pertinent to this document include, but are not limited to:

1. Submission view with pagination
2. Sort via any field, e.g. last name, organization, submission date
3. Search across data in any field
4. Submission CSV export, specifying individual fields

As the Boneless CMS is currently built on serverless infrastructure with Lambda and DynamoDB, the goal for forms is also effortless horizontal scaling. While our current scenario does not have problems with people rushing to register for events, we can only hope to need such capabilities. Our largest registration form to date has approximately 3,000 submissions.

## DynamoDB

Each submission enters into DynamoDB with the following layout:

| PK      | SK            | Metadata | Data           |
|---------|---------------|----------|----------------|
| Form ID | Submission ID | Metadata | Submitted Data |

All IDs are [xid](https://github.com/rs/xid) to allow for maintaining the same order of entry. Submitted data may be a map or JSON document.

Satisfying the requirements requires some trade-offs, though:

1. Pagination with hard per-page values involves scanning up to the last viewable entry before returning the result subset
2. Sorting requires unmarshalling data for every entry and then sorting by the field(s) requested in memory
3. Searching, like sorting, will require iterating over each unmarshalled entry to find matches, then sorting the results; or querying a copy of entries in an indexing service (Open Search)
4. Export is very straightforward

## SQLite on S3

Each submission is stored as a JSON file in an S3 bucket. A recurring Lambda function reconciles all of the JSON documents into a single SQLite file on a regular basis. There are a few different ways to lay the database out - either with a single table or two tables (entry & attributes).

1. Pagination available with `LIMIT` and `OFFSET`
2. Sort available with `ORDER BY`
3. Search available with full text indexes available inside SQLite
4. Export is very straightforward

## Relational Database

The engine does not matter, but the cost does. Some of the sites Boneless CMS will support may be microsites that stand up, take some traffic, take some event registrations, then sit idle as an archive of the event the site supported.

1. Pagination available with `LIMIT` and `OFFSET`
2. Sort available with `ORDER BY`
3. Search available with full text indexes or rudimentary `LIKE` clauses
4. Export is very straightforward
