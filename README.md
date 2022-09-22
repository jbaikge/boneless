<p align="center">
  <img src="https://github.com/jbaikge/boneless/raw/main/assets/images/Boneless-stacked-512x512.png" alt="Boneless CMS" width="512" height="512">
</p>

# Boneless CMS

The content management system with no bones!

## Bring Your Own Definition

Upon initial installation, nothing is defined - not even regular pages. You are in control of every piece of HTML, every data field, and the content within. This CMS aims to solve the problems with people who change their mind and often want the power of a content management system but the flexibility of bespoke HTML.

## Technology Stack

Aside from the challenge of creating a CMS with no real initial structure, an additional layer came into play with trying to deploy it on AWS with serverless components. This includes:

  - __API Gateway__ to route API and frontend requests
  - __Lambda__ to handle requests from __API Gateway__
  - __DynamoDB__ for metadata storage and sorting
  - __S3__ for data and file storage
