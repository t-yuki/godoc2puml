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

`godoc2puml` generates a diagram for packages. You can also specify output format if you installed the above.

    $ godoc2puml net/http > net.http.puml
    $ java -jar plantuml.jar -pipe -tpng < net.http.puml > net.http.png

    $ # NOT IMPLEMENTED YET # godoc2puml -t=png net/http > net.http.png

Other options:

```
Usage of godoc2puml:
  -t="puml": output format
        puml:  write PlantUML format.
  -h=false: show this help
```

Examples
--------
The below is output example of "image" package. For more examples, see #1

[image](https://cloud.githubusercontent.com/assets/3804806/3258061/1a0a6f32-f235-11e3-8648-89b9e9abd326.png)

Known Problems
--------------
Many, but...

* enum pattern is not recognized
* noisy
* elementType of map is not recognized as an association
* outer package type is always recognized as classes, except implemented interfaces
* ...

References
----------
* [Plant UML](http://plantuml.sourceforge.net/)

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
