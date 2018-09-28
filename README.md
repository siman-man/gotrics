## Gotrics

Gotrics is go source analyzer.


## Demo

https://gotrics.herokuapp.com/


## Available Metrics

| Name           | Description                                             |
|:---------------|:--------------------------------------------------------|
| FuncLength     | calc function length                                    |
| ParameterList  | calc function parameter count (not include `_`)         |
| FuncNesting    | calc function nesting level                             |
| ABCSize        | calc function [abc size](http://wiki.c2.com/?AbcMetric) |


### FuncLength

`gotrics` calculated function length under the following conditions.

```
Line count of from start `{` line number to end `}` line number.
```