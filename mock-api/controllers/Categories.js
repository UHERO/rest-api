'use strict';

var url = require('url');


var Categories = require('./CategoriesService');


module.exports.categoryChildrenGET = function categoryChildrenGET (req, res, next) {
  Categories.categoryChildrenGET(req.swagger.params, res, next);
};

module.exports.categoryGET = function categoryGET (req, res, next) {
  Categories.categoryGET(req.swagger.params, res, next);
};

module.exports.categorySeriesGET = function categorySeriesGET (req, res, next) {
  Categories.categorySeriesGET(req.swagger.params, res, next);
};
