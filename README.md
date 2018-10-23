## Gotrics

Gotrics is go source analyzer.


## Demo

https://gotrics.herokuapp.com/

Demo page works on Heroku free dyno. so, if dyno is sleeping mode, maybe page launch time is slow. (about 30 sec)


## Available Metrics

| Name           | Description                                             |
|:---------------|:--------------------------------------------------------|
| FuncLength     | calc function length                                    |
| ParameterCount | calc function parameter count (not include `_`)         |
| FuncNesting    | calc function nesting level                             |
| ABCSize        | calc function [abc size](http://wiki.c2.com/?AbcMetric) |


### FuncLength

`gotrics` calculated function length under the following conditions.

```
Line count of from start `{` line number to end `}` line number.
```


### ParameterCount

Count a number of function arguments.


### FuncNesting

Count a number of the leading tab space.


### ABCSize

Calc abc size.
