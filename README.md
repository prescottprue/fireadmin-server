# Fireadmin Server

Golang server for storing and configuring Fireadmin accounts. Modified auth token logins are the initial functionality this server will be handling, but is planned to contain all server-side functionality for Fireadmin.

## Setup
1. Clone into your go path repository with:
  `go get github.com/prescottprue/Fireadmin`.
2. Make sure that Google Cloud Platform CLI is installed by running: `gcloud --help`.

  **NOTE**: If you don't have it installed follow the [Google Cloud CLI Quick start Guide](https://cloud.google.com/sdk/#Quick_Start)

## Development

Preview and Dev Server can be used interchangeably

### Dev Server
App engine has an equivalent to the standard go tool called `goapp`. Go app uses the python tool to create a dev server/environment.
To start the development server navigate to project root folder then run:

```bash
goapp serve .
```
More info available in [Dev Server Docs](https://cloud.google.com/appengine/docs/go/tools/devserver)

### Preview
  Previewing server on Google Cloud directly simulates running on App Engine.

  To preview the app run the following command in the main project directory (Fireadmin) folder:

  ```bash
  gcloud preview app run .

  ```

### Server Ports
  **Preview of Server**: [localhost:8080](http://localhost:8000)

  **Admin Panel**: [localhost:8000](http://localhost:8000)

  **Api Server**: [localhost:65391](http://localhost:65391)

  More info available in [Preview Docs](https://cloud.google.com/sdk/gcloud/reference/preview/).


## Deploy
To deploy the app to App Engine you can run the following command in the main project folder folder:

```bash
goapp deploy -oauth -application fireadmin-server
```
**Note**: Only approved accounts can actually deploy. [Contact Me](mailto:sprue.dev@gmail.com) if you would like to become an active developer on the project.

## Endpoints

1. ###`/setup`
  **Description**: Configures Fireadmin App for usage with other endpoints.

  **Params**:
  * `{string}` secret - Secret to associate with Fireadmin app
  * `{string}` fbUrl - Firebase url of app

2. ###`/auth`

  **Description**: Responds with auth object created using secret stored for app.

  **Params**:
  * `{string}` fbUrl - Firebase url of app

  *More Coming Soon...*
