godoc2puml
==========

godoc2puml converts godoc formatted documentation into plantuml format.

Installation
------------

Just type the following to install the program:

    $ go get -u github.com/t-yuki/godoc2puml/...

If you want to draw a diagram as a image file, it depends on java and graphviz. you must install them:

    $ sudo yum install java graphviz
    $ # or, etc...
    $ sudo apt-get install java graphviz

It also depends on plantuml. Well, it is already attached so you do not need to install separately.

Usage
-----

`godoc2puml` generates the diagram for a package. You can also specify output format if you installed the above.

    $ godoc2puml net/http > net.http.puml
    $ # NOT IMPLEMENTED YET # godoc2puml -t=png net/http > net.http.png

Other options:

```
Usage of godoc2puml:
  -t="puml": output format
        puml:  write PlantUML format.
  -h=false: show this help
```

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
