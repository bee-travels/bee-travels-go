# Bee Travels Destination V2 Service - Go

The destination service is a microservice designed to provide information about various destination locations for the Bee Travels travel application.

## Data
> ***NOTE:*** All data being used is made up and used for the purpose of this demo application

The destination cities used are a subset of cities with a population over 1 million people and consists of the following data for various destination locations around the world:

* City ID
* City name
* Latitude
* Longitude
* Country
* Population
* Description about the city
* Images of the city ([Hosted on IBM Cloud Cloud Object Storage](https://www.ibm.com/cloud/object-storage))

The source of the destination service data is provided from a database. The following databases are currently supported: PostgreSQL. Check out [this](https://github.com/bee-travels/data-generator/tree/master/src/destination) for more info on data generation and populating a database with destination data.

## Environment Variables

* `PG_CONN_STRING` - postgresql standard connection string according to [LIBPG](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)

* `SERVER_ADDRESS` - Overrides the listening address (`host:port`)

## Basic Usage

* [Run](#run)
* [Test](#test)
* [Deploy to the Cloud](#deploy-to-the-cloud)

To use the car rental service navigate to the `destination-v2` directory:

```bash
git clone https://github.com/bee-travels/bee-travels-go
cd services/destination-v2/
```

### Run

The car rental service runs on port `9001`

#### Local without container

```bash
go run .
```

#### Local with container

```bash
docker build -t beetravels-go-destination-v2 .
docker run -it beetravels-go-destination-v2
```

### Deploy to the Cloud

Bee Travels currently supports deploying to the Cloud using the following configurations:

* Helm
* K8s
* OpenShift

For instructions on how to deploy the car rental service to the Cloud, check out the [config](https://github.com/bee-travels/config) repo for the Bee Travels project.