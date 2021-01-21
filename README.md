# osoba
static site deployment manager

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

# scripts path
There are some script to make or run project in `/script`.
If you cloned osoba in `/home/alice`, you can add this path to `$PATH` by bellow command.
```
export PATH=/home/alice/osoba/script:$PATH
```
