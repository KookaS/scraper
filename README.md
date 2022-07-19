# scraper

Online scraper for building a dataset for ML.


## installation

Install Golang and MongoDB

    git clone git@github.com:KookaS/scraper.git

If pbm with package `<package>: command not found`:

    export GOPATH="$HOME/go"
    PATH="$GOPATH/bin:$PATH"

    
## run

    go run src/main.go

## build

    go build -o scraper src/main.go

## License

must share photos generated with https://creativecommons.org/licenses/by-sa/2.0/

## .env

    TAGS_UNWANTED=tagsUnwanted
    TAGS_WANTED=tagsWanted
    IMAGES_WANTED=imagesWanted
    IMAGES_PENDING=imagesPending
    IMAGES_UNWANTED=imagesUnwanted
    USERS_UNWANTED=usersUnwanted
    IMAGE_PATH=<absolutePath>
    FLICKR_PRIVATE_KEY=***
    FLICKR_PUBLIC_KEY=***
    UNSPLASH_PRIVATE_KEY=***
    UNSPLASH_PUBLIC_KEY=***
    PEXELS_PUBLIC_KEY=***

## linter

https://github.com/mgechev/revive

    revive -config revive.toml

## Docker

    docker build -t scraper-img .
    docker run -it --rm --name scraper-run scraper-img


