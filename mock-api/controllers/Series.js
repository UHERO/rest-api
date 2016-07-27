'use strict';

var url = require('url');


var Series = require('./SeriesService');


module.exports.seriesCategoriesGET = function seriesCategoriesGET (req, res, next) {
  Series.seriesCategoriesGET(req.swagger.params, res, next);
};

module.exports.seriesGET = function seriesGET (req, res, next) {
  Series.seriesGET(req.swagger.params, res, next);
};

module.exports.seriesObservationsGET = function seriesObservationsGET (req, res, next) {
  Series.seriesObservationsGET(req.swagger.params, res, next);
};
