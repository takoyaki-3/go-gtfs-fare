package gogtfsfare

import (
	"errors"

	csvtag "github.com/artonge/go-csv-tag/v2"
)

type FareAttribute struct {
	FareId           string  `csv:"fare_id"`
	Price            float64 `csv:"price"`
	CurrentType      string  `csv:"currency_type"`
	PaymentMethod    int     `csv:"payment_method"`
	Transfers        int     `csv:"transfer"`
	AgencyId         string  `csv:"agency_id"`
	TransferDuration string  `csv:"transfer_duration"`
}

type FareRule struct {
	FareId        string `csv:"fare_id"`
	RouteId       string `csv:"route_id"`
	OriginId      string `csv:"origin_id"`
	DestinationId string `csv:"destination_id"`
	ContainsId    string `csv:"contains_id"`
}

type GtfsFareData struct {
	FareAttributes []FareAttribute
	FareRules []FareRule
}

func LoadGTFS(dirPath string)(g *GtfsFareData,err error){
	g = &GtfsFareData{}
	if err = csvtag.LoadFromPath(dirPath+"/fare_attributes.txt",&g.FareAttributes);err != nil{
		return g,err
	}
	if err = csvtag.LoadFromPath(dirPath+"/fare_rules.txt",&g.FareRules);err != nil{
		return g,err
	}
	return g,err
}

func fareAttribute(gf *GtfsFareData,fareId string)(FareAttribute,error){
	for _,v:=range gf.FareAttributes{
		if v.FareId == fareId {
			return v,nil
		}
	}
	return FareAttribute{}, errors.New("muched fare_attribute not found.")
}

func GetFareAttribute(gf *GtfsFareData,originId string,destinationId string,routeId string)(FareAttribute,error){

	for _,v:=range gf.FareRules{

		// 条件に一致するか判定
		isOriginOk := (originId == v.OriginId)
		isDestinationOk := (destinationId == v.DestinationId) || ("*" == v.DestinationId) || ("" == v.DestinationId)
		isRouteIdOk := (v.RouteId == routeId) || (v.RouteId == "*") || (v.RouteId == "")

		if isOriginOk && isDestinationOk && isRouteIdOk{
			return fareAttribute(gf,v.FareId)
		}
	}

	return FareAttribute{}, errors.New("muched fare_rule not found.")
}
