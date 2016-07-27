'use strict';

exports.seriesCategoriesGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  **/
    var examples = {};
  examples['application/json'] = {
  "categories" : [ {
    "id" : 3,
    "name" : "Major Economic Indicator Summary",
    "parent" : 1
  }, {
    "id" : 4,
    "name" : "Summary of External Indicators",
    "parent" : 1
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

exports.seriesGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * id (BigDecimal)
  * search_text (String)
  * freq (String)
  * geo (String)
  **/
    var examples = {};
  examples['application/json'] = {
  "series" : [ {
    "id" : "E_NF@HI.A",
    "title" : "NonFarm Jobs",
    "observationStart" : "1958-01-01",
    "observationEnd" : "2015-01-01",
    "frequency" : "Annual",
    "frequencyShort" : "A",
    "units" : "Thousands",
    "unitsShort" : "Thou.",
    "SeasonalAdjustment" : "Not Seasonally Adjusted",
    "seasonalAdjustmentShort" : "NSA",
    "source" : "U.S. Dept. of Labor, Bureau of Labor Statistics",
    "sourceLink" : "http://data.bls.gov/cgi-bin/srgate"
  }, {
    "id" : "E_NF@HI.M",
    "title" : "NonFarm Jobs",
    "observationStart" : "1958-01-01",
    "observationEnd" : "2016-06-01",
    "frequency" : "Monthly",
    "frequencyShort" : "M",
    "units" : "Thousands",
    "unitsShort" : "Thou.",
    "SeasonalAdjustment" : "Seasonally Adjusted",
    "seasonalAdjustmentShort" : "SA",
    "source" : "U.S. Dept. of Labor, Bureau of Labor Statistics",
    "sourceLink" : "http://data.bls.gov/cgi-bin/srgate"
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

exports.seriesObservationsGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * id (BigDecimal)
  * transforms (List)
  **/
    var examples = {};
  examples['application/json'] = {
  "observationStart" : "2014-01-01",
  "observationEnd" : "2015-01-01",
  "orderBy" : "observationDate",
  "sortOrder" : "desc",
  "transformationResults" : [ {
    "transformation" : "lvl",
    "observations" : [ {
      "date" : "2015-01-01",
      "value" : "636.867"
    }, {
      "date" : "2014-01-01",
      "value" : "627.225"
    } ]
  }, {
    "transformation" : "pc1",
    "observations" : [ {
      "date" : "2015-01-01",
      "value" : "1.5"
    }, {
      "date" : "2014-01-01",
      "value" : "1.4"
    } ]
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

