'use strict';

exports.geoGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * fips (String)
  * search_text (String)
  **/
    var examples = {};
  examples['application/json'] = {
  "geographies" : [ {
    "fips" : "15",
    "name" : "State of Hawaii",
    "handle" : "HI"
  }, {
    "fips" : "JA",
    "name" : "Japan",
    "handle" : "JP"
  }, {
    "fips" : "15001",
    "name" : "Hawaii County",
    "handle" : "HAW"
  }, {
    "fips" : "15003",
    "name" : "Honolulu",
    "handle" : "HON"
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

