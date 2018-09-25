## Gotrics

Gotrics is go source analyzer.


## Demo

https://gotrics.herokuapp.com/


## Available Metrics

| Name           | Description                                                       |
|:---------------|:------------------------------------------------------------------|
| MethodLength   | calc function or method length                                    |
| ParameterList  | calc function or method parameter count (not include `_`)         |
| MethodNesting  | calc function or method nesting level                             |
| ABCSize        | calc function or method [abc size](http://wiki.c2.com/?AbcMetric) |


### MethodLength

`gotrics` calculated method or function length under the following conditions.

```
Line count of from start `{` line number to end `}` line number.
```