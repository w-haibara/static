# osoba
static site deployment manager

# make
$ # get source code
$ git clone git@github.com:w-haibara/osoba.git
$ # build and run the container to make source
$ cd osoba
$ ./script/osoba-make-setup
$ # you can execute "make init" or "make test" or "make" in container
$ ./script/osoba-make init
$ ./script/osoba-make test
$ ./script/osoba-make

# run
$ ./script/osoba-run
