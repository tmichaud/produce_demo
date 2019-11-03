# Overview
The Produce Demo API provides a RESTful interface to produce inventory.  Produce items can be fetched, added or deleted.  This API is written in Go and is stored in Github.  
It is designed to be run in a Docker container and will be built and tested automatically via CIRCLECI.  



# Infrastructure
This code resides in Github. When the master branch is updated, CIRCLECI will pull the master branch, build, and run all unit tests via the Dockerfile.  Any error that occurs during building will stop the build and the developer(s) will be alerted via email.  If the build succeeds, developer(s) will be notified via email and the container will be pushed to Dockerhub tmichaud/produce_demo and tagged with both the github annotation tag and with "latest".

It is the responsibility of developer(s) to ensure that the code is always appropriately tagged. Note that the annotated tag should be using semantic versioning (https://semver.org/).  GitFlow (https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) is assumed to be used and work will be done on feature branches.  The annotated tag will be merged into master automatically when the feature branch is merged to master.

```
To create an annotated tag on a feature branch:
    git tag -a '<major>.<minor>.<build>' -m "<comment>"   
    git push origin '<major>.<minor>.<build>'
where major.minor.build is the next version in semantic versioning.
```

All other branches besides master are ignored by CIRCLECI. 

The Docker build for this API uses a 2 stage build. The 2nd stage utilizes the FROM SCRATCH which allows for a smaller container to be built but developers cannot log into the container.  

# How to run Produce Demo

Log into DockerHub with :  \
```	docker login ```

Pull Docker image:  \
```	docker pull tmichaud/produce_demo  ```

Run Docker:  \
``` 	docker run -p 8080:8080 tmichaud/produce_demo > produce_demo.log 2>&1   ```

## Produce Details

Produce is defined by 3 fields.
* Produce Name is the name of the produce.  It is alphanumeric and may contain spaces.  Produce Name may not start or end with spaces.  
* Produce Code uniquely identifies the produce.  It must be 19 characters long, consisting of 4 groups of alphanumeric characters separated by dashes, and is case insensitive. 
* Unit Price is the cost of the produce.  Unit price may optionally start with a '$' (automatically removed), must contain a decimal point and will be padded with a leading zero before the decimal point and padded with up to 2 trailing zeros after the decimal point.


## API calls
### Fetching:
Produce items can be fetched via GET to /produce. This will result in all produce items being fetched.

GET /produce/(Produce Code) will fetch an array containing only the produce item with that particular Produce Code.  

Note that the Produce Code will be validated before attempting to fetch.  Fetching with an invalid Produce Code will result in an error being returned. 

```
Fetching examples:
	curl http://127.0.0.1:8080/produce
	curl http://127.0.0.1:8080/produce/AAAA-1111-2222-3333

Possible Returns:
	(StatusOK|200) 			{"Produce":[{"Produce Code":"A12T-4GH7-QPL9-3N4M","Name":"Lettuce","Unit Price":"3.46"}]}
        (StatusNoContent|204)		{"Error":"No produce found"}
	(StatusBadRequest|400)		{"Error":"Bad Produce Code"}
	(StatusInternalServerError|500)	{"Error":"Internal Error detected"}
```
 
### Adding:
Produce items can be added by calling /produce with a JSON object of either Produce or an array of Produce. \
All three fields of Produce ("Produce Code", "Name" and "Unit Price") must be defined.  

The return will contain an array of Produce (if any) that have been added to database. \
The return will contain an array of Rejected Produce (if any) and the associated errors. \
If the Produce array could not be determined, the Reject Produce will also be returned without a Produce and with appropriate errors.

```
Adding:
	curl -d '[ {"Produce Code": "AAAA-1111-2222-3333", "Name": "Fuji Apples", "Unit Price": "200.60" } ]' -X POST http://127.0.0.1:8080/produce
	curl -d '[ {"Produce Code": "AAAA-1111-2222-3333", "Name": "Fuji Apples", "Unit Price": "200.60" }, {"Produce Code": "AAAB-1111-2222-3333", "Name": "Celery", "Unit Price": ".45" }, {"Produce Code": "AAAC-1111-2222-3333", "Name": "Corn", "Unit Price": "$.5" } ]' -X POST http://127.0.0.1:8080/produce
	curl -d '{"Produce Code": "AAAA-1111-2222-3333", "Name": "Fuji Apples", "Unit Price": "200.60" }' -X POST http://127.0.0.1:8080/produce

Possible Returns:
	(StatusOK|200)			{"Produce":[{"Produce Code":"BBBB-1111-2222-3333","Name":"Black Truffles","Unit Price":"200.6"}]}
	(StatusBadRequest|400)		{"Rejected Produce":[{"Errors":["Failed to read request body"]}]}
        (StatusBadRequest|400)          {"Rejected Produce":[{"Errors":["Failed to unmarshal request body"]}]}
        (StatusBadRequest|400)  	{"Rejected Produce":[{"Produce":{"Produce Code":"BBBB-1111-2222-3333-","Name":" Black Truffles ","Unit Price":"200.645"},"Errors":["Detected error for Produce Code (BBBB-1111-2222-3333-)","Detected error for Produce Name ( Black Truffles )","Detected error for Produce Unit Price (200.645)"]}]}
        (StatusPartialContent|206)	{"Rejected Produce":[{"Produce":{"Produce Code":"AAAA-1111-2222-3333","Name":"Black Truffles","Unit Price":"200.6"},"Errors":["AAAA-1111-2222-3333 already exists"]}]}
        (StatusPartialContent|206)      {"Produce":[{"Produce Code":"AAAA-1111-2222-9999","Name":"Red Peppers","Unit Price":"10.60"}],"Rejected Produce":[{"Produce":{"Produce Code":"AAAA-1111-2222-3333","Name":"Black Truffles","Unit Price":"200.6"},"Errors":["AAAA-1111-2222-3333 already exists"]}]}
```

### Deleting:
Produce items can be removed by calling /produce/(Produce Code).   

Note that the Produce Code will be validated before attempting to delete.  Deleting with an invalid Produce Code will result in an error being returned. 

```
Delete:
	curl -X "DELETE" http://127.0.0.1:8080/produce/AAAA-1111-2222-3333

Possible Returns:
	(StatusOK|200)          	{"Msg":"Produce AAAA-1111-2222-3333 deleted"}
	(StatusBadRequest|400)  	{"Error":"Bad Produce Code"}
	(StatusNotFound|404)    	{"Error":"Produce not found"}
```

# Assumptions
* The echo framework is acceptable for this API.
* Produce Code is unique for all produce items
  * Produce Code is 19 characters long, consisting of 4 groups of alphanumeric characters separated by dashes (case insensitive).  If a Produce Code does not meet this format, it will be rejected.
  * Produce Code is case insenitive, however it will be internally stored in uppercase format.
* Produce Unit Price may optionally start with a '$' (automatically removed) and will be padded with a leading zero before the decimal point and padded with up to 2 trailing zeros after the decimal point.
  * Produce Unit Price lacking a decimal point will be rejected
  * Produce Unit Price with more than 2 digits after the decimal point will be rejected
  * Produce Unit Price with symbols other than digits, a single decimal point and an optionally leading with a '$' will be rejected
  * Produce Unit Price is stored as string.  It is assumed that no arthimetic operations will occur on Produce Unit Price.
* Produce Name is alphanumeric and may contain spaces.  Produce Name may not start or end with spaces.  Produce Name not meeting this format will be rejected.
* All work will be done in feature branches and those branches will be appropriately tagged with annotated tags before merging into master.  
  * Annotation tags will utilize semantic versioning (https://semver.org/).
  * It is the responsibility of developers to ensure the code is appropriately tagged.
  * Developer will not need to log into the container itself (ie: docker exec -it /bin/bash)  
  * Assumption: GitFlow (https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) is to be used.
* Assumption: A external logging mechanism is not needed.

