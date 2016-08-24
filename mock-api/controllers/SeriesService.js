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

exports.seriesObservationsGET = function(args, res, next) {
  /**
   * parameters expected in the args:
  * id (BigDecimal)
  * transforms (List)
  **/
    var examples = {};
  examples['application/json'] = {
  "observationStart" : "1995-01-01",
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
    }, {
      "date" : "2013-01-01",
      "value" : "618.575"
    }, {
      "date" : "2012-01-01",
      "value" : "606.283"
    }, {
      "date" : "2011-01-01",
      "value" : "593.392"
    }, {
      "date" : "2010-01-01",
      "value" : "586.908"
    }, {
      "date" : "2009-01-01",
      "value" : "591.492"
    }, {
      "date" : "2008-01-01",
      "value" : "619.233"
    }, {
      "date" : "2007-01-01",
      "value" : "624.875"
    }, {
      "date" : "2006-01-01",
      "value" : "617.108"
    }, {
      "date" : "2005-01-01",
      "value" : "601.642"
    }, {
      "date" : "2004-01-01",
      "value" : "583.408"
    }, {
      "date" : "2003-01-01",
      "value" : "567.633"
    }, {
      "date" : "2002-01-01",
      "value" : "556.792"
    }, {
      "date" : "2001-01-01",
      "value" : "555.000"
    }, {
      "date" : "2000-01-01",
      "value" : "551.358"
    }, {
      "date" : "1999-01-01",
      "value" : "535.025"
    }, {
      "date" : "1998-01-01",
      "value" : "531.283"
    }, {
      "date" : "1997-01-01",
      "value" : "531.617"
    }, {
      "date" : "1996-01-01",
      "value" : "530.733"
    }, {
      "date" : "1995-01-01",
      "value" : "532.883"
    } ]
  }, {
    "transformation" : "pc1",
    "observations" : [ {
      "date" : "2015-01-01",
      "value" : "1.5"
    }, {
      "date" : "2014-01-01",
      "value" : "1.4"
    }, {
      "date" : "2013-01-01",
      "value" : "2.0"
    }, {
      "date" : "2012-01-01",
      "value" : "2.2"
    }, {
      "date" : "2011-01-01",
      "value" : "1.1"
    }, {
      "date" : "2010-01-01",
      "value" : "-0.8"
    }, {
      "date" : "2009-01-01",
      "value" : "-4.5"
    }, {
      "date" : "2008-01-01",
      "value" : "-0.9"
    }, {
      "date" : "2007-01-01",
      "value" : "1.3"
    }, {
      "date" : "2006-01-01",
      "value" : "2.6"
    }, {
      "date" : "2005-01-01",
      "value" : "3.1"
    }, {
      "date" : "2004-01-01",
      "value" : "2.8"
    }, {
      "date" : "2003-01-01",
      "value" : "1.9"
    }, {
      "date" : "2002-01-01",
      "value" : "0.3"
    }, {
      "date" : "2001-01-01",
      "value" : "0.7"
    }, {
      "date" : "2000-01-01",
      "value" : "3.1"
    }, {
      "date" : "1999-01-01",
      "value" : "0.7"
    }, {
      "date" : "1998-01-01",
      "value" : "-0.1"
    }, {
      "date" : "1997-01-01",
      "value" : "0.2"
    }, {
      "date" : "1996-01-01",
      "value" : "-0.4"
    }, {
      "date" : "1995-01-01",
      "value" : "-0.6"
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

