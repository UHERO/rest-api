'use strict';

exports.transformsGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  **/
    var examples = {};
  examples['application/json'] = {
  "transformations" : [ {
    "id" : "lvl",
    "description" : "Value in levels",
    "formula" : "x(t)"
  }, {
    "id" : "chg",
    "description" : "Change since last value",
    "formula" : "x(t) - x(t-1)"
  }, {
    "id" : "pc1",
    "description" : "Percent change since last year",
    "formula" : "(x(t)/x(t-n) - 1) * 100, where n is the number of observations in a year"
  } ]
};
  if(Object.keys(examples).length > 0) {
    res.setHeader('Content-Type', 'application/json');
    res.end(JSON.stringify(examples[Object.keys(examples)[0]] || {}, null, 2));
  }
  else {
    res.end();
  }
  
}

