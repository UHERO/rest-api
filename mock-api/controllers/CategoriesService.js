'use strict';

exports.categoryChildrenGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * id (BigDecimal)
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

exports.categoryGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * id (BigDecimal)
  * search_text (String)
  * top_level (Boolean)
  **/
    var examples = {};
  examples['application/json'] = {
  "categories" : [ {
    "id" : 1,
    "name" : "Summary"
  }, {
    "id" : 2,
    "name" : "Income"
  }, {
    "id" : 3,
    "name" : "Major Economic Indicator Summary",
    "parent" : 1
  }, {
    "id" : 4,
    "name" : "Summary of External Indicators",
    "parent" : 1
  }, {
    "id" : 5,
    "name" : "Personal Income Summary",
    "parent" : 2
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

exports.categorySeriesGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * id (BigDecimal)
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
  }, {
    "id" : "VIS@HON.A",
    "title" : "Total Visitor Arrivals",
    "observationStart" : "1990-01-01",
    "observationEnd" : "2015-01-01",
    "frequency" : "Annual",
    "frequencyShort" : "A",
    "units" : "Thousands",
    "unitsShort" : "Thou.",
    "SeasonalAdjustment" : "Not Seasonally Adjusted",
    "seasonalAdjustmentShort" : "NSA",
    "source" : "Hawaii Dept. of Business, Economic Development and Tourism",
    "sourceLink" : "http://www.hawaii.gov/dbedt/info/visitor-stats/"
  }, {
    "id" : "VIS@HON.M",
    "title" : "Total Visitor Arrivals",
    "observationStart" : "1990-01-01",
    "observationEnd" : "2015-01-01",
    "frequency" : "Monthly",
    "frequencyShort" : "M",
    "units" : "Thousands",
    "unitsShort" : "Thou.",
    "SeasonalAdjustment" : "Not Seasonally Adjusted",
    "seasonalAdjustmentShort" : "NSA",
    "source" : "Hawaii Dept. of Business, Economic Development and Tourism",
    "sourceLink" : "http://www.hawaii.gov/dbedt/info/visitor-stats/"
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

