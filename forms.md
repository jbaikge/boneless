# Form Builder API & Submission Management

The Boneless CMS needs to collect form submissions which may or may not include personally identifiable information and credit card data (credit cards are processed with Authorize.net and only last four digits are stored for dispute resolution). For the sake of the target use case, these forms will capture event registrations and contact information in front of gated content.

Form processing may be its own API on its own domain supporting all Boneless CMS deployments (`forms.somesite.com`). If advised to do so, each Boneless CMS deployment may have its own form processing deployment (`forms.deployed-domain.com`).

The client requests for most, if not all, features of a WordPress plugin, [GravityForms](https://www.gravityforms.com/), which is currently used in the above scenarios. These features pertinent to this document include, but are not limited to:

1. Submission view with pagination
2. Sort via any field, e.g. last name, organization, submission date
3. Search across data in any field
4. Submission CSV export, specifying individual fields

As the Boneless CMS is currently built on serverless infrastructure with Lambda and DynamoDB, the goal for forms is also effortless horizontal scaling. While our current WordPress sites do not have problems with people rushing to register for events, we can only hope to need such capabilities. Our largest registration form to date has approximately 2,000 submissions, while an award nomination form has 11,000.

## DynamoDB

A single DynamoDB table will hold all forms and their submissions. All IDs are [xid](https://github.com/rs/xid) to allow for maintaining the same order of entry.

Each form is defined by a JSON object which is stored in a file on S3. The following information is stored in DynamoDB to make it easier to generate a list of forms in the API. Metadata is a map with a `title` key (for now).

| PK     | Metadata |
| ------ | -------- |
| FormID | Metadata |

Each submission enters into the database with the following layout; again Metadata is a map with keys like `date`, `ip`, etc. Submitted data may be a map or JSON document.

| PK                  | Metadata | Data           |
| ------------------- | -------- | -------------- |
| FormID#SubmissionID | Metadata | Submitted Data |

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

## Access Patterns

Table below taken from [identify your data access patterns](https://docs.aws.amazon.com/prescriptive-guidance/latest/dynamodb-data-modeling/step3.html) on AWS docs.

| Access Pattern | Priority | R/W | Description | Type (Single, Multiple, All) | Key | Filters | Ordering |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Create Form | Low | Write | Admin creates a new form | Single | Form ID | NA | NA |
| List Forms | Low | Read | Admin lists configured forms | All | NA | NA | Name or Date Created |
| View Form Config | Low | Read | Admin views form layout & config | Single | Form ID | Form ID | NA |
| Update Form | Low | Write | Admin updates form layout or config | Single | Form ID | Form ID | NA |
| Delete Form | Low | Write | Admin deletes form and associated entries | Multiple | Form ID | Form ID | NA |
| Form Submission | High | Write | User submits validated data from form | Single | Submission ID | NA | NA |
| List Submissions | Medium | Read | Admin lists user submissions from a form | Multiple | Form ID | Form ID | Date Created or Any field |
| Export Submissions | Low | Read | Admin exports submission data to CSV | Multiple | Form ID | Form ID | NA |
| Filter Submissions | Medium | Read | Admin filters submissions | Multiple | NA | Any field | Date Created or Any field |
| View Submission | Medium | Read | Admin views individual form submission | Single | Submission ID | Submission ID | NA |
| Edit Submission | Low | Write | Modify a user submission | Single | Submission ID | Submission ID | NA |
| Delete Submission | Low | Write | Delete a user submission | Single | Submission ID | Submission ID | NA |
| Delete Submissions | Low | Write | Delete a selection of submissions | Multiple | Submission ID | Submission IDs | NA |
