# osoba
Static site deployment manager

# make
```
$ # get source code
$ git clone git@github.com:w-haibara/osoba.git

$ # build and run the container to make source
$ cd osoba
$ ./script/osoba-make-setup

$ # you can execute "make init" or "make test" or "make" in container
$ ./script/osoba-make init
$ ./script/osoba-make test
$ ./script/osoba-make
```

# run
```
$ ./script/osoba-run
```

# scripts
## path
There are some script to make or run project in `/script`.
You can add this path to `$PATH` by bellow command.
```
export PATH=/path/to/the/osoba/script:$PATH
```
## reference
| script name       | description                                                                                                                         | 
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------- | 
| osoba-make-setup  | Build and run the container to use `make` command.                                                                                  | 
| osoba-make        | Exec make command with options in the container (ex. `osoba-make init`, `osoba-make test`).<br>The containers name is `osoba-make`. | 
| osoba-make-clear  | Clear the container and image relative to `make`.                                                                                   | 
| osoba-data-create | Create data volume contaiiners from each images.                                                                                    | 
| osoba-data-clear  | Clear the data volume containers                                                                                                    | 
| osoba-run<br>     | Build and run the container to exec osoba.<br>The containers is named `osoba`.                                                      | 
| osoba-clear       | Clear the container to exec osoba.                                                                                                  | 
