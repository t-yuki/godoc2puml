godoc2puml
==========

godoc2puml converts godoc formatted documentation into plantuml format.

Installation
------------

Just type the following to install the program:

    $ go get -u github.com/t-yuki/godoc2puml

If you want to draw a diagram as a image file, it depends on java and graphviz. you must install them:

    $ sudo yum install java graphviz
    $ # or, etc...
    $ sudo apt-get install openjdk-7-jre-headless graphviz

It also depends on plantuml. <del>Well, it is already attached so you do not need to install separately.</del>

Usage
-----

`godoc2puml` generates a diagram for packages. You can also specify output format if you installed the above.

    $ godoc2puml net/http > net.http.puml
    $ java -jar plantuml.jar -pipe -tpng < net.http.puml > net.http.png

    $ # NOT IMPLEMENTED YET # godoc2puml -t=png net/http > net.http.png

Other options:

```
Usage of godoc2puml:
  -dont-ignore string
        white-list for ignore. default/empty value means packages of arg
  -field string
        set package names in comma-separated strings that use field relationship instead of association
  -h    show this help
  -ignore string
        name filter to ignore. default value removes fmt.String and private declarations except specified packages (default "(fmt\\.Stringer|\\w+\\.[a-z][\\w]*)$")
  -lolipop string
        set package names in comma-separated strings that use lolipop-style interface relationship instead of implementation
  -scope string
        set analysis scope (main) package. if it is omitted, scope is tests in the same directory
  -t string
        output format
        puml:  write PlantUML format (default "text")
```

Examples
--------
The below is output example of "image" package. For more examples, see #1

![image](https://cloud.githubusercontent.com/assets/3804806/3258061/1a0a6f32-f235-11e3-8648-89b9e9abd326.png)

Known Problems
--------------
Many, but...

* enum pattern is not recognized
* noisy
* elementType of map is not recognized as an association
* interface extensions is recognized only when explicit extensions, not implicit extension
* cant parse array- or nested-pointer- embed structs
* cant detect constructors
* ...

References
----------
* [Plant UML](http://plantuml.sourceforge.net/)

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
