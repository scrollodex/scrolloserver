# scrolloserver

A web frontend.



## Build the image

    docker build --tag scrolloserver:latest .

## Run the image

    docker run --publish 8080:8080 scrolloserver

# Run from hugo-bi:
(not tested)

    docker run --rm -v $(pwd)/public:/hugo-bi/public/ -e AIRTABLE_APIKEY=REDACTED -e AIRTABLE_BASE_ID=REDACTED -it hugo-bi
